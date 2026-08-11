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

package config

import (
	"bytes"
	"fmt"
)

const AwlessASCIILogo = `
 █████╗  ██╗    ██╗ ██╗     ██████  ██████╗ ██████╗     
██╔══██╗ ██║    ██║ ██║     ██╔══╝  ██╔═══╝ ██╔═══╝    
███████║ ██║ █╗ ██║ ██║     ████╗   ██████  ██████   
██╔══██║ ██║███╗██║ ██║     ██╔═╝       ██╗     ██╗   
██║  ██║ ╚███╔███╔╝ ██████╗ ██████╗ ██████║ ██████║   
╚═╝  ╚═╝  ╚══╝╚══╝  ╚═════╝ ╚═════╝ ╚═════╝ ╚═════╝
`

var (
	Version  = "v1.1.1"
	BuildFor string

	buildSha, buildDate, buildArch, buildOS string
)

type BuildInfo struct {
	Version, Sha, Date, Arch, OS, For string
}

func (b BuildInfo) String() string {
	var buff bytes.Buffer
	fmt.Fprintf(&buff, "version=%s", b.Version)

	if b.Sha != "" {
		fmt.Fprintf(&buff, ", commit=%s", b.Sha)
	}
	if b.Date != "" {
		fmt.Fprintf(&buff, ", build-date=%s", b.Date)
	}
	if b.Arch != "" {
		fmt.Fprintf(&buff, ", build-arch=%s", b.Arch)
	}
	if b.OS != "" {
		fmt.Fprintf(&buff, ", build-os=%s", b.OS)
	}
	if b.For != "" {
		fmt.Fprintf(&buff, ", build-for=%s", b.For)
	}
	return buff.String()
}

var CurrentBuildInfo = BuildInfo{
	Version: Version,
	For:     BuildFor,
	Sha:     buildSha,
	Date:    buildDate,
	Arch:    buildArch,
	OS:      buildOS,
}
