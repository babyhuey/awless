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
	"net"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bootswithdefer/awless/config"
	"github.com/bootswithdefer/awless/logger"
	"github.com/bootswithdefer/awless/web"
)

var (
	webPortFlag      string
	webListenAllFlag bool
)

func init() {
	RootCmd.AddCommand(webCmd)

	webCmd.Flags().StringVar(&webPortFlag, "port", ":8080", "Web UI port to listen on")
	webCmd.Flags().BoolVar(&webListenAllFlag, "listen-all", false, "Listen on all interfaces instead of loopback only. The web UI has no authentication: only use this on a trusted network")
}

// webListenAddr builds the listen address for the web UI.
//
// It defaults to loopback: this UI serves the local synced inventory (instance
// IDs, private IPs, security group rules, VPC topology) and has no
// authentication, so binding every interface by default exposed it to the whole
// network. listenAll is the explicit opt-out.
func webListenAddr(port string, listenAll bool) string {
	host := "127.0.0.1"
	if listenAll {
		host = "0.0.0.0"
	}
	return net.JoinHostPort(host, strings.TrimPrefix(port, ":"))
}

var webCmd = &cobra.Command{
	Use:               "web",
	Hidden:            true,
	Short:             "[Experimental] Browse your locally synced data through the web",
	PersistentPreRunE: applyHooks(initLoggerHook, initAwlessEnvHook),

	RunE: func(cmd *cobra.Command, args []string) error {
		if webListenAllFlag {
			logger.Warning("web UI listening on all interfaces with no authentication; anyone who can reach this port can browse your synced AWS inventory")
		}

		server := web.New(webListenAddr(webPortFlag, webListenAllFlag), config.GetAWSProfile())
		if err := server.Start(); err != nil {
			return err
		}
		return nil
	},
}
