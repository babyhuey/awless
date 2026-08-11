package template

import (
	"fmt"
	"os"
	"testing"
)

func TestCorpusBaseline(t *testing.T) {
	values := []string{
		"22-80", "1-10", "0-65535", "22", "0022", "3.14", "-1",
		"12345678-1234-1234-1234-123456789012",
		"abcd-1234", "abcd1234-5678", "1-2-3", "22-80abc", "1-2.3",
		"web-server-01", "10.0.0.0/16", "arn:aws:iam::123456789012:role/Admin",
		"t2.micro", "sg-1234abcd", "vpc-0a1b2c3d", "i-0123456789abcdef0",
		"2024-02-09", "2024-02-09T10:30:00Z",
		"'22-80'", "\"22-80\"", "'12345678-1234-1234-1234-123456789012'",
		"[22-80,443]", "[1-10]", "a,b,c",
		"{myhole}", "$myref", "@myalias",
	}
	out, _ := os.Create("/tmp/baseline.txt")
	defer out.Close()
	for _, v := range values {
		tpl, err := Parse("create instance portrange=" + v)
		if err != nil {
			fmt.Fprintf(out, "%-42s REJECTED\n", v)
			continue
		}
		// Record the resulting value and its Go type, which is what a change must not alter.
		var got string
		for _, cmd := range tpl.CommandNodesIterator() {
			for k, pv := range cmd.ToDriverParams() {
				got += fmt.Sprintf("%s=%v(%T)", k, pv, pv)
			}
		}
		fmt.Fprintf(out, "%-42s OK %s\n", v, got)
	}
}
