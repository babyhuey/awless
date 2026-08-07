package config

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestIsSemverUpgradeExtra(t *testing.T) {
	tcases := []struct {
		current, latest string
		exp             bool
	}{
		// minor upgrade
		{current: "v0.1.0", latest: "v0.2.0", exp: true},
		// same version
		{current: "v0.1.0", latest: "v0.1.0", exp: false},
		// downgrade
		{current: "v0.2.0", latest: "v0.1.0", exp: false},
		// major upgrade
		{current: "v1.0.0", latest: "v2.0.0", exp: true},
		// empty strings
		{current: "", latest: "", exp: false},
		// malformed versions
		{current: "abc", latest: "def", exp: false},
		{current: "v1.0", latest: "v2.0", exp: false},
		{current: "v1.0.0.0", latest: "v2.0.0.0", exp: false},
		{current: "v1.x.0", latest: "v2.0.0", exp: false},
	}
	for _, tc := range tcases {
		if got, want := IsSemverUpgrade(tc.current, tc.latest), tc.exp; got != want {
			t.Fatalf("IsSemverUpgrade(%q, %q): got %t, want %t", tc.current, tc.latest, got, want)
		}
	}
}

func TestSetVolatileGetRoundtrip(t *testing.T) {
	f, e := os.MkdirTemp(".", "test")
	if e != nil {
		t.Fatal(e)
	}
	defer os.RemoveAll(f)
	os.Setenv("__AWLESS_HOME", f)

	origConfigDefs := configDefinitions
	origDefaultsDefs := defaultsDefinitions
	defer func() {
		configDefinitions = origConfigDefs
		defaultsDefinitions = origDefaultsDefs
	}()

	configDefinitions = map[string]*Definition{
		"aws.region": {help: "AWS region", defaultValue: "us-east-1"},
	}
	defaultsDefinitions = map[string]*Definition{
		"instance.type": {defaultValue: "t2.micro"},
	}

	if err := InitConfig(map[string]string{}); err != nil {
		t.Fatal(err)
	}

	// SetVolatile on a config key
	if err := SetVolatile("aws.region", "ap-southeast-1"); err != nil {
		t.Fatal(err)
	}
	v, ok := Get("aws.region")
	if !ok {
		t.Fatal("expected aws.region to be set")
	}
	if got := fmt.Sprint(v); got != "ap-southeast-1" {
		t.Fatalf("got %s, want ap-southeast-1", got)
	}

	// SetVolatile on a defaults key
	if err := SetVolatile("instance.type", "t3.large"); err != nil {
		t.Fatal(err)
	}
	v, ok = Get("instance.type")
	if !ok {
		t.Fatal("expected instance.type to be set")
	}
	if got := fmt.Sprint(v); got != "t3.large" {
		t.Fatalf("got %s, want t3.large", got)
	}

	// Verify volatile change is not persisted after reload
	if err := LoadConfig(); err != nil {
		t.Fatal(err)
	}
	v, ok = Get("aws.region")
	if !ok {
		t.Fatal("expected aws.region after reload")
	}
	if got := fmt.Sprint(v); got != "us-east-1" {
		t.Fatalf("after reload got %s, want us-east-1", got)
	}
}

func TestGetConfigWithPrefixExtra(t *testing.T) {
	origConfig := Config
	defer func() { Config = origConfig }()

	Config = map[string]interface{}{
		"aws.region":     "us-east-1",
		"aws.profile":    "prod",
		"aws.infra.sync": true,
		"other.setting":  "value",
	}

	t.Run("aws prefix", func(t *testing.T) {
		got := GetConfigWithPrefix("aws.")
		if len(got) != 3 {
			t.Fatalf("expected 3 entries with aws. prefix, got %d: %+v", len(got), got)
		}
		if got["aws.region"] != "us-east-1" {
			t.Fatalf("expected aws.region=us-east-1, got %v", got["aws.region"])
		}
	})

	t.Run("no match prefix", func(t *testing.T) {
		got := GetConfigWithPrefix("nonexistent.")
		if len(got) != 0 {
			t.Fatalf("expected 0 entries, got %d", len(got))
		}
	})

	t.Run("empty prefix returns all", func(t *testing.T) {
		got := GetConfigWithPrefix("")
		if len(got) != len(Config) {
			t.Fatalf("expected %d entries, got %d", len(Config), len(got))
		}
	})
}

func TestGetAWSRegion(t *testing.T) {
	origConfig := Config
	origDefaults := Defaults
	defer func() {
		Config = origConfig
		Defaults = origDefaults
	}()

	t.Run("from config", func(t *testing.T) {
		Config = map[string]interface{}{RegionConfigKey: "eu-west-1"}
		Defaults = map[string]interface{}{}
		if got := GetAWSRegion(); got != "eu-west-1" {
			t.Fatalf("got %s, want eu-west-1", got)
		}
	})

	t.Run("from defaults legacy key", func(t *testing.T) {
		Config = map[string]interface{}{}
		Defaults = map[string]interface{}{"region": "us-west-2"}
		if got := GetAWSRegion(); got != "us-west-2" {
			t.Fatalf("got %s, want us-west-2", got)
		}
	})

	t.Run("empty when not set", func(t *testing.T) {
		Config = map[string]interface{}{}
		Defaults = map[string]interface{}{}
		if got := GetAWSRegion(); got != "" {
			t.Fatalf("got %s, want empty", got)
		}
	})
}

func TestGetAWSProfile(t *testing.T) {
	origConfig := Config
	origDefaults := Defaults
	defer func() {
		Config = origConfig
		Defaults = origDefaults
	}()

	t.Run("from config", func(t *testing.T) {
		Config = map[string]interface{}{ProfileConfigKey: "staging"}
		Defaults = map[string]interface{}{}
		if got := GetAWSProfile(); got != "staging" {
			t.Fatalf("got %s, want staging", got)
		}
	})

	t.Run("default profile", func(t *testing.T) {
		Config = map[string]interface{}{}
		Defaults = map[string]interface{}{}
		if got := GetAWSProfile(); got != "default" {
			t.Fatalf("got %s, want default", got)
		}
	})
}

func TestGetAutosync(t *testing.T) {
	origConfig := Config
	origDefaults := Defaults
	defer func() {
		Config = origConfig
		Defaults = origDefaults
	}()

	t.Run("true from config", func(t *testing.T) {
		Config = map[string]interface{}{"autosync": true}
		Defaults = map[string]interface{}{}
		if got := GetAutosync(); !got {
			t.Fatal("expected true")
		}
	})

	t.Run("false from config", func(t *testing.T) {
		Config = map[string]interface{}{"autosync": false}
		Defaults = map[string]interface{}{}
		if got := GetAutosync(); got {
			t.Fatal("expected false")
		}
	})

	t.Run("default true when not set", func(t *testing.T) {
		Config = map[string]interface{}{}
		Defaults = map[string]interface{}{}
		if got := GetAutosync(); !got {
			t.Fatal("expected default true")
		}
	})
}

func TestDefaultParser(t *testing.T) {
	tcases := []struct {
		input string
		exp   interface{}
	}{
		{"42", 42},
		{"true", true},
		{"false", false},
		{"hello", "hello"},
		{"0", 0},
	}
	for _, tc := range tcases {
		got, err := defaultParser(tc.input)
		if err != nil {
			t.Fatalf("defaultParser(%q) error: %v", tc.input, err)
		}
		if got != tc.exp {
			t.Fatalf("defaultParser(%q): got %v (%T), want %v (%T)", tc.input, got, got, tc.exp, tc.exp)
		}
	}
}

func TestBuildInfoString(t *testing.T) {
	t.Run("version only", func(t *testing.T) {
		b := BuildInfo{Version: "v1.0.0"}
		got := b.String()
		if got != "version=v1.0.0" {
			t.Fatalf("got %s, want version=v1.0.0", got)
		}
	})

	t.Run("all fields", func(t *testing.T) {
		b := BuildInfo{
			Version: "v1.0.0",
			Sha:     "abc123",
			Date:    "2024-01-01",
			Arch:    "amd64",
			OS:      "linux",
			For:     "zip",
		}
		got := b.String()
		exp := "version=v1.0.0, commit=abc123, build-date=2024-01-01, build-arch=amd64, build-os=linux, build-for=zip"
		if got != exp {
			t.Fatalf("got %s, want %s", got, exp)
		}
	})

	t.Run("partial fields sha and date only", func(t *testing.T) {
		b := BuildInfo{Version: "v2.0.0", Sha: "def456", Date: "2025-06-15"}
		got := b.String()
		exp := "version=v2.0.0, commit=def456, build-date=2025-06-15"
		if got != exp {
			t.Fatalf("got %s, want %s", got, exp)
		}
	})

	t.Run("partial fields arch and os only", func(t *testing.T) {
		b := BuildInfo{Version: "v3.0.0", Arch: "arm64", OS: "darwin"}
		got := b.String()
		exp := "version=v3.0.0, build-arch=arm64, build-os=darwin"
		if got != exp {
			t.Fatalf("got %s, want %s", got, exp)
		}
	})

	t.Run("for field only", func(t *testing.T) {
		b := BuildInfo{Version: "v4.0.0", For: "brew"}
		got := b.String()
		exp := "version=v4.0.0, build-for=brew"
		if got != exp {
			t.Fatalf("got %s, want %s", got, exp)
		}
	})
}

func TestParseBool(t *testing.T) {
	tcases := []struct {
		input string
		exp   interface{}
		isErr bool
	}{
		{"true", true, false},
		{"false", false, false},
		{"1", true, false},
		{"0", false, false},
		{"TRUE", true, false},
		{"FALSE", false, false},
		{"toto", nil, true},
		{"", nil, true},
		{"yes", nil, true},
	}
	for _, tc := range tcases {
		got, err := parseBool(tc.input)
		if tc.isErr {
			if err == nil {
				t.Fatalf("parseBool(%q): expected error, got nil", tc.input)
			}
		} else {
			if err != nil {
				t.Fatalf("parseBool(%q): unexpected error: %v", tc.input, err)
			}
			if got != tc.exp {
				t.Fatalf("parseBool(%q): got %v, want %v", tc.input, got, tc.exp)
			}
		}
	}
}

func TestParseInt(t *testing.T) {
	tcases := []struct {
		input string
		exp   interface{}
		isErr bool
	}{
		{"0", 0, false},
		{"1", 1, false},
		{"-1", -1, false},
		{"42", 42, false},
		{"999999", 999999, false},
		{"abc", nil, true},
		{"1.5", nil, true},
		{"", nil, true},
	}
	for _, tc := range tcases {
		got, err := parseInt(tc.input)
		if tc.isErr {
			if err == nil {
				t.Fatalf("parseInt(%q): expected error, got nil", tc.input)
			}
		} else {
			if err != nil {
				t.Fatalf("parseInt(%q): unexpected error: %v", tc.input, err)
			}
			if got != tc.exp {
				t.Fatalf("parseInt(%q): got %v, want %v", tc.input, got, tc.exp)
			}
		}
	}
}

func TestGetCheckUpgradeFrequency(t *testing.T) {
	origConfig := Config
	defer func() { Config = origConfig }()

	t.Run("default 8 hours when not set", func(t *testing.T) {
		Config = map[string]interface{}{}
		got := getCheckUpgradeFrequency()
		if got != 8*time.Hour {
			t.Fatalf("got %v, want %v", got, 8*time.Hour)
		}
	})

	t.Run("custom value", func(t *testing.T) {
		Config = map[string]interface{}{"upgrade.checkfrequency": 24}
		got := getCheckUpgradeFrequency()
		if got != 24*time.Hour {
			t.Fatalf("got %v, want %v", got, 24*time.Hour)
		}
	})

	t.Run("negative value", func(t *testing.T) {
		Config = map[string]interface{}{"upgrade.checkfrequency": -1}
		got := getCheckUpgradeFrequency()
		if got != -1*time.Hour {
			t.Fatalf("got %v, want %v", got, -1*time.Hour)
		}
	})

	t.Run("wrong type returns default", func(t *testing.T) {
		Config = map[string]interface{}{"upgrade.checkfrequency": "not-int"}
		got := getCheckUpgradeFrequency()
		if got != 8*time.Hour {
			t.Fatalf("got %v, want %v", got, 8*time.Hour)
		}
	})
}

func TestGetAutosyncLegacyKey(t *testing.T) {
	origConfig := Config
	origDefaults := Defaults
	defer func() {
		Config = origConfig
		Defaults = origDefaults
	}()

	t.Run("from legacy defaults key sync.auto true", func(t *testing.T) {
		Config = map[string]interface{}{}
		Defaults = map[string]interface{}{"sync.auto": true}
		if got := GetAutosync(); !got {
			t.Fatal("expected true from legacy key")
		}
	})

	t.Run("from legacy defaults key sync.auto false", func(t *testing.T) {
		Config = map[string]interface{}{}
		Defaults = map[string]interface{}{"sync.auto": false}
		if got := GetAutosync(); got {
			t.Fatal("expected false from legacy key")
		}
	})

	t.Run("non-bool type in config returns default true", func(t *testing.T) {
		Config = map[string]interface{}{"autosync": "not-a-bool"}
		Defaults = map[string]interface{}{}
		if got := GetAutosync(); !got {
			t.Fatal("expected default true when type assertion fails")
		}
	})
}

func TestGetAWSProfileLegacyKey(t *testing.T) {
	origConfig := Config
	origDefaults := Defaults
	defer func() {
		Config = origConfig
		Defaults = origDefaults
	}()

	t.Run("from legacy defaults key", func(t *testing.T) {
		Config = map[string]interface{}{}
		Defaults = map[string]interface{}{ProfileConfigKey: "legacy-profile"}
		if got := GetAWSProfile(); got != "legacy-profile" {
			t.Fatalf("got %s, want legacy-profile", got)
		}
	})

	t.Run("config takes priority over defaults", func(t *testing.T) {
		Config = map[string]interface{}{ProfileConfigKey: "config-profile"}
		Defaults = map[string]interface{}{ProfileConfigKey: "default-profile"}
		if got := GetAWSProfile(); got != "config-profile" {
			t.Fatalf("got %s, want config-profile", got)
		}
	})

	t.Run("empty config key falls to defaults", func(t *testing.T) {
		Config = map[string]interface{}{ProfileConfigKey: ""}
		Defaults = map[string]interface{}{ProfileConfigKey: "fallback"}
		if got := GetAWSProfile(); got != "fallback" {
			t.Fatalf("got %s, want fallback", got)
		}
	})
}

func TestCompareSemverErrors(t *testing.T) {
	tcases := []struct {
		a, b string
	}{
		{"", ""},
		{"1.0", "2.0"},
		{"abc", "def"},
		{"1.a.0", "1.b.0"},
		{"1.0.0.0", "2.0.0.0"},
	}
	for _, tc := range tcases {
		_, err := CompareSemver(tc.a, tc.b)
		if !errors.Is(err, SemverInvalidFormatErr) {
			t.Fatalf("CompareSemver(%q, %q): expected SemverInvalidFormatErr, got %v", tc.a, tc.b, err)
		}
	}
}

func TestGetNotFound(t *testing.T) {
	origConfig := Config
	origDefaults := Defaults
	defer func() {
		Config = origConfig
		Defaults = origDefaults
	}()

	Config = map[string]interface{}{}
	Defaults = map[string]interface{}{}

	_, ok := Get("nonexistent.key")
	if ok {
		t.Fatal("expected ok=false for nonexistent key")
	}
}

func TestGetFromConfigBeforeDefaults(t *testing.T) {
	origConfig := Config
	origDefaults := Defaults
	defer func() {
		Config = origConfig
		Defaults = origDefaults
	}()

	Config = map[string]interface{}{"shared.key": "from-config"}
	Defaults = map[string]interface{}{"shared.key": "from-defaults"}

	v, ok := Get("shared.key")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if got := fmt.Sprint(v); got != "from-config" {
		t.Fatalf("got %s, want from-config", got)
	}
}

func TestGetFromDefaultsWhenNotInConfig(t *testing.T) {
	origConfig := Config
	origDefaults := Defaults
	defer func() {
		Config = origConfig
		Defaults = origDefaults
	}()

	Config = map[string]interface{}{}
	Defaults = map[string]interface{}{"only.in.defaults": "default-val"}

	v, ok := Get("only.in.defaults")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if got := fmt.Sprint(v); got != "default-val" {
		t.Fatalf("got %s, want default-val", got)
	}
}

func TestSetProfileCallback(t *testing.T) {
	f, e := os.MkdirTemp(".", "test")
	if e != nil {
		t.Fatal(e)
	}
	defer os.RemoveAll(f)
	os.Setenv("__AWLESS_HOME", f)

	origConfigDefs := configDefinitions
	origDefaultsDefs := defaultsDefinitions
	defer func() {
		configDefinitions = origConfigDefs
		defaultsDefinitions = origDefaultsDefs
	}()

	configDefinitions = map[string]*Definition{
		ProfileConfigKey: {help: "AWS profile", defaultValue: "default"},
	}
	defaultsDefinitions = map[string]*Definition{}

	if err := InitConfig(map[string]string{}); err != nil {
		t.Fatal(err)
	}

	if err := SetProfileCallback("myprofile"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := GetAWSProfile(); got != "myprofile" {
		t.Fatalf("got %s, want myprofile", got)
	}
}

func TestRunSyncWithUpdatedRegion(t *testing.T) {
	origConfig := Config
	origTrigger := TriggerSyncOnConfigUpdate
	defer func() {
		Config = origConfig
		TriggerSyncOnConfigUpdate = origTrigger
	}()

	t.Run("does not trigger when autosync is false", func(t *testing.T) {
		Config = map[string]interface{}{"autosync": false}
		TriggerSyncOnConfigUpdate = false
		runSyncWithUpdatedRegion("us-east-1")
		if TriggerSyncOnConfigUpdate {
			t.Fatal("expected TriggerSyncOnConfigUpdate to remain false when autosync is disabled")
		}
	})

	t.Run("does not trigger for invalid region", func(t *testing.T) {
		Config = map[string]interface{}{"autosync": true}
		TriggerSyncOnConfigUpdate = false
		runSyncWithUpdatedRegion("invalid-region-xyz")
		if TriggerSyncOnConfigUpdate {
			t.Fatal("expected TriggerSyncOnConfigUpdate to remain false for invalid region")
		}
	})

	t.Run("triggers for valid region with autosync", func(t *testing.T) {
		Config = map[string]interface{}{"autosync": true}
		TriggerSyncOnConfigUpdate = false
		runSyncWithUpdatedRegion("us-east-1")
		if !TriggerSyncOnConfigUpdate {
			t.Fatal("expected TriggerSyncOnConfigUpdate to be true")
		}
	})
}

func TestSetVolatileWithAWSPrefix(t *testing.T) {
	f, e := os.MkdirTemp(".", "test")
	if e != nil {
		t.Fatal(e)
	}
	defer os.RemoveAll(f)
	os.Setenv("__AWLESS_HOME", f)

	origConfigDefs := configDefinitions
	origDefaultsDefs := defaultsDefinitions
	origConfig := Config
	origDefaults := Defaults
	defer func() {
		configDefinitions = origConfigDefs
		defaultsDefinitions = origDefaultsDefs
		Config = origConfig
		Defaults = origDefaults
	}()

	configDefinitions = map[string]*Definition{}
	defaultsDefinitions = map[string]*Definition{}

	Config = map[string]interface{}{}
	Defaults = map[string]interface{}{}

	// Key with "aws." prefix but not in any definition should go to Config
	if err := SetVolatile("aws.custom.setting", "value"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := Config["aws.custom.setting"]; !ok {
		t.Fatal("expected key with aws. prefix to be in Config")
	}
	if _, ok := Defaults["aws.custom.setting"]; ok {
		t.Fatal("expected key with aws. prefix NOT to be in Defaults")
	}
}

func TestSetVolatileWithoutAWSPrefix(t *testing.T) {
	f, e := os.MkdirTemp(".", "test")
	if e != nil {
		t.Fatal(e)
	}
	defer os.RemoveAll(f)
	os.Setenv("__AWLESS_HOME", f)

	origConfigDefs := configDefinitions
	origDefaultsDefs := defaultsDefinitions
	origConfig := Config
	origDefaults := Defaults
	defer func() {
		configDefinitions = origConfigDefs
		defaultsDefinitions = origDefaultsDefs
		Config = origConfig
		Defaults = origDefaults
	}()

	configDefinitions = map[string]*Definition{}
	defaultsDefinitions = map[string]*Definition{}

	Config = map[string]interface{}{}
	Defaults = map[string]interface{}{}

	// Key without "aws." prefix and not in any definition should go to Defaults
	if err := SetVolatile("custom.setting", "value"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := Defaults["custom.setting"]; !ok {
		t.Fatal("expected key without aws. prefix to be in Defaults")
	}
	if _, ok := Config["custom.setting"]; ok {
		t.Fatal("expected key without aws. prefix NOT to be in Config")
	}
}

func TestNotifyIfUpgradeBrew(t *testing.T) {
	origBuildFor := BuildFor
	defer func() { BuildFor = origBuildFor }()

	BuildFor = "brew"
	tserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"URL":"https://github.com/bootswithdefer/awless/releases/latest","Version":"1000.0.0"}`))
	}))
	defer tserver.Close()

	var buff bytes.Buffer
	if err := notifyIfUpgrade(tserver.URL, &buff); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buff.String(), "brew upgrade awless") {
		t.Fatalf("expected brew upgrade message, got %s", buff.String())
	}
}

func TestNotifyIfUpgradeGoGet(t *testing.T) {
	origBuildFor := BuildFor
	defer func() { BuildFor = origBuildFor }()

	BuildFor = "" // default, not brew or zip
	tserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"URL":"https://github.com/bootswithdefer/awless/releases/latest","Version":"1000.0.0"}`))
	}))
	defer tserver.Close()

	var buff bytes.Buffer
	if err := notifyIfUpgrade(tserver.URL, &buff); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buff.String(), "go get -u") {
		t.Fatalf("expected go get message, got %s", buff.String())
	}
}

func TestNotifyIfUpgradeNoUpgrade(t *testing.T) {
	tserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"URL":"","Version":"0.0.0"}`))
	}))
	defer tserver.Close()

	var buff bytes.Buffer
	if err := notifyIfUpgrade(tserver.URL, &buff); err != nil {
		t.Fatal(err)
	}
	if buff.Len() != 0 {
		t.Fatalf("expected no output for no-upgrade, got %s", buff.String())
	}
}

func TestNotifyIfUpgradeInvalidJSON(t *testing.T) {
	tserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`not json`))
	}))
	defer tserver.Close()

	var buff bytes.Buffer
	// Should not error even with invalid JSON
	if err := notifyIfUpgrade(tserver.URL, &buff); err != nil {
		t.Fatal(err)
	}
	if buff.Len() != 0 {
		t.Fatalf("expected no output for invalid json, got %s", buff.String())
	}
}

func TestNotifyIfUpgradeServerDown(t *testing.T) {
	var buff bytes.Buffer
	err := notifyIfUpgrade("http://localhost:1", &buff)
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
