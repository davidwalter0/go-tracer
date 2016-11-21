package tracer

import (
	"fmt"
	"testing"
)

var tracer *Tracer = New()

// use deeper call to demo recursive calls
func deeper(depth int, n int) {
	defer tracer.ScopedTrace(fmt.Sprintf("depth %4d %4d", depth, n))()
	if n > 0 {
		deeper(depth, n-1)
	}
}

func recursive_trace(n int) {
	defer tracer.ScopedTrace(fmt.Sprintf(">>%d<<", n))()
	if n > 0 {
		if n%2 == 1 {
			recursive_trace(n - 1)
			deeper(tracer.depth, n+1)
		} else {
			recursive_trace(n - 2)
		}
	}
}

func TestTracerRecurse(t *testing.T) {
	fmt.Println()
	defer tracer.ScopedTrace()()
	if tracer != nil {
		recursive_trace(3)
	}
}

func TestTracerRecurseDetail(t *testing.T) {
	tracer.Detailed(true)
	fmt.Println()
	defer tracer.ScopedTrace()()
	if tracer != nil {
		recursive_trace(3)
	}
}

func TestTracerOff(t *testing.T) {
	fmt.Println()
	if tracer != nil {
		tracer.Off()
		defer tracer.ScopedTrace()()
	}
}

func TestTracerOn(t *testing.T) {
	fmt.Println()
	tracer.On()
	defer tracer.ScopedTrace()()
	deeper(tracer.depth, 0)
}

func TestTracerDetailed(t *testing.T) {
	fmt.Println()
	tracer.Detailed(true)
	tracer.On()
	defer tracer.ScopedTrace()()
	deeper(tracer.depth, 0)
}

func TestTracerOnOffTracer(t *testing.T) {
	fmt.Println()
	tracer.Reset()
	tracer.On()
	{
		defer tracer.ScopedTrace("scoped", "in", "braces")()
		tracer.Reset()
	}
	defer tracer.ScopedTrace("scoped", "by", "func()")()
	tracer.Reset()
	tracer.Off()
	{
		defer tracer.ScopedTrace("scoped", "in", "braces")()
		tracer.Reset()
	}
	tracer.Off()
}

func TestTracerOnOffOnTracer(t *testing.T) {
	fmt.Println()
	tracer.Reset()
	tracer.On()
	defer tracer.ScopedTrace("scoped", "fnctn", "braces")()
	{
		defer tracer.ScopedTrace("scoped", "in", "braces")()
	}
	defer tracer.ScopedTrace("scoped", "by", "func()")()
	{
		defer tracer.ScopedTrace("scoped", "in", "braces")()
	}
}
