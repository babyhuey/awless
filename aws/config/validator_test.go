package awsconfig

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestIsValidRegion(t *testing.T) {
	tcases := []struct {
		name   string
		input  string
		expect bool
	}{
		// Valid standard regions
		{"valid us-east-1", "us-east-1", true},
		{"valid eu-west-1", "eu-west-1", true},
		{"valid ap-southeast-2", "ap-southeast-2", true},
		{"valid af-south-1", "af-south-1", true},
		{"valid me-south-1", "me-south-1", true},
		{"valid il-central-1", "il-central-1", true},
		{"valid cn-north-1", "cn-north-1", true},
		// Valid compound regions
		{"valid us-gov-west-1", "us-gov-west-1", true},
		{"valid us-iso-east-1", "us-iso-east-1", true},
		{"valid us-isob-east-1", "us-isob-east-1", true},
		{"valid eu-isoe-west-1", "eu-isoe-west-1", true},
		// Invalid regions
		{"invalid word", "invalid", false},
		{"incomplete us-east", "us-east", false},
		{"no prefix east-1", "east-1", false},
		{"empty string", "", false},
		{"underscore separator", "us_east_1", false},
		{"uppercase", "US-EAST-1", false},
	}
	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsValidRegion(tc.input); got != tc.expect {
				t.Errorf("IsValidRegion(%q) = %t, want %t", tc.input, got, tc.expect)
			}
		})
	}
}

func TestParseRegion(t *testing.T) {
	t.Run("valid region returns value and nil error", func(t *testing.T) {
		val, err := ParseRegion("us-east-1")
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if val != "us-east-1" {
			t.Errorf("expected value 'us-east-1', got %v", val)
		}
	})

	t.Run("invalid region returns value and error", func(t *testing.T) {
		val, err := ParseRegion("invalid")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if val != "invalid" {
			t.Errorf("expected value 'invalid', got %v", val)
		}
	})
}

func TestIsValidInstanceType(t *testing.T) {
	tcases := []struct {
		name   string
		input  string
		expect bool
	}{
		{"valid t2.micro", "t2.micro", true},
		{"valid m5.xlarge", "m5.xlarge", true},
		{"valid c5.2xlarge", "c5.2xlarge", true},
		{"invalid word", "invalid", false},
		{"empty string", "", false},
		{"no suffix t2", "t2", false},
		{"no prefix .micro", ".micro", false},
	}
	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isValidInstanceType(tc.input); got != tc.expect {
				t.Errorf("isValidInstanceType(%q) = %t, want %t", tc.input, got, tc.expect)
			}
		})
	}
}

func TestParseInstanceType(t *testing.T) {
	t.Run("valid instance type returns value and nil error", func(t *testing.T) {
		val, err := ParseInstanceType("t2.micro")
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if val != "t2.micro" {
			t.Errorf("expected value 't2.micro', got %v", val)
		}
	})

	t.Run("invalid instance type returns value and error", func(t *testing.T) {
		val, err := ParseInstanceType("invalid")
		if err == nil {
			t.Error("expected error, got nil")
		}
		if val != "invalid" {
			t.Errorf("expected value 'invalid', got %v", val)
		}
	})
}

// #302: newer AWS regions like af-south-1 must appear in allRegions()
func TestNewerRegionsInAllRegions(t *testing.T) {
	regions := allRegions()
	newer := []string{
		"af-south-1",
		"ap-east-1",
		"ap-northeast-3",
		"ap-southeast-3",
		"eu-north-1",
		"eu-south-1",
		"me-south-1",
		"me-central-1",
		"il-central-1",
		"ca-west-1",
	}
	for _, r := range newer {
		found := false
		for _, ar := range regions {
			if ar == r {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("allRegions() missing newer region %q", r)
		}
	}
}

func TestAllRegions(t *testing.T) {
	regions := allRegions()

	t.Run("returns non-empty slice", func(t *testing.T) {
		if len(regions) == 0 {
			t.Fatal("allRegions() returned empty slice")
		}
	})

	t.Run("slice is sorted", func(t *testing.T) {
		if !sort.StringsAreSorted(regions) {
			t.Error("allRegions() returned unsorted slice")
		}
	})

	t.Run("all entries are valid regions", func(t *testing.T) {
		for _, r := range regions {
			if !IsValidRegion(r) {
				t.Errorf("allRegions() contains invalid region %q", r)
			}
		}
	})
}

func TestIsValidProfile(t *testing.T) {
	tmpDir := t.TempDir()

	origFunc := awsHomeFunc
	awsHomeFunc = func() string { return tmpDir }
	t.Cleanup(func() { awsHomeFunc = origFunc })

	configContent := []byte("[default]\nregion = us-east-1\n\n[profile myprofile]\nregion = eu-west-1\n")
	if err := os.WriteFile(filepath.Join(tmpDir, "config"), configContent, 0600); err != nil {
		t.Fatal(err)
	}

	tcases := []struct {
		name   string
		input  string
		expect bool
	}{
		{"default profile is valid", "default", true},
		{"myprofile is valid", "myprofile", true},
		{"nonexistent profile is invalid", "nonexistent", false},
		{"empty string is invalid", "", false},
	}
	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsValidProfile(tc.input); got != tc.expect {
				t.Errorf("IsValidProfile(%q) = %t, want %t", tc.input, got, tc.expect)
			}
		})
	}
}

func TestAllProfiles(t *testing.T) {
	tmpDir := t.TempDir()

	origFunc := awsHomeFunc
	awsHomeFunc = func() string { return tmpDir }
	t.Cleanup(func() { awsHomeFunc = origFunc })

	configContent := []byte("[default]\nregion = us-east-1\n\n[profile staging]\nregion = eu-west-1\n")
	credentialsContent := []byte("[default]\naws_access_key_id = EXAMPLE\n\n[production]\naws_access_key_id = EXAMPLE2\n")

	if err := os.WriteFile(filepath.Join(tmpDir, "config"), configContent, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "credentials"), credentialsContent, 0600); err != nil {
		t.Fatal(err)
	}

	profiles := AllProfiles()

	expected := []string{"default", "staging", "default", "production"}
	if len(profiles) != len(expected) {
		t.Fatalf("AllProfiles() returned %d profiles %v, want %d %v", len(profiles), profiles, len(expected), expected)
	}
	for i, p := range expected {
		if profiles[i] != p {
			t.Errorf("AllProfiles()[%d] = %q, want %q", i, profiles[i], p)
		}
	}
}

func TestAllProfilesMissingFiles(t *testing.T) {
	tmpDir := t.TempDir()

	origFunc := awsHomeFunc
	awsHomeFunc = func() string { return tmpDir }
	t.Cleanup(func() { awsHomeFunc = origFunc })

	profiles := AllProfiles()
	if len(profiles) != 0 {
		t.Errorf("AllProfiles() with no files returned %v, want empty", profiles)
	}
}

func TestStringInSlice(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}

	tcases := []struct {
		name   string
		input  string
		expect bool
	}{
		{"found at start", "apple", true},
		{"found in middle", "banana", true},
		{"found at end", "cherry", true},
		{"not found", "durian", false},
		{"empty string not found", "", false},
	}
	for _, tc := range tcases {
		t.Run(tc.name, func(t *testing.T) {
			if got := stringInSlice(tc.input, slice); got != tc.expect {
				t.Errorf("stringInSlice(%q, ...) = %t, want %t", tc.input, got, tc.expect)
			}
		})
	}

	t.Run("empty slice", func(t *testing.T) {
		if stringInSlice("anything", nil) {
			t.Error("stringInSlice should return false for nil slice")
		}
	})
}
