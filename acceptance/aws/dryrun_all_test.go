package awsat

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"

	awsspec "github.com/bootswithdefer/awless/aws/spec"
	"github.com/bootswithdefer/awless/template/params"
)

// Drives every registered command through its generated dryRun path.
//
// dryRun is a per-command generated function that injects the params, sets DryRun on
// the AWS input and expects the API to reject the call — so nothing about it is
// exercised by the ordinary run path. A command whose param injection panics or whose
// awsName tag does not match its input struct fails here, which is how the
// awsstringslice and awsCall bugs would have been caught much earlier.
//
// Params are built from each command's own ParamsSpec pattern, so the combination is
// one the validator accepts.
func TestEveryCommandDryRuns(t *testing.T) {
	names := make([]string, 0, len(awsspec.AWSTemplatesDefinitions))
	for name := range awsspec.AWSTemplatesDefinitions {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		def := awsspec.AWSTemplatesDefinitions[name]

		t.Run(name, func(t *testing.T) {
			if reason, skip := skipDryRun[name]; skip {
				t.Skip(reason)
			}
			// check commands poll AWS for a state rather than mutating anything, so
			// they have no dryRun path and a dry-run mock is the wrong harness.
			if def.Action == "check" {
				t.Skip("check commands poll for a state; no dryRun path to exercise")
			}

			line, ok := templateLineFor(def)
			if !ok {
				t.Skipf("no sample value for one of %s's params", name)
			}

			// A dry run that reaches AWS is a success; anything before that — param
			// validation, struct injection, an awsName that does not exist on the
			// input — is a failure worth reporting.
			if err := Template(line).Mock(NewMock()).DryRun().RunExpectingError(t); err != nil {
				t.Errorf("%s\n  template: %s\n  pattern:  %s", err, line, def.Params)
			}
		})
	}
}

// templateLineFor builds a runnable one-liner for a command from its params pattern.
func templateLineFor(def awsspec.Definition) (string, bool) {
	required, ok := requiredFromPattern(fmt.Sprint(def.Params))
	if !ok {
		return "", false
	}

	pairs := make([]string, 0, len(required))
	for _, p := range required {
		v, known := sampleValues[p]
		if !known {
			v, known = entitySample(p, def.Entity)
		}
		if !known {
			return "", false
		}
		pairs = append(pairs, p+"="+v)
	}

	line := def.Action + " " + def.Entity
	if len(pairs) > 0 {
		line += " " + strings.Join(pairs, " ")
	}
	return line, true
}

// Commands whose dry run needs setup beyond a params line. Listed with the reason
// rather than silently omitted, so the gap is visible.
var skipDryRun = map[string]string{
	"createkeypair":             "writes a private key to __AWLESS_KEYS_DIR; covered by TestCreateKeypair",
	"createlaunchconfiguration": "resolves a distro through DescribeImages, which needs its own mock",
	"importimage":               "resolves a snapshot or bucket first, which needs its own mock",
	"deletecontainertask":       "lists task definitions first, which needs its own mock",
	"stopcontainertask":         "requires either a run-arn or a deployment-name, neither derivable from the pattern",
	"detachnetworkinterface":    "needs an existing attachment id discovered through DescribeNetworkInterfaces",
	"createapigatewayroute":     "route-key contains a space, which the generated one-liner cannot quote",
}

var groupRe = regexp.MustCompile(`\(([^()]*)\)`)

// requiredFromPattern extracts one accepted combination of required params, taking
// the first branch of any alternative so the result is a combination the validator
// accepts.
func requiredFromPattern(pattern string) ([]string, bool) {
	pattern = regexp.MustCompile(`\[[^\]]*\]`).ReplaceAllString(pattern, "")

	for strings.Contains(pattern, "(") {
		m := groupRe.FindStringSubmatchIndex(pattern)
		if m == nil {
			return nil, false
		}
		branch := strings.TrimSpace(strings.Split(pattern[m[2]:m[3]], "|")[0])
		pattern = pattern[:m[0]] + branch + pattern[m[1]:]
	}

	var out []string
	for _, p := range strings.Split(pattern, "+") {
		p = strings.TrimSpace(p)
		if p != "" && p != "none" {
			out = append(out, p)
		}
	}
	return out, true
}

// entitySample covers params whose sensible value depends on the entity.
func entitySample(param, entity string) (string, bool) {
	switch param {
	case "id", "ids":
		return "my-" + entity, true
	case "name", "names":
		return "my-" + entity, true
	case "type":
		if v, ok := typeSamples[entity]; ok {
			return v, true
		}
		return "", false
	}
	return "", false
}

var typeSamples = map[string]string{
	"instance":            "t3.micro",
	"record":              "A",
	"targetgroup":         "instance",
	"launchconfiguration": "t3.micro",
	"containertask":       "task",
	"loadbalancer":        "application",
	"ssmparameter":        "String",
	"dynamodbtable":       "S",
	"listener":            "HTTP",
	"database":            "db.t3.micro",
}

// sampleValues holds a plausible value per param name, shaped like real AWS input so
// that validators which check CIDRs, ARNs and enums are satisfied.
var sampleValues = map[string]string{
	"access":              "AKIAIOSFODNN7EXAMPLE",
	"account":             "123456789012",
	"acl":                 "private",
	"action":              "s3:GetObject",
	"action-arn":          "arn:aws:sns:us-west-2:123456789012:alerts",
	"actiontype":          "forward",
	"adjustment-scaling":  "1",
	"adjustment-type":     "ChangeInCapacity",
	"api":                 "abc123",
	"arn":                 "arn:aws:iam::123456789012:policy/my-policy",
	"association":         "eipassoc-0123456789abcdef0",
	"attachment":          "eni-attach-0123456789abcdef0",
	"availabilityzone":    "us-west-2a",
	"bucket":              "my-bucket",
	"callerreference":     "ref-2024-01-01",
	"cidr":                "10.0.0.0/24",
	"cluster":             "my-cluster",
	"container-name":      "web",
	"desired-count":       "2",
	"device":              "/dev/sdh",
	"device-index":        "1",
	"dimension":           "ecs:service:DesiredCount",
	"distro":              "amazonlinux",
	"domains":             "example.com",
	"effect":              "Allow",
	"elasticip-id":        "eipalloc-0123456789abcdef0",
	"endpoint":            "ops@example.com",
	"engine":              "postgres",
	"file":                "/dev/null",
	"gateway":             "igw-0123456789abcdef0",
	"group":               "my-group",
	"handler":             "index.handler",
	"image":               "ami-0123456789abcdef0",
	"instance":            "i-0123456789abcdef0",
	"instanceprofile":     "my-instance-profile",
	"ip":                  "52.10.20.30",
	"key":                 "Env",
	"kms-key":             "arn:aws:kms:us-west-2:123456789012:key/abcd",
	"launchconfiguration": "my-launchconfig",
	"loadbalancer":        "arn:aws:elasticloadbalancing:us-west-2:123456789012:loadbalancer/app/my-lb/abcd",
	"max-size":            "4",
	"memory-hard-limit":   "512",
	"mfa-code-1":          "123456",
	"mfa-code-2":          "654321",
	"min-size":            "1",
	"partition-key":       "id",
	"password":            "S3curePassw0rd",
	"port":                "80",
	"principal-service":   "ec2.amazonaws.com",
	"protocol":            "HTTP",
	"resource":            "arn:aws:s3:::my-bucket/*",
	"role":                "arn:aws:iam::123456789012:role/my-role",
	"route-key":           "GET /items",
	"runtime":             "python3.12",
	"s3object":            "images/disk.vmdk",
	"scalinggroup":        "my-scalinggroup",
	"secret":              "s3cret",
	"securitygroups":      "sg-0123456789abcdef0",
	"service":             "s3",
	"service-namespace":   "ecs",
	"size":                "20",
	"snapshot":            "snap-0123456789abcdef0",
	"state":               "available",
	"subnet":              "subnet-0123456789abcdef0",
	"subnets":             "subnet-0123456789abcdef0,subnet-0fedcba9876543210",
	"table":               "rtb-0123456789abcdef0",
	"targetgroup":         "arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/my-tg/abcd",
	"template-file":       "/dev/null",
	"timeout":             "30",
	"topic":               "arn:aws:sns:us-west-2:123456789012:alerts",
	"ttl":                 "300",
	"url":                 "https://sqs.us-west-2.amazonaws.com/123456789012/my-queue",
	"user":                "my-user",
	"username":            "my-user",
	"value":               "production",
	"values":              "1.2.3.4",
	"version":             "1",
	"volume":              "vol-0123456789abcdef0",
	"vpc":                 "vpc-0123456789abcdef0",
	"zone":                "/hostedzone/Z3P5QSUBK4POTI",
}

// Guards the helper above: a pattern with alternatives must yield one branch, not all.
func TestRequiredFromPattern(t *testing.T) {
	for _, tc := range []struct {
		pattern string
		want    []string
	}{
		{"cidr + [name]", []string{"cidr"}},
		{"(ids | id)", []string{"ids"}},
		{"name + ttl + type + (values | value) + zone + [comment]", []string{"name", "ttl", "type", "values", "zone"}},
		{"(user | role | group) + (arn | access + service)", []string{"user", "arn"}},
		{"none", nil},
	} {
		got, ok := requiredFromPattern(tc.pattern)
		if !ok {
			t.Errorf("%q: not parsed", tc.pattern)
			continue
		}
		if strings.Join(got, ",") != strings.Join(tc.want, ",") {
			t.Errorf("%q: got %v, want %v", tc.pattern, got, tc.want)
		}
	}
}

var _ = params.List // keep the import if the helper changes shape
