package ast

import "strings"

// RedactedValue replaces sensitive parameter values when a command is
// serialized for persistence.
//
// It must remain parseable by the template grammar, because persisted command
// lines are re-parsed when reading the log back (see
// template.TemplateExecution.UnmarshalJSON). '*' is in the UnquotedParam
// character class of awless-template-syntax.peg, so this round-trips.
const RedactedValue = "*****"

// sensitiveParamKeys holds template parameter names whose values must never be
// persisted or logged.
//
// Keys are matched exactly against the template parameter name, i.e. the
// `templateName:"..."` struct tag in aws/spec, not the AWS API field name.
//
// This lives here rather than in aws/spec because the template packages cannot
// import aws/spec (aws/spec depends on template). Use RegisterSensitiveParams
// to add entries from higher-level packages.
var sensitiveParamKeys = map[string]struct{}{
	// aws/spec/loginprofile.go (create + update), aws/spec/database.go (create)
	"password": {},
}

// RegisterSensitiveParams marks additional template parameter names as
// sensitive. Intended for use from package init functions.
//
// Call this when adding a command that accepts a secret as a parameter,
// otherwise the value will be written to the local template log in plaintext.
func RegisterSensitiveParams(keys ...string) {
	for _, k := range keys {
		sensitiveParamKeys[k] = struct{}{}
	}
}

// IsSensitiveParam reports whether values for the given template parameter
// name must be redacted before being persisted or logged.
func IsSensitiveParam(key string) bool {
	_, ok := sensitiveParamKeys[key]
	return ok
}

// StringRedacted renders the whole AST with sensitive parameter values
// replaced, mirroring String(). Use it when building any text that will be
// persisted — notably the template execution Message, which is stored
// alongside the command lines in the local log.
func (a *AST) StringRedacted() string {
	var all []string
	for _, stat := range a.Statements {
		if cmd, ok := stat.Node.(*CommandNode); ok {
			all = append(all, cmd.StringRedacted())
			continue
		}
		all = append(all, stat.String())
	}
	return strings.Join(all, "\n")
}
