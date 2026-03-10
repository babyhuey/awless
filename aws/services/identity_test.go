package awsservices

import (
	"strings"
	"testing"
)

func TestIdentityIsRoot(t *testing.T) {
	tests := []struct {
		name     string
		identity Identity
		want     bool
	}{
		{
			name:     "root user is root",
			identity: Identity{Resource: "root"},
			want:     true,
		},
		{
			name:     "IAM user is not root",
			identity: Identity{Resource: "johndoe"},
			want:     false,
		},
		{
			name:     "empty resource is not root",
			identity: Identity{},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.identity.IsRoot(); got != tt.want {
				t.Errorf("IsRoot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIdentityIsUserType(t *testing.T) {
	tests := []struct {
		name     string
		identity Identity
		want     bool
	}{
		{
			name:     "user type returns true",
			identity: Identity{ResourceType: "user"},
			want:     true,
		},
		{
			name:     "role type returns false",
			identity: Identity{ResourceType: "role"},
			want:     false,
		},
		{
			name:     "assumed-role type returns false",
			identity: Identity{ResourceType: "assumed-role"},
			want:     false,
		},
		{
			name:     "empty type returns false",
			identity: Identity{},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.identity.IsUserType(); got != tt.want {
				t.Errorf("IsUserType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArnResourceInfoRegex(t *testing.T) {
	tests := []struct {
		name             string
		arn              string
		wantResourceType string
		wantResource     string
		wantIsRoot       bool
		wantIsUser       bool
	}{
		{
			name:             "root user ARN",
			arn:              "arn:aws:iam::123456789012:root",
			wantResourceType: "user",
			wantResource:     "root",
			wantIsRoot:       true,
			wantIsUser:       true,
		},
		{
			name:             "IAM user ARN",
			arn:              "arn:aws:iam::123456789012:user/johndoe",
			wantResourceType: "user",
			wantResource:     "johndoe",
			wantIsRoot:       false,
			wantIsUser:       true,
		},
		{
			name:             "IAM user with path ARN",
			arn:              "arn:aws:iam::123456789012:user/division_abc/subdivision_xyz/johndoe",
			wantResourceType: "user",
			wantResource:     "division_abc/subdivision_xyz/johndoe",
			wantIsRoot:       false,
			wantIsUser:       true,
		},
		{
			name:             "role ARN",
			arn:              "arn:aws:iam::123456789012:role/myrole",
			wantResourceType: "role",
			wantResource:     "myrole",
			wantIsRoot:       false,
			wantIsUser:       false,
		},
		{
			name:             "assumed role ARN",
			arn:              "arn:aws:sts::123456789012:assumed-role/myrole/mysession",
			wantResourceType: "assumed-role",
			wantResource:     "myrole/mysession",
			wantIsRoot:       false,
			wantIsUser:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the ARN parsing logic from GetIdentity
			ident := &Identity{Arn: tt.arn}

			splits := strings.Split(tt.arn, ":")
			if l := len(splits); l > 0 {
				ident.ResourcePath = splits[l-1]
				matches := arnResourceInfoRegex.FindStringSubmatch(ident.ResourcePath)
				if len(matches) == 4 {
					if matches[1] == "root" {
						ident.Resource = "root"
						ident.ResourceType = "user"
					} else {
						ident.ResourceType = matches[2]
						ident.Resource = matches[3]
					}
				}
			}

			if ident.ResourceType != tt.wantResourceType {
				t.Errorf("ResourceType = %q, want %q", ident.ResourceType, tt.wantResourceType)
			}
			if ident.Resource != tt.wantResource {
				t.Errorf("Resource = %q, want %q", ident.Resource, tt.wantResource)
			}
			if ident.IsRoot() != tt.wantIsRoot {
				t.Errorf("IsRoot() = %v, want %v", ident.IsRoot(), tt.wantIsRoot)
			}
			if ident.IsUserType() != tt.wantIsUser {
				t.Errorf("IsUserType() = %v, want %v", ident.IsUserType(), tt.wantIsUser)
			}
		})
	}
}
