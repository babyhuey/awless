package awsservices

import (
	"strings"
	"testing"

	"github.com/bootswithdefer/awless/cloud"
	gen "github.com/bootswithdefer/awless/gen/aws"
	"github.com/bootswithdefer/awless/logger"
)

// Every service is generated from gen/aws/generators/services.go, so its accessors
// are identical in shape and none were exercised. A service registered with the wrong
// name, an empty resource-type list, or a resource type claimed by two services all
// break sync and `awless list` in ways nothing else catches.
// initServices populates cloud.ServiceRegistry. Init needs a region but resolves
// credentials lazily, so it succeeds without any AWS access.
func initServices(t *testing.T) {
	t.Helper()
	if len(cloud.ServiceRegistry) > 0 {
		return
	}

	// Without these the credential chain reaches for EC2 instance metadata and waits
	// out its timeout, which added 15 seconds to the suite. Static credentials keep
	// the test hermetic; no request is ever sent.
	t.Setenv("AWS_EC2_METADATA_DISABLED", "true")
	t.Setenv("AWS_ACCESS_KEY_ID", "AKIAIOSFODNN7EXAMPLE")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY")
	t.Setenv("AWS_REGION", "us-west-2")
	if err := Init("", "us-west-2", map[string]any{}, logger.DiscardLogger, func(string) error { return nil }, false); err != nil {
		t.Skipf("cannot initialize AWS services: %s", err)
	}
}

func TestRegisteredServices(t *testing.T) {
	initServices(t)
	if len(cloud.ServiceRegistry) == 0 {
		t.Fatal("no services registered; init should have populated the registry")
	}

	seenTypes := make(map[string]string)

	for name, svc := range cloud.ServiceRegistry {
		t.Run(name, func(t *testing.T) {
			// The registry is keyed by Name(), so a mismatch means lookups by name
			// silently miss.
			if got := svc.Name(); got != name {
				t.Errorf("registered under %q but Name() reports %q", name, got)
			}
			if strings.TrimSpace(name) == "" {
				t.Error("service registered under an empty name")
			}

			types := svc.ResourceTypes()
			if len(types) == 0 {
				t.Error("no resource types; the service can never be synced or listed")
			}

			for _, rt := range types {
				if strings.TrimSpace(rt) == "" {
					t.Errorf("empty resource type among %v", types)
					continue
				}
				// Two services claiming one type makes which one handles it depend on
				// map iteration order.
				if other, dup := seenTypes[rt]; dup && other != name {
					t.Errorf("resource type %q is claimed by both %s and %s", rt, other, name)
				}
				seenTypes[rt] = name
			}

			// These are read on every command through the hooks, so they must not panic
			// on a service built without credentials.
			_ = svc.Region()
			_ = svc.Profile()
			_ = svc.IsSyncDisabled()
		})
	}
}

// Each resource type must be reachable from exactly one service, since that is how
// `awless list <type>` and sync decide who fetches it.
func TestEveryResourceTypeIsServedOnce(t *testing.T) {
	initServices(t)
	owners := make(map[string][]string)
	for name, svc := range cloud.ServiceRegistry {
		for _, rt := range svc.ResourceTypes() {
			owners[rt] = append(owners[rt], name)
		}
	}

	if len(owners) == 0 {
		t.Fatal("no resource types across any service")
	}

	for rt, svcs := range owners {
		if len(svcs) > 1 {
			t.Errorf("resource type %q served by %v", rt, svcs)
		}
	}
}

// Every service must stay registered. A service can be generated and compiled while
// never being added to the registry in Init, in which case it is invisible at runtime
// — nothing lists or syncs it. This is the list as registered today; ECR, ECS, ACM
// and Application Auto Scaling resources are deliberately absent because they are
// served by infra and monitoring rather than by services of their own.
func TestNewServicesAreRegistered(t *testing.T) {
	initServices(t)
	for _, want := range []string{
		"infra", "access", "storage", "messaging", "dns", "lambda", "monitoring",
		"cdn", "cloudformation", "eks", "dynamodb", "secretsmanager", "apigateway",
		"ssm", "efs", "cloudtrail", "cloudwatchlogs", "elasticache", "eventbridge", "stepfunctions", "waf", "configservice", "kinesis", "redshift", "codepipeline", "codebuild", "beanstalk", "codedeploy", "glue",
	} {
		if _, ok := cloud.ServiceRegistry[want]; !ok {
			t.Errorf("service %q is not registered", want)
		}
	}
}

// The list above is hand-maintained, so on its own it cannot catch a service that is
// added to the generator but never registered — the very mistake the comment warns
// about. This derives the expectation from the generator definitions instead, so a new
// fetchersDef that nobody wires into Init fails here rather than going unnoticed.
func TestEveryGeneratedServiceIsRegistered(t *testing.T) {
	initServices(t)
	for _, def := range gen.FetchersDefs {
		if _, ok := cloud.ServiceRegistry[def.Name]; !ok {
			t.Errorf("service %q has a fetchers definition but is not in the registry; "+
				"add it to Init in aws/services/init.go", def.Name)
		}
	}
	// And the converse, so a registry entry cannot outlive its definition.
	defined := make(map[string]bool, len(gen.FetchersDefs))
	for _, def := range gen.FetchersDefs {
		defined[def.Name] = true
	}
	for name := range cloud.ServiceRegistry {
		if !defined[name] {
			t.Errorf("service %q is registered but has no fetchers definition", name)
		}
	}
}
