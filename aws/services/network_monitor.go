package awsservices

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/wallix/awless/console"
)

var DefaultNetworkMonitor = &NetworkMonitor{requests: make(map[string]*req)}

type NetworkMonitor struct {
	requests map[string]*req
	l        sync.Mutex
}

type req struct {
	name    string
	from    time.Time
	to      time.Time
	retries []time.Time
}

func (n *NetworkMonitor) DisplayStats(w io.Writer) {
	fmt.Fprintf(w, "\n%d requests sent:\n", len(n.requests))

	var sorted []*req

	var min, max time.Time
	var maxFunctionNameLength int
	for _, r := range n.requests {
		if min.IsZero() || r.from.Before(min) {
			min = r.from
		}
		if max.IsZero() || r.to.After(max) {
			max = r.to
		}
		if len(r.name) > maxFunctionNameLength {
			maxFunctionNameLength = len(r.name)
		}
		sorted = append(sorted, r)
	}
	sort.Slice(sorted, func(i int, j int) bool {
		if sorted[i].from.Equal(sorted[j].from) {
			return sorted[i].to.Before(sorted[j].to)
		}
		return sorted[i].from.Before(sorted[j].from)
	})

	maxwidth := uint(console.GetTerminalWidth() - maxFunctionNameLength - 11) // 11 = '['+']'+' '+'('+4+'m'+'s'+')'
	maxDuration := max.Sub(min)

	for _, r := range sorted {
		if len(r.retries) > 0 {
			drawRequest(w, r.name, min, r.from, r.retries[0], maxwidth, maxDuration, "[", "X")
			for i := 0; i < len(r.retries)-1; i++ {
				drawRequest(w, r.name, min, r.retries[i], r.retries[i+1], maxwidth, maxDuration, "o", "X")
			}
			drawRequest(w, r.name, min, r.retries[len(r.retries)-1], r.to, maxwidth, maxDuration, "o", "]")
		} else {
			drawRequest(w, r.name, min, r.from, r.to, maxwidth, maxDuration, "[", "]")
		}
	}
}

func drawRequest(w io.Writer, name string, min, from, to time.Time, maxwidth uint, maxduration time.Duration, startChar, stopChar string) {
	duration := to.Sub(from)
	width := uint(duration) * maxwidth / uint(maxduration)
	before := uint(from.Sub(min)) * maxwidth / uint(maxduration)
	after := maxwidth - width - before
	fmt.Fprintf(w, "%s%s%s%s%s %s(%dms)\n", strings.Repeat(" ", int(before)), startChar, strings.Repeat("-", int(width)), stopChar, strings.Repeat(" ", int(after)), name, duration/(1000*1000))
}

func (n *NetworkMonitor) AddRequest(name string) {
	n.l.Lock()
	defer n.l.Unlock()
	if r, ok := n.requests[name]; ok {
		r.retries = append(r.retries, time.Now().UTC())
	} else {
		n.requests[name] = &req{name: name, from: time.Now().UTC()}
	}
}

func (n *NetworkMonitor) SetRequestEnd(name string) {
	n.l.Lock()
	defer n.l.Unlock()
	r, ok := n.requests[name]
	if !ok {
		fmt.Fprintf(os.Stderr, "request '%s' not found\n", name)
		return
	}
	r.to = time.Now().UTC()
}
