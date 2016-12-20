package tracer

import (
	"fmt"
	"runtime"
	"strings"
)

type Tracer struct {
	on     bool
	depth  int
	detail bool
	mutex  *Mutex
}

func New() (tracer *Tracer) {
	tracer = &Tracer{on: true, depth: 0, detail: false, mutex: new(Mutex)}
	return tracer
}

func (tracer *Tracer) Reset() *Tracer {
	tracer.on = false
	tracer.depth = 0
	tracer.detail = false
	return tracer
}

func (tracer *Tracer) Disable() *Tracer {
	tracer.on = false
	return tracer
}

func (tracer *Tracer) Enable(enable bool) *Tracer {
	tracer.on = enable
	return tracer
}

func (tracer *Tracer) On() *Tracer {
	tracer.on = true
	return tracer
}

func CurrentScopeTraceDetail() {
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	for i := 0; i < 10; i++ {
		if pc[i] == 0 {
			break
		}
		f := runtime.FuncForPC(pc[i])
		file, line := f.FileLine(pc[i])
		fmt.Printf("%s:%d %s\n", file, line, f.Name())
	}
}

func (tracer *Tracer) Off() *Tracer {
	tracer.on = false
	return tracer
}

func (tracer *Tracer) Detailed(set bool) *Tracer {
	tracer.detail = set
	return tracer
}

func (tracer *Tracer) Space() (spaces string) {
	if tracer.on {
		if tracer.depth > 0 {
			spaces = fmt.Sprintf("%*s", 2*tracer.depth, " ")
		}
	}
	return
}

func (tracer *Tracer) Printf(format string, args ...interface{}) {
	if tracer.on {
		fmt.Printf(format, args...)
	}
}

func (tracer *Tracer) Println(c rune, a ...interface{}) {
	if tracer.on {
		fmt.Printf("%s%c ", tracer.Space(), c)
		fmt.Println(a...)
	}
}

func CallerInfo(detailed bool) (where string) {
	pc := make([]uintptr, 1)
	runtime.Callers(3, pc)
	n := 0
	f := runtime.FuncForPC(pc[n])
	file, line := f.FileLine(pc[n])
	if pc[n] != 0 {
		path := strings.Split(f.Name(), "/")
		if detailed {
			where = fmt.Sprintf("%s:%d %s", file, line, path[len(path)-1])
		} else {
			where = fmt.Sprintf("%s", path[len(path)-1])
		}
	} else {
		where = "unknown"
	}
	return
}

func prepend(e interface{}, elems ...interface{}) []interface{} {
	return append([]interface{}{e}, elems...)
}

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

func (tracer *Tracer) ScopedTrace(a ...interface{}) (exitScopeTrace func()) {
	if tracer.on {
		where := CallerInfo(tracer.detail)
		var args []interface{}

		exitScopeTrace = func() {
			tracer.depth--
			if tracer.detail {
				tracer.Printf("%s ", where)
				tracer.Println('<', a...)
			} else {
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
