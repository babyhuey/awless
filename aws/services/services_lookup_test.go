package awsservices

import (
	"context"
	"testing"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/graph"
)

// mockService implements cloud.Service for testing.
type mockService struct {
	name string
}

func (m *mockService) Region() string                                { return "us-east-1" }
func (m *mockService) Profile() string                               { return "default" }
func (m *mockService) Name() string                                  { return m.name }
func (m *mockService) ResourceTypes() []string                       { return nil }
func (m *mockService) IsSyncDisabled() bool                          { return false }
func (m *mockService) Fetch(context.Context) (cloud.GraphAPI, error) { return graph.NewGraph(), nil }
func (m *mockService) FetchByType(context.Context, string) (cloud.GraphAPI, error) {
	return graph.NewGraph(), nil
}

func setupServiceRegistry(t *testing.T) func() {
	t.Helper()
	// Save original registry contents
	origRegistry := make(map[string]cloud.Service)
	for k, v := range cloud.ServiceRegistry {
		origRegistry[k] = v
	}

	// Populate with mock services
	serviceNames := []string{"infra", "access", "storage", "messaging", "dns", "lambda", "monitoring", "cdn", "cloudformation"}
	for _, name := range serviceNames {
		cloud.ServiceRegistry[name] = &mockService{name: name}
	}

	// Return cleanup function
	return func() {
		for k := range cloud.ServiceRegistry {
			delete(cloud.ServiceRegistry, k)
		}
		for k, v := range origRegistry {
			cloud.ServiceRegistry[k] = v
		}
	}
}

func TestGetCloudServicesForAPIs(t *testing.T) {
	cleanup := setupServiceRegistry(t)
	defer cleanup()

	tests := []struct {
		name      string
		apis      []string
		wantNames []string
	}{
		{
			name:      "ec2 maps to infra",
			apis:      []string{"ec2"},
			wantNames: []string{"infra"},
		},
		{
			name:      "iam maps to access",
			apis:      []string{"iam"},
			wantNames: []string{"access"},
		},
		{
			name:      "s3 maps to storage",
			apis:      []string{"s3"},
			wantNames: []string{"storage"},
		},
		{
			name:      "multiple APIs from same service deduplicated",
			apis:      []string{"ec2", "elbv2", "rds"},
			wantNames: []string{"infra"},
		},
		{
			name:      "multiple APIs from different services",
			apis:      []string{"ec2", "iam", "s3"},
			wantNames: []string{"infra", "access", "storage"},
		},
		{
			name:      "unknown API returns nothing",
			apis:      []string{"nonexistent"},
			wantNames: nil,
		},
		{
			name:      "empty input returns nothing",
			apis:      []string{},
			wantNames: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			services := GetCloudServicesForAPIs(tt.apis...)
			gotNames := make([]string, len(services))
			for i, s := range services {
				gotNames[i] = s.Name()
			}

			if len(gotNames) != len(tt.wantNames) {
				t.Fatalf("got %d services %v, want %d services %v", len(gotNames), gotNames, len(tt.wantNames), tt.wantNames)
			}

			// Check all expected names are present (order may vary for multi-service cases)
			wantSet := make(map[string]bool)
			for _, n := range tt.wantNames {
				wantSet[n] = true
			}
			for _, n := range gotNames {
				if !wantSet[n] {
					t.Errorf("unexpected service name %q in result", n)
				}
			}
		})
	}
}

func TestGetCloudServicesForTypes(t *testing.T) {
	cleanup := setupServiceRegistry(t)
	defer cleanup()

	tests := []struct {
		name      string
		types     []string
		wantNames []string
	}{
		{
			name:      "instance maps to infra",
			types:     []string{"instance"},
			wantNames: []string{"infra"},
		},
		{
			name:      "user maps to access",
			types:     []string{"user"},
			wantNames: []string{"access"},
		},
		{
			name:      "bucket maps to storage",
			types:     []string{"bucket"},
			wantNames: []string{"storage"},
		},
		{
			name:      "multiple types from same service deduplicated",
			types:     []string{"instance", "subnet", "vpc"},
			wantNames: []string{"infra"},
		},
		{
			name:      "multiple types from different services",
			types:     []string{"instance", "user", "bucket"},
			wantNames: []string{"infra", "access", "storage"},
		},
		{
			name:      "unknown type returns nothing",
			types:     []string{"nonexistent"},
			wantNames: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			services := GetCloudServicesForTypes(tt.types...)
			gotNames := make([]string, len(services))
			for i, s := range services {
				gotNames[i] = s.Name()
			}

			if len(gotNames) != len(tt.wantNames) {
				t.Fatalf("got %d services %v, want %d services %v", len(gotNames), gotNames, len(tt.wantNames), tt.wantNames)
			}

			wantSet := make(map[string]bool)
			for _, n := range tt.wantNames {
				wantSet[n] = true
			}
			for _, n := range gotNames {
				if !wantSet[n] {
					t.Errorf("unexpected service name %q in result", n)
				}
			}
		})
	}
}

func TestResourceTypesPerServiceName(t *testing.T) {
	result := ResourceTypesPerServiceName()

	// Verify infra has expected types
	infraTypes := result["infra"]
	if len(infraTypes) == 0 {
		t.Fatal("expected infra service to have resource types")
	}

	infraSet := make(map[string]bool)
	for _, rt := range infraTypes {
		infraSet[rt] = true
	}

	expectedInfraTypes := []string{"instance", "subnet", "vpc", "keypair", "securitygroup", "volume"}
	for _, et := range expectedInfraTypes {
		if !infraSet[et] {
			t.Errorf("expected infra to contain resource type %q", et)
		}
	}

	// Verify access has expected types
	accessTypes := result["access"]
	accessSet := make(map[string]bool)
	for _, rt := range accessTypes {
		accessSet[rt] = true
	}

	expectedAccessTypes := []string{"user", "group", "role", "policy"}
	for _, et := range expectedAccessTypes {
		if !accessSet[et] {
			t.Errorf("expected access to contain resource type %q", et)
		}
	}

	// Verify storage has bucket
	storageTypes := result["storage"]
	storageSet := make(map[string]bool)
	for _, rt := range storageTypes {
		storageSet[rt] = true
	}
	if !storageSet["bucket"] {
		t.Error("expected storage to contain resource type 'bucket'")
	}
}

func TestServicePerAPIMappings(t *testing.T) {
	expectedMappings := map[string]string{
		"ec2":            "infra",
		"iam":            "access",
		"s3":             "storage",
		"sns":            "messaging",
		"route53":        "dns",
		"lambda":         "lambda",
		"cloudwatch":     "monitoring",
		"cloudfront":     "cdn",
		"cloudformation": "cloudformation",
		"sts":            "access",
		"elbv2":          "infra",
		"rds":            "infra",
	}

	for api, wantService := range expectedMappings {
		if got, ok := ServicePerAPI[api]; !ok {
			t.Errorf("ServicePerAPI missing key %q", api)
		} else if got != wantService {
			t.Errorf("ServicePerAPI[%q] = %q, want %q", api, got, wantService)
		}
	}
}

func TestServicePerResourceTypeMappings(t *testing.T) {
	expectedMappings := map[string]string{
		"instance":     "infra",
		"subnet":       "infra",
		"vpc":          "infra",
		"user":         "access",
		"group":        "access",
		"role":         "access",
		"policy":       "access",
		"bucket":       "storage",
		"topic":        "messaging",
		"zone":         "dns",
		"function":     "lambda",
		"metric":       "monitoring",
		"distribution": "cdn",
		"stack":        "cloudformation",
	}

	for resType, wantService := range expectedMappings {
		if got, ok := ServicePerResourceType[resType]; !ok {
			t.Errorf("ServicePerResourceType missing key %q", resType)
		} else if got != wantService {
			t.Errorf("ServicePerResourceType[%q] = %q, want %q", resType, got, wantService)
		}
	}
}

func TestAPIPerResourceTypeMappings(t *testing.T) {
	expectedMappings := map[string]string{
		"instance":     "ec2",
		"subnet":       "ec2",
		"loadbalancer": "elbv2",
		"database":     "rds",
		"user":         "iam",
		"bucket":       "s3",
		"zone":         "route53",
		"function":     "lambda",
		"metric":       "cloudwatch",
		"distribution": "cloudfront",
		"stack":        "cloudformation",
	}

	for resType, wantAPI := range expectedMappings {
		if got, ok := APIPerResourceType[resType]; !ok {
			t.Errorf("APIPerResourceType missing key %q", resType)
		} else if got != wantAPI {
			t.Errorf("APIPerResourceType[%q] = %q, want %q", resType, got, wantAPI)
		}
	}
}
