package config

import (
	"fmt"
	"os"
	"testing"
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
		"scheduler.url":  "http://localhost:8082",
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

func TestGetSchedulerURL(t *testing.T) {
	origConfig := Config
	defer func() { Config = origConfig }()

	t.Run("when set", func(t *testing.T) {
		Config = map[string]interface{}{"scheduler.url": "http://example.com:9090"}
		if got := GetSchedulerURL(); got != "http://example.com:9090" {
			t.Fatalf("got %s, want http://example.com:9090", got)
		}
	})

	t.Run("empty when not set", func(t *testing.T) {
		Config = map[string]interface{}{}
		if got := GetSchedulerURL(); got != "" {
			t.Fatalf("got %s, want empty", got)
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
}
