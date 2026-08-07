/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func exitOn(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, color.RedString("[error]  "), err)
		os.Exit(1)
	}
}

// capitalize upper-cases the first character of s.
//
// Replaces strings.Title, deprecated in Go 1.18 because it applies Unicode word
// boundaries and title-cases every word. Every input here is a single ASCII
// token — an AWS API name, resource type, template action, or policy effect —
// so this is both correct and narrower than the deprecated behavior.
func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
