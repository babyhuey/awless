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
	"errors"
	"strings"
)

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

// ErrExitZero reports that a command has already told the user everything they
// need and should end successfully without printing an error.
//
// Exists so those paths can return through cobra like any other outcome instead of
// calling os.Exit(0) mid-call, which skipped deferred cleanup and bypassed the
// single exit point in main.
var ErrExitZero = errors.New("awless: nothing further to report")

// ErrReported marks a failure whose explanation has already been printed, usually
// through logger.Errorf together with a suggested command.
//
// Returning it exits non-zero without printing a second message, which is what
// these paths used os.Exit(1) for. Keeping them as errors means deferred cleanup
// runs and there is still only one exit point.
var ErrReported = errors.New("awless: command failed")
