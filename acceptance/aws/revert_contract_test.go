package awsat

import (
	"strings"
	"testing"

	awsspec "github.com/bootswithdefer/awless/aws/spec"
	"github.com/bootswithdefer/awless/template"
	"github.com/bootswithdefer/awless/template/params"
)

// A create's revert must be a command its own delete accepts.
//
// template/revert.go picks the revert param from an explicit per-entity switch and falls
// back to `id=<result>`. Nothing tied that switch to what the delete command actually
// takes, so a resource keyed on `name` or `arn` got an `id=` it would reject — and the
// failure only appears when someone runs a revert. Eight resources shipped that way, and
// two more were added the same week before this test existed.
//
// This walks every create with a matching delete, generates the revert, and checks the
// resulting params against the delete's own ParamsSpec. It is the guarantee that replaces
// finding these one at a time.
func TestEveryRevertIsAcceptedByItsDelete(t *testing.T) {
	deletes := map[string]params.Rule{}
	for _, def := range awsspec.AWSTemplatesDefinitions {
		if def.Action == "delete" {
			deletes[def.Entity] = def.Params
		}
	}

	for name, def := range awsspec.AWSTemplatesDefinitions {
		if def.Action != "create" {
			continue
		}
		deleteParams, ok := deletes[def.Entity]
		if !ok {
			continue // nothing to revert to
		}

		t.Run(name, func(t *testing.T) {
			// Build a minimal create from the required params so the revert has
			// something to read. Values are placeholders; only the keys matter.
			required, _, _ := params.List(def.Params)

			// params.List flattens across branches, so for a create with mutually
			// exclusive params the set above can be self-contradictory — `create
			// record` comes back wanting both `value` and `values`. That is a limit of
			// synthesizing the command rather than a defect in the revert, so skip
			// instead of reporting it.
			if err := params.Run(def.Params, required); err != nil {
				t.Skipf("cannot synthesize a valid create for %s: %s", def.Entity, err)
			}

			var fields []string
			for _, k := range required {
				fields = append(fields, k+"=placeholder-"+k)
			}
			line := "create " + def.Entity
			if len(fields) > 0 {
				line += " " + strings.Join(fields, " ")
			}

			tpl, err := template.Parse(line)
			if err != nil {
				t.Skipf("cannot synthesize a create for %s: %s", def.Entity, err)
			}
			for _, cmd := range tpl.CommandNodesIterator() {
				cmd.CmdResult = "placeholder-result"
			}

			reverted, err := tpl.Revert()
			if err != nil {
				t.Skipf("not revertible: %s", err)
			}

			for _, cmd := range reverted.CommandNodesIterator() {
				// Only the delete of this same entity is in scope; reverts may also
				// emit `check` commands and teardown of related resources.
				if cmd.Action != "delete" || cmd.Entity != def.Entity {
					continue
				}
				var keys []string
				for k := range cmd.ParamNodes {
					keys = append(keys, k)
				}
				// Mirror compilation: a required param the revert omits becomes a
				// prompted hole, which is legitimate. An unknown param is not.
				keys = append(keys, deleteParams.Missing(keys)...)
				if err := params.Run(deleteParams, keys); err != nil {
					t.Errorf("revert emits `%s` but `delete %s` rejects it: %s\n  accepts: %s",
						cmd.String(), def.Entity, err, deleteParams)
				}
			}
		})
	}
}
