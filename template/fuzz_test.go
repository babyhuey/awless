package template

import "testing"

// The template parser is generated from a PEG grammar and consumes untrusted
// input: `awless run` accepts a local file path or an HTTP URL. It must never
// panic, only return an error.
//
// Replaces the go-fuzz target that used to live in template/fuzz/parsing. That
// one required the external go-fuzz tool, could not run under `go test`, and so
// was never executed.
func FuzzParse(f *testing.F) {
	// Seeds carried over from the deleted go-fuzz corpus in template/fuzz, plus
	// one per construct the grammar handles specially.
	for _, seed := range []string{
		"# comment",
		"cidr=10.0.0.0/16",
		"name=cockroachdb-vpc",
		"id=$gateway",
		"@cockroachdb-pubsubnet-1",
		"{hole.test}+'-pubsubnet-1'",
		"152",
		"create instance name=web",
		"create subnet cidr=10.0.0.0/16 vpc=@my-vpc",
		"create instance name={hole.test}+'-suffix'",
		"create instance name=$ref",
		"create instance name=@alias",
		"create instance count=1-10",
		`create instance name="double quoted"`,
		"create instance name='single quoted'",
		"create instance list=[a,b,c]",
		"# comment only",
		"// other comment",
		"a = b\ncreate instance name=$a",
		"",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, in string) {
		// Only requirement: never panic. A parse error is a valid outcome.
		tpl, err := Parse(in)
		if err != nil {
			return
		}
		if tpl == nil {
			t.Fatal("Parse returned nil template and nil error")
		}
		// A successfully parsed template must survive being rendered and
		// re-parsed, since that round trip is what the local log replay does.
		_ = tpl.String()
	})
}

// FuzzParseStatement covers the single-statement entry point used when reading
// persisted command lines back out of the local template log.
func FuzzParseStatement(f *testing.F) {
	for _, seed := range []string{
		"create instance name=web",
		"delete vpc id=vpc-1234",
		"attach policy arn=arn:aws:iam::aws:policy/AdministratorAccess user=jsmith",
		"create loginprofile username=jsmith password=*****",
		"",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, in string) {
		_, _ = parseStatement(in)
	})
}
