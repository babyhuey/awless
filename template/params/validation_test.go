package params_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wallix/awless/template/params"
)

func TestValidation(t *testing.T) {
	vals := params.Validators{
		"one": params.MaxLengthOf(3),
		"two": params.MinLengthOf(2),
	}

	err := params.Validate(vals, map[string]interface{}{"one": "morethan3", "two": "o"})
	if err == nil {
		t.Fatal("expected error got none")
	}
	msg := err.Error()
	if got, want := msg, "param validation:"; !strings.Contains(got, want) {
		t.Fatalf("expected '%s' to contains: %s", got, want)
	}
	if got, want := msg, "param 'one': expected max length of 3"; !strings.Contains(got, want) {
		t.Fatalf("expected '%s' to contains: %s", got, want)
	}
	if got, want := msg, "param 'two': expected min length of 2"; !strings.Contains(got, want) {
		t.Fatalf("expected '%s' to contains: %s", got, want)
	}
}

func TestValidationAllPass(t *testing.T) {
	vals := params.Validators{
		"one": params.MaxLengthOf(10),
		"two": params.MinLengthOf(1),
	}
	err := params.Validate(vals, map[string]interface{}{"one": "short", "two": "ok"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidationNoMatchingKeys(t *testing.T) {
	vals := params.Validators{
		"one": params.MaxLengthOf(3),
	}
	err := params.Validate(vals, map[string]interface{}{"other": "value"})
	if err != nil {
		t.Fatalf("expected no error when no keys match, got %v", err)
	}
}

func TestValidationEmptyValidators(t *testing.T) {
	vals := params.Validators{}
	err := params.Validate(vals, map[string]interface{}{"any": "value"})
	if err != nil {
		t.Fatalf("expected no error with empty validators, got %v", err)
	}
}

func TestIsInEnumIgnoreCase(t *testing.T) {
	validator := params.IsInEnumIgnoreCase("tcp", "udp", "icmp")

	t.Run("exact match", func(t *testing.T) {
		if err := validator("tcp", nil); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("case insensitive match", func(t *testing.T) {
		if err := validator("TCP", nil); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("mixed case match", func(t *testing.T) {
		if err := validator("Udp", nil); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("not in enum", func(t *testing.T) {
		err := validator("http", nil)
		if err == nil {
			t.Fatal("expected error for non-matching value")
		}
		if !strings.Contains(err.Error(), "http") {
			t.Fatalf("expected error to mention 'http', got %v", err)
		}
	})

	t.Run("non-string type", func(t *testing.T) {
		err := validator(42, nil)
		if err == nil {
			t.Fatal("expected error for non-string type")
		}
	})
}

func TestMaxLengthOf(t *testing.T) {
	validator := params.MaxLengthOf(5)

	t.Run("under limit", func(t *testing.T) {
		if err := validator("abc", nil); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("at limit", func(t *testing.T) {
		if err := validator("abcde", nil); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("over limit", func(t *testing.T) {
		err := validator("abcdef", nil)
		if err == nil {
			t.Fatal("expected error for over-limit string")
		}
		if !strings.Contains(err.Error(), "max length of 5") {
			t.Fatalf("expected error about max length, got %v", err)
		}
	})

	t.Run("empty string", func(t *testing.T) {
		if err := validator("", nil); err != nil {
			t.Fatalf("expected no error for empty string, got %v", err)
		}
	})

	t.Run("non-string type", func(t *testing.T) {
		err := validator(123, nil)
		if err == nil {
			t.Fatal("expected error for non-string type")
		}
	})
}

func TestMinLengthOf(t *testing.T) {
	validator := params.MinLengthOf(3)

	t.Run("over limit", func(t *testing.T) {
		if err := validator("abcdef", nil); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("at limit", func(t *testing.T) {
		if err := validator("abc", nil); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("under limit", func(t *testing.T) {
		err := validator("ab", nil)
		if err == nil {
			t.Fatal("expected error for under-limit string")
		}
		if !strings.Contains(err.Error(), "min length of 3") {
			t.Fatalf("expected error about min length, got %v", err)
		}
	})

	t.Run("non-string type", func(t *testing.T) {
		err := validator(123, nil)
		if err == nil {
			t.Fatal("expected error for non-string type")
		}
	})
}

func TestIsCIDR(t *testing.T) {
	t.Run("valid CIDR", func(t *testing.T) {
		if err := params.IsCIDR("10.0.0.0/24", nil); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("valid CIDR v6", func(t *testing.T) {
		if err := params.IsCIDR("::1/128", nil); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("invalid CIDR", func(t *testing.T) {
		err := params.IsCIDR("not-a-cidr", nil)
		if err == nil {
			t.Fatal("expected error for invalid CIDR")
		}
	})

	t.Run("IP without mask", func(t *testing.T) {
		err := params.IsCIDR("10.0.0.1", nil)
		if err == nil {
			t.Fatal("expected error for IP without mask")
		}
	})

	t.Run("non-string type", func(t *testing.T) {
		err := params.IsCIDR(42, nil)
		if err == nil {
			t.Fatal("expected error for non-string type")
		}
	})
}

func TestIsIP(t *testing.T) {
	t.Run("valid IPv4", func(t *testing.T) {
		if err := params.IsIP("192.168.1.1", nil); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("valid IPv6", func(t *testing.T) {
		if err := params.IsIP("::1", nil); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("invalid IP", func(t *testing.T) {
		err := params.IsIP("not-an-ip", nil)
		if err == nil {
			t.Fatal("expected error for invalid IP")
		}
		if !strings.Contains(err.Error(), "not-an-ip") {
			t.Fatalf("expected error to mention input, got %v", err)
		}
	})

	t.Run("empty string", func(t *testing.T) {
		err := params.IsIP("", nil)
		if err == nil {
			t.Fatal("expected error for empty string")
		}
	})

	t.Run("non-string type", func(t *testing.T) {
		err := params.IsIP(42, nil)
		if err == nil {
			t.Fatal("expected error for non-string type")
		}
	})
}

func TestIsFilepath(t *testing.T) {
	t.Run("existing file", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test-filepath")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())
		tmpFile.Close()

		if err := params.IsFilepath(tmpFile.Name(), nil); err != nil {
			t.Fatalf("expected no error for existing file, got %v", err)
		}
	})

	t.Run("nonexistent file", func(t *testing.T) {
		err := params.IsFilepath("/nonexistent/path/to/file.txt", nil)
		if err == nil {
			t.Fatal("expected error for nonexistent file")
		}
		if !strings.Contains(err.Error(), "cannot find file") {
			t.Fatalf("expected 'cannot find file' error, got %v", err)
		}
	})

	t.Run("directory", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "test-filepath-dir")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)

		err = params.IsFilepath(dir, nil)
		if err == nil {
			t.Fatal("expected error for directory")
		}
		if !strings.Contains(err.Error(), "is a directory") {
			t.Fatalf("expected 'is a directory' error, got %v", err)
		}
	})

	t.Run("non-string type", func(t *testing.T) {
		err := params.IsFilepath(42, nil)
		if err == nil {
			t.Fatal("expected error for non-string type")
		}
	})
}

func TestIsFilepathWithSubdirectory(t *testing.T) {
	dir, err := os.MkdirTemp("", "test-filepath-sub")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	subFile := filepath.Join(dir, "nested.txt")
	if err := os.WriteFile(subFile, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := params.IsFilepath(subFile, nil); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateWithMixedResults(t *testing.T) {
	vals := params.Validators{
		"name":  params.MinLengthOf(1),
		"cidr":  params.IsCIDR,
		"proto": params.IsInEnumIgnoreCase("tcp", "udp"),
	}

	// All valid
	err := params.Validate(vals, map[string]interface{}{
		"name":  "my-resource",
		"cidr":  "10.0.0.0/16",
		"proto": "tcp",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// One invalid
	err = params.Validate(vals, map[string]interface{}{
		"name":  "my-resource",
		"cidr":  "invalid",
		"proto": "tcp",
	})
	if err == nil {
		t.Fatal("expected error for invalid cidr")
	}
	if !strings.Contains(err.Error(), "cidr") {
		t.Fatalf("expected error to mention 'cidr', got %v", err)
	}
}
