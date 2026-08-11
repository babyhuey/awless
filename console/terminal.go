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

package console

import (
	"fmt"
	"io"
	"os"
	"os/signal"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

// termGetSize looks up the terminal dimensions. Replaceable in tests, where stdout is
// a pipe rather than a tty and the real call fails.
var termGetSize = func() (width, height int, err error) {
	return term.GetSize(int(os.Stdout.Fd()))
}

func GetTerminalWidth() int {
	w, _, err := termGetSize()
	if err != nil {
		return 0
	}
	return w
}

func GetTerminalHeight() int {
	_, h, err := termGetSize()
	if err != nil {
		return 0
	}
	return h
}

// defaultPtyDimension is used when the real terminal size is unavailable, which is the
// case whenever stdout is redirected.
const defaultPtyDimension = 100

// ptySize resolves the dimensions to request for a remote pty, substituting a default
// for either axis the terminal could not report. Split out from InteractiveTerminal so
// the fallback is reachable without a live SSH session; a zero slipping through here
// would ask the remote for a 0x0 terminal.
func ptySize() (width, height int) {
	width = GetTerminalWidth()
	if width == 0 {
		width = defaultPtyDimension
	}
	height = GetTerminalHeight()
	if height == 0 {
		height = defaultPtyDimension
	}
	return width, height
}

func InteractiveTerminal(client *ssh.Client) error {
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}
	go func() { _, _ = io.Copy(stdin, os.Stdin) }()

	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	go func() { _, _ = io.Copy(os.Stdout, stdout) }()

	stderr, err := session.StderrPipe()
	if err != nil {
		return err
	}
	go func() { _, _ = io.Copy(os.Stderr, stderr) }()

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	width, height := ptySize()
	if err := session.RequestPty("xterm", height, width, modes); err != nil {
		return err
	}

	// Start remote shell
	if err := session.Shell(); err != nil {
		return err
	}

	signalc := make(chan os.Signal, 1)
	defer func() {
		signal.Reset()
		close(signalc)
	}()
	go propagateSignals(signalc, session, stdin)
	signal.Notify(signalc, os.Interrupt)
	return session.Wait()
}

func propagateSignals(signalc chan os.Signal, _ *ssh.Session, stdin io.WriteCloser) {
	for s := range signalc {
		switch s {
		case os.Interrupt:
			fmt.Fprint(stdin, "\x03")
		}
	}
}
