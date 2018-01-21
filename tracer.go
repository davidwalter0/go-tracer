// Trace with call context info using defer to execute auto scoped
// entry exit for block level tracing

// Chained calls enabled for many receivers by returning the Tracer
// object while setting a state in that scope

// defer tracer.Detailed(true).On().ScopedTrace()()

package tracer

import (
	"fmt"
	"runtime"
	"strings"
)

// Tracer state info for trace
type Tracer struct {
	on     bool
	depth  int
	detail bool
	mutex  *Mutex
}

// New initialize and return a new Tracer struct
func New() (tracer *Tracer) {
	tracer = &Tracer{on: true, depth: 0, detail: false, mutex: new(Mutex)}
	return tracer
}

// Reset a Tracer struct to disable tracing and depth to zero and no
// detail
func (tracer *Tracer) Reset() *Tracer {
	tracer.on = false
	tracer.depth = 0
	tracer.detail = false
	return tracer
}

// Disable tracing for this Tracer object
func (tracer *Tracer) Disable() *Tracer {
	tracer.on = false
	return tracer
}

// Enabled returns tracing enabled state
func (tracer *Tracer) Enabled() bool {
	return tracer.on
}

// Detail returns tracing detail state
func (tracer *Tracer) Detail() bool {
	return tracer.detail
}

// Enable sets to arg enable state and returns *Tracer enabled state
func (tracer *Tracer) Enable(enable bool) *Tracer {
	tracer.on = enable
	return tracer
}

// On turns on and returns *Tracer
func (tracer *Tracer) On() *Tracer {
	tracer.on = true
	return tracer
}

// CurrentScopeTraceDetail dump to stdout the stack location
func CurrentScopeTraceDetail() {
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	for i := 0; i < 10; i++ {
		if pc[i] == 0 {
			break
		}
		f := runtime.FuncForPC(pc[i])
		file, line := f.FileLine(pc[i])
		fmt.Printf("%s:%d:%s\n", file, line, f.Name())
	}
}

// Off turns off and returns *Tracer
func (tracer *Tracer) Off() *Tracer {
	tracer.on = false
	return tracer
}

// Detailed set detail state to the set value return Tracer
func (tracer *Tracer) Detailed(set bool) *Tracer {
	tracer.detail = set
	return tracer
}

// Space create a space filler string
func (tracer *Tracer) Space() (spaces string) {
	if tracer.on {
		if tracer.depth > 0 {
			spaces = fmt.Sprintf("%*s", 2*tracer.depth, " ")
		}
	}
	return
}

// Printf print formatted string if tracer is on (true)
func (tracer *Tracer) Printf(format string, args ...interface{}) {
	if tracer.on {
		fmt.Printf(format, args...)
	}
}

// Println print depth spaces before a char c (rune)
func (tracer *Tracer) Println(c rune, a ...interface{}) {
	if tracer.on {
		fmt.Printf("%s%c ", tracer.Space(), c)
		fmt.Println(a...)
	}
}

// CallerInfo build the process trace call tree context info for a
// trace.
func CallerInfo(detailed bool) (where string) {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	n := 0
	f := runtime.FuncForPC(pc[n])
	file, line := f.FileLine(pc[n])
	if pc[n] != 0 {
		path := strings.Split(f.Name(), "/")
		if detailed {
			where = fmt.Sprintf("%s:%d:%s", file, line, path[len(path)-1])
		} else {
			where = fmt.Sprintf("%s", path[len(path)-1])
		}
	} else {
		where = "unknown"
	}
	return
}

// prepend internal utility like append but before
func prepend(e interface{}, elems ...interface{}) []interface{} {
	return append([]interface{}{e}, elems...)
}

// GuardedTrace execute print calls in lock guarded blocks, enter a
// scope with trace info and return a function that prints the end of
// trace info, call with:
// defer tracer.GuardedTrace()()
// to cause the initial call and deferred calls to be handled
// transparently.
func (tracer *Tracer) GuardedTrace(a ...interface{}) (exitScopeTrace func()) {
	if tracer.on {
		where := CallerInfo(tracer.detail)
		var args []interface{}

		exitScopeTrace = func() {
			// tracer.mutex.Lock()
			tracer.depth--
			if tracer.detail {
				tracer.Printf("%s ", where)
				tracer.Println('<', a...)
			} else {
				tracer.Println('<', args...)
			}
			// tracer.mutex.Unlock()
		}

		defer tracer.mutex.Monitor()()
		// tracer.mutex.Lock()
		if tracer.detail {
			tracer.Printf("%s ", where)
			tracer.Println('>', a...)
		} else {
			args = prepend(where, a...)
			tracer.Println('>', args...)
		}
		tracer.depth++
		// tracer.mutex.Unlock()
	} else {
		exitScopeTrace = func() {}
	}
	return
}

// ScopedTrace enter a scope with trace info and return a function
// that prints the end of trace info, call with
// defer tracer.ScopedTrace()()
// to cause the initial call and deferred calls to be handled
// transparently.
func (tracer *Tracer) ScopedTrace(a ...interface{}) (exitScopeTrace func()) {
	if tracer.on {
		where := CallerInfo(tracer.detail)
		var args []interface{}

		exitScopeTrace = func() {
			var where = CallerInfo(tracer.detail)
			tracer.depth--
			if tracer.detail {
				tracer.Printf("%s ", where)
				tracer.Println('<', a...)
			} else {
				args = prepend(where, a...)
				tracer.Println('<', args...)
			}
		}

		if tracer.detail {
			tracer.Printf("%s ", where)
			tracer.Println('>', a...)
		} else {
			args = prepend(where, a...)
			tracer.Println('>', args...)
		}
		tracer.depth++
	} else {
		exitScopeTrace = func() {}
	}
	return
}
