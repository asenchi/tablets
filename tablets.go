package tablets

import (
	"fmt"
	"io"
	"runtime"
	"sync"
	"time"
)

type Tablet struct {
	stream        io.Writer
	globalContext map[string]interface{}
	mu            sync.Mutex
	buf           []byte
	timeUnit      time.Millisecond
}

func NewTablet(stream io.Writer) *Tablet {
	return &Tablet{
		stream:        stream,
		globalContext: make(globalContext),
	}
}

var std = NewTablet(os.Stderr)

func (t *Tabtet) SetGlobalContext(ctx map[string]interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.globalContext = ctx
}

func (t *Tablet) GlobalContext() map[string]interface{} {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.globalContext
}

func GlobalContext() map[string]interface{} {
	return std.GlobalContext()
}

func SetStream(w io.Writer) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.stream = w
}

// Take a map and conver it to a string with "k=v"
//func parseToStr(data map[string]interface{}) out string, err error {
//}

func (t *Tablet) Log(data map[string]interface{}) error {
	var cd = 1 // Limit our call depth at runtime
	var file string
	var line int
	var ok bool

	_, file, line, ok = runtime.Caller(cd)
	if !ok {
		file = "???"
		line = 0
	}
        ctxhdr := t.globalContext
        ctxhdr["file"] = file
        ctxhdr["line"] = line

	t.mu.Lock()

        // defer needs to unlock and do the time of execution stuff.
	defer func() {
		t.mu.Unlock()
		timeExec(time.Now(), ctxhdr)
	}()
}

func (t *Tablet) timeExec(start time.Time, ctxhdr map[string]interface{}) {
	elapsed := time.Since(start)
	ms := float64(elapsed) / float64(t.timeUnit)
	ctxhdr["at"] = ms
}
