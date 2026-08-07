package parsing

import (
	"fmt"

	"github.com/bootswithdefer/awless/template"
)

func Fuzz(data []byte) int {

	if _, err := template.Parse(fmt.Sprintf("none none %s", data)); err != nil {
		return 0
	}

	return 1
}
