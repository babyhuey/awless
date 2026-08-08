package ast

type Entity string

// IsInvalidEntity reports whether s is outside the set the template parser
// accepts. The set itself is generated into gen_entities.go from the `entity:`
// struct tags in aws/spec, so it cannot drift from the registered commands.
func IsInvalidEntity(s string) bool {
	_, ok := entities[Entity(s)]
	return !ok
}
