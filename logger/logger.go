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

package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync/atomic"

	"github.com/fatih/color"
)

var DefaultLogger *Logger = New("", 0)
var DiscardLogger *Logger = New("", 0, io.Discard)

const (
	VerboseF = 1 << iota
	ExtraVerboseF
)

type Logger struct {
	verbose uint32 // atomic
	out     *log.Logger
	w       io.Writer
}

var (
	infoPrefix         = color.GreenString("[info]   ")
	errorPrefix        = color.RedString("[error]  ")
	warningPrefix      = color.YellowString("[warning]")
	verbosePrefix      = color.CyanString("[verbose]")
	extraVerbosePrefix = color.MagentaString("[extra]  ")
)

func New(prefix string, flag int, w ...io.Writer) *Logger {
	var out io.Writer = os.Stderr
	if len(w) > 0 {
		out = w[0]
	}
	return &Logger{out: log.New(out, prefix, flag), w: out}
}

func (l *Logger) Verbosef(format string, v ...any) {
	if l.verbosity() > 0 {
		l.out.Println(prepend(verbosePrefix, fmt.Sprintf(format, v...))...)
	}
}

func (l *Logger) Verbose(v ...any) {
	if l.verbosity() > 0 {
		l.out.Println(prepend(verbosePrefix, v...)...)
	}
}

func (l *Logger) ExtraVerbosef(format string, v ...any) {
	if l.verbosity() > 1 {
		l.out.Println(prepend(extraVerbosePrefix, fmt.Sprintf(format, v...))...)
	}
}

func (l *Logger) ExtraVerbose(v ...any) {
	if l.verbosity() > 1 {
		l.out.Println(prepend(extraVerbosePrefix, v...)...)
	}
}

func (l *Logger) Info(v ...any) {
	l.out.Println(prepend(infoPrefix, v...)...)
}

func (l *Logger) Infof(format string, v ...any) {
	l.out.Println(prepend(infoPrefix, fmt.Sprintf(format, v...))...)
}

func (l *Logger) InteractiveInfof(format string, v ...any) {
	fmt.Fprint(l.w, prepend("\r\033[K"+infoPrefix, " ", fmt.Sprintf(format, v...))...)
}

func (l *Logger) Error(v ...any) {
	l.out.Println(prepend(errorPrefix, v...)...)
}

func (l *Logger) Errorf(format string, v ...any) {
	l.out.Println(prepend(errorPrefix, fmt.Sprintf(format, v...))...)
}

func (l *Logger) MultiLineError(err error) {
	if err != nil {
		for _, msg := range formatMultiLineErrMsg(err.Error()) {
			l.out.Println(color.New(color.FgRed).Sprint(msg))
		}
	}
}

func (l *Logger) Warning(v ...any) {
	l.out.Println(prepend(warningPrefix, v...)...)
}

func (l *Logger) Warningf(format string, v ...any) {
	l.out.Println(prepend(warningPrefix, fmt.Sprintf(format, v...))...)
}

func (l *Logger) Println() {
	l.out.Println()
}

func (l *Logger) SetVerbose(level int) {
	atomic.StoreUint32(&l.verbose, uint32(level))
}

func (l *Logger) verbosity() uint32 {
	return atomic.LoadUint32(&l.verbose)
}

func Verbosef(format string, v ...any) {
	DefaultLogger.Verbosef(format, v...)
}

func Verbose(v ...any) {
	DefaultLogger.Verbose(v...)
}

func ExtraVerbosef(format string, v ...any) {
	DefaultLogger.ExtraVerbosef(format, v...)
}

func ExtraVerbose(v ...any) {
	DefaultLogger.ExtraVerbose(v...)
}

func Info(v ...any) {
	DefaultLogger.Info(v...)
}

func Infof(format string, v ...any) {
	DefaultLogger.Infof(format, v...)
}

func Error(v ...any) {
	DefaultLogger.Error(v...)
}

func Errorf(format string, v ...any) {
	DefaultLogger.Errorf(format, v...)
}

func Warning(v ...any) {
	DefaultLogger.Warning(v...)
}

func Warningf(format string, v ...any) {
	DefaultLogger.Warningf(format, v...)
}

func MultiLineError(err error) {
	DefaultLogger.MultiLineError(err)
}

func prepend(s any, v ...any) []any {
	return append([]any{s}, v...)
}

func formatMultiLineErrMsg(msg string) []string {
	notabs := strings.ReplaceAll(msg, "\t", "")
	var indented []string
	for _, line := range strings.Split(notabs, "\n") {
		indented = append(indented, fmt.Sprintf("          %s", line))
	}
	return indented
}
