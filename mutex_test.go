package tracer

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

var N int = 5
var started = make(chan int, N)
var done = make(chan int, N)

func (tracer *Tracer) recurse(n int) {
	defer tracer.ScopedTrace(fmt.Sprintf(">>%d<<", n))()
	if n > 0 {
		if n%2 == 1 {
			tracer.recurse(n - 1)
		} else {
			tracer.recurse(n - 2)
		}
	}
}

func Monitor(n int) {
	defer tracer.GuardedTrace(fmt.Sprintf("Monitor%d", n))()
	started <- 0
	TimedWait()
	{
		defer tracer.mutex.Monitor()()
		tracer.recurse(n + 2)
	}
	done <- 0
}

func TimedWait() {
	t := rand.Intn(250)
	select {
	case <-time.After(time.Duration(t) * time.Millisecond):
	}
}

func ChanWait(c chan int, max int) {
	for i := 0; i < max; i++ {
		select {
		case <-c:
		}
	}
}

func WaitStart() {
	ChanWait(started, cap(started))
}

func TestMonitor0(t *testing.T) {
	rand.Seed(42)
	runtime.GOMAXPROCS(runtime.NumCPU())
	{
		for i := 0; i < N; i++ {
			go func(n int) {
				defer Monitor(n)
			}(i)
		}
		defer tracer.GuardedTrace()()
		started <- 0
		WaitStart()
		TimedWait()
		{
			defer tracer.mutex.Monitor()()
			tracer.recurse(3)
			fmt.Println("Exiting Function TestMonitor0")
			done <- 0
		}
	}
}

func TestLast(t *testing.T) {
	ChanWait(done, cap(done))
}
