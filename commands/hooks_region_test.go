package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeSharedConfig points the AWS environment variables at scratch files for the
// duration of a test.
func writeSharedConfig(t *testing.T, config, credentials string) {
	t.Helper()
	dir := t.TempDir()

	cfgPath := filepath.Join(dir, "config")
	if err := os.WriteFile(cfgPath, []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}
	credPath := filepath.Join(dir, "credentials")
	if err := os.WriteFile(credPath, []byte(credentials), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("AWS_CONFIG_FILE", cfgPath)
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", credPath)
	// Otherwise a lookup that falls through waits out the instance metadata timeout.
	t.Setenv("AWS_EC2_METADATA_DISABLED", "true")
}

func TestHasEmbeddedRegionReadsTheProfileRegion(t *testing.T) {
	writeSharedConfig(t, `[profile withregion]
region = eu-west-3
`, "")

	region, embedded, err := hasEmbeddedRegionInSharedConfigForProfile("withregion")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !embedded {
		t.Error("expected the region to be reported as embedded")
	}
	if region != "eu-west-3" {
		t.Errorf("got %q, want eu-west-3", region)
	}
}

func TestHasEmbeddedRegionReportsAbsence(t *testing.T) {
	writeSharedConfig(t, `[profile noregion]
output = json
`, "")

	region, embedded, err := hasEmbeddedRegionInSharedConfigForProfile("noregion")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if embedded {
		t.Errorf("expected no embedded region, got %q", region)
	}
}

// The regression this replaced: reading a region used to resolve the whole credential
// chain, so a profile requiring MFA failed with "AssumeRoleTokenProvider session option
// not set" before awless could even establish which region to use. Looking up a region
// must not touch credentials.
func TestHasEmbeddedRegionWorksForAnMFAProfile(t *testing.T) {
	writeSharedConfig(t, `[profile needsmfa]
role_arn = arn:aws:iam::123456789012:role/Admin
source_profile = base
mfa_serial = arn:aws:iam::123456789012:mfa/jsmith
region = us-west-2
`, `[base]
aws_access_key_id = AKIAIOSFODNN7EXAMPLE
aws_secret_access_key = wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
`)

	region, embedded, err := hasEmbeddedRegionInSharedConfigForProfile("needsmfa")
	if err != nil {
		if strings.Contains(err.Error(), "AssumeRoleTokenProvider") {
			t.Fatalf("reading a region resolved credentials and hit the MFA requirement: %s", err)
		}
		t.Fatalf("unexpected error: %s", err)
	}
	if !embedded || region != "us-west-2" {
		t.Errorf("got region %q embedded=%v, want us-west-2 true", region, embedded)
	}
}

// AWS_CONFIG_FILE has to be honored explicitly, because LoadSharedConfigProfile does not
// consult it the way LoadDefaultConfig does. Getting this wrong silently reads the
// developer's real ~/.aws/config instead.
func TestHasEmbeddedRegionHonorsTheConfigFileEnvVar(t *testing.T) {
	writeSharedConfig(t, `[profile scratch]
region = ap-southeast-2
`, "")

	region, _, err := hasEmbeddedRegionInSharedConfigForProfile("scratch")
	if err != nil {
		t.Fatalf("the profile from AWS_CONFIG_FILE was not found: %s", err)
	}
	if region != "ap-southeast-2" {
		t.Errorf("got %q, want ap-southeast-2 from the scratch config", region)
	}
}
