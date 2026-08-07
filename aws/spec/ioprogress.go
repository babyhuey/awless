package awsspec

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// Minimal inlining of the parts of github.com/mitchellh/ioprogress that awless
// uses: a Reader wrapper that periodically reports progress, terminal drawing,
// and byte formatting.
//
// The upstream package is a ~100 LOC orphan with no tagged release and no
// commits since February 2018, so vendoring this much of it is cheaper than
// depending on it.

// progressDrawFunc renders progress. Returning an empty string draws nothing.
type progressDrawFunc func(progress, total int64) string

// drawTerminalf writes progress to w, overwriting the current line each time so
// output stays on one line, and clears it when the transfer completes.
func drawTerminalf(w io.Writer, draw progressDrawFunc) progressDrawFunc {
	var maxLen int
	return func(progress, total int64) string {
		line := draw(progress, total)
		if line == "" {
			return ""
		}
		if len(line) > maxLen {
			maxLen = len(line)
		}
		// \r returns to line start; pad to the widest line drawn so far so
		// leftovers from a longer previous line are overwritten.
		fmt.Fprintf(w, "\r%-*s", maxLen, line)
		if progress >= total {
			fmt.Fprint(w, "\n")
		}
		return line
	}
}

const (
	byteKB = 1024
	byteMB = byteKB * 1024
	byteGB = byteMB * 1024
)

// drawTextFormatBytes renders progress as human-readable byte counts.
func drawTextFormatBytes(progress, total int64) string {
	return fmt.Sprintf("%s/%s", humanBytes(progress), humanBytes(total))
}

func humanBytes(n int64) string {
	switch {
	case n >= byteGB:
		return fmt.Sprintf("%.2f GB", float64(n)/float64(byteGB))
	case n >= byteMB:
		return fmt.Sprintf("%.2f MB", float64(n)/float64(byteMB))
	case n >= byteKB:
		return fmt.Sprintf("%.2f KB", float64(n)/float64(byteKB))
	default:
		return fmt.Sprintf("%d B", n)
	}
}

// progressReader wraps a Reader and invokes DrawFunc at DrawInterval as data is
// read. A zero DrawInterval defaults to one second.
type progressReader struct {
	Reader       io.Reader
	Size         int64
	DrawFunc     progressDrawFunc
	DrawInterval time.Duration

	mu       sync.Mutex
	progress int64
	once     sync.Once
	stopCh   chan struct{}
}

func (r *progressReader) Read(p []byte) (int, error) {
	r.once.Do(r.start)

	n, err := r.Reader.Read(p)
	if n > 0 {
		r.mu.Lock()
		r.progress += int64(n)
		r.mu.Unlock()
	}
	if err == io.EOF {
		r.stop()
	}
	return n, err
}

func (r *progressReader) start() {
	if r.DrawFunc == nil {
		return
	}
	interval := r.DrawInterval
	if interval == 0 {
		interval = time.Second
	}
	r.stopCh = make(chan struct{})

	go func() {
		t := time.NewTicker(interval)
		defer t.Stop()
		for {
			select {
			case <-r.stopCh:
				r.draw()
				return
			case <-t.C:
				r.draw()
			}
		}
	}()
}

func (r *progressReader) draw() {
	r.mu.Lock()
	progress := r.progress
	r.mu.Unlock()
	r.DrawFunc(progress, r.Size)
}

func (r *progressReader) stop() {
	if r.stopCh == nil {
		return
	}
	select {
	case <-r.stopCh: // already closed
	default:
		close(r.stopCh)
	}
}
