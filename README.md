**tracer**

enable tracing of function calls with scoped entry/exit indentation

https://github.com/davidwalter0/go-tracer.git

```
go get github.com/davidwalter0/go-tracer
```

---
**Deferable resource acquisition and release**

A facet of scoped syntax is often the scoped resource acquisition,
release semantic of a file handle or file reader

One method to achieve this is to have the acquisition function return
a release function

```
    func AcquireRelease( resource *Resource ) ( release func() ) {
        resource.acquire()
        return func() { resource.release() }
    }

    // Can be called then like 
    defer AcquireRelease(&resource)()
    // which acquires at the point of the defer call, and is released
    // when the deferred returned method leaves the defer scope

```

This deferable pattern is used to setup a single scoped call from the
defer function call point in source.

Another example to acquire and release a lock

```
    import (
        "github.com/davidwalter0/go-mutex"
    )
    // package scoped var
    var monitor = mutex.NewMonitor()

    func ProtectedBlock() {
      defer monitor()()
      // ... protected shared resources
    }

```

**Tracing**

Execute a scoped print from point of defer call to exit of a block with 0..n
args to print for the trace entry and exit text.

For a thread safe version, use defer GuardedTrace()() be aware of
locking hierarchical tracing with the same guard.

chained configurable option settings added for enable and detailed e.g.

```
    var detailed=false
    var enable=false
...
    defer tracer.Detailed(detailed).Enable(enable).ScopedTrace()()

```

Create an instance of tracer and call the receiver method with a defer
call, for example:

*Create*

```
    var tracer *tracer.Tracer = tracer.New()
```

*Call*

```
    defer tracer.ScopedTrace()()
```

*Call with trace text args*
```
	defer tracer.ScopedTrace("scope process", "more", "text"))()
```

The unit tests dump some sample output
 
```
> tracer.TestTracerRecurse
  > tracer.recursive_trace >>3<<
    > tracer.recursive_trace >>2<<
      > tracer.recursive_trace >>0<<
      < tracer.recursive_trace >>0<<
    < tracer.recursive_trace >>2<<
    > tracer.deeper depth    2    4
      > tracer.deeper depth    2    3
        > tracer.deeper depth    2    2
          > tracer.deeper depth    2    1
            > tracer.deeper depth    2    0
            < tracer.deeper depth    2    0
          < tracer.deeper depth    2    1
        < tracer.deeper depth    2    2
      < tracer.deeper depth    2    3
    < tracer.deeper depth    2    4
  < tracer.recursive_trace >>3<<
< tracer.TestTracerRecurse

```


