package tracer

import (
	"fmt"
	"sync"
)

type Mutex sync.Mutex

// Lock the mutex
func (mutex *Mutex) Lock() {
	(*sync.Mutex)(mutex).Lock()
}

// Unlock the mutex
func (mutex *Mutex) Unlock() {
	(*sync.Mutex)(mutex).Unlock()
}

// MonitorTrace:
// block scoped mutex with depth print

// defer mutex.MonitorTrace()()
// prefer to use example from tests with defer GuardedTrace()()
func (mutex *Mutex) MonitorTrace(args ...interface{}) func() {
	mutex.Lock()
	fmt.Printf("> ")
	fmt.Println(args...)
	return func() {
		fmt.Printf("< ")
		fmt.Println(args...)
		mutex.Unlock()
	}
}

// Monitor: block scoped mutex use
// defer mutex.Guard()()
func (mutex *Mutex) Monitor() func() {
	mutex.Lock()
	return func() {
		mutex.Unlock()
	}
}

// Guard: alias of Monitor block scoped mutex use
// defer mutex.Guard()()
func (mutex *Mutex) Guard() func() {
	return mutex.Monitor()
}
