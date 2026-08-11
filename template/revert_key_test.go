package template

import (
	"strings"
	"testing"
)

// The revert builder emits `id=<result>` by default. Eight entities have a delete
// command that accepts only `name`, so their reverts were generated invalid and failed
// at validation with "unexpected param id" — the failure only shows up when a revert is
// actually run, which is why it went unnoticed across several service additions.
//
// This pins the key each of them reverts on. It is a table rather than one test per
// entity so that adding a service with a name-keyed delete has an obvious place to
// register.
func TestRevertUsesTheKeyDeleteAccepts(t *testing.T) {
	for _, tc := range []struct {
		name     string
		template string
		result   string
		want     []string
		notWant  []string
	}{
		{
			name:     "loggroup",
			template: `create loggroup name=/aws/lambda/fn`,
			result:   "/aws/lambda/fn",
			want:     []string{"delete loggroup", "name=/aws/lambda/fn"},
			notWant:  []string{"id="},
		},
		{
			name:     "trail",
			template: `create trail name=org-audit bucket=audit-logs`,
			result:   "org-audit",
			want:     []string{"delete trail", "name=org-audit"},
			notWant:  []string{"id="},
		},
		{
			name:     "ssmparameter",
			template: `create ssmparameter name=/app/db/host value=db.internal type=String`,
			result:   "/app/db/host",
			want:     []string{"delete ssmparameter", "name=/app/db/host"},
			notWant:  []string{"id="},
		},
		{
			name:     "dynamodbtable",
			template: `create dynamodbtable name=users partition-key=id partition-type=S`,
			result:   "users",
			want:     []string{"delete dynamodbtable", "name=users"},
			notWant:  []string{"id="},
		},
		{
			name:     "ekscluster",
			template: `create ekscluster name=prod role=arn:aws:iam::1:role/eks subnets=subnet-1,subnet-2`,
			result:   "prod",
			want:     []string{"delete ekscluster", "name=prod"},
			notWant:  []string{"id="},
		},
		{
			name:     "cachesubnetgroup",
			template: `create cachesubnetgroup name=cache-private subnets=subnet-1,subnet-2`,
			result:   "cache-private",
			want:     []string{"delete cachesubnetgroup", "name=cache-private"},
			notWant:  []string{"id="},
		},
		{
			// Compound key: the delete needs the parent cluster as well as the name.
			name:     "eksnodegroup",
			template: `create eksnodegroup name=workers cluster=prod role=arn:aws:iam::1:role/ng subnets=subnet-1`,
			result:   "workers",
			want:     []string{"delete eksnodegroup", "name=workers", "cluster=prod"},
			notWant:  []string{"id="},
		},
		{
			name:     "apigatewaystage",
			template: `create apigatewaystage api=api-1234 name=prod`,
			result:   "prod",
			want:     []string{"delete apigatewaystage", "name=prod", "api=api-1234"},
			notWant:  []string{"id="},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tpl, err := Parse(tc.template)
			if err != nil {
				t.Fatalf("parsing %q: %s", tc.template, err)
			}
			// Revert only considers commands that recorded a result.
			for _, cmd := range tpl.CommandNodesIterator() {
				cmd.CmdResult = tc.result
			}

			reverted, err := tpl.Revert()
			if err != nil {
				t.Fatalf("reverting: %s", err)
			}
			got := reverted.String()
			for _, w := range tc.want {
				if !strings.Contains(got, w) {
					t.Errorf("revert %q does not contain %q", got, w)
				}
			}
			for _, nw := range tc.notWant {
				if strings.Contains(got, nw) {
					t.Errorf("revert %q still contains %q", got, nw)
				}
			}
		})
	}
}
