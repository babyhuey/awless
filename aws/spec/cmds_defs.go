package awsspec

import "github.com/bootswithdefer/awless/template/params"

type Definition struct {
	Action, Entity, Api string
	Params              params.Rule
}

func AWSLookupDefinitions(key string) (t Definition, ok bool) {
	t, ok = AWSTemplatesDefinitions[key]
	return
}
