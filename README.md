---
**tracer**

enable tracing of function calls with scoped entry/exit indentation

https://github.com/davidwalter0/tracer.git

Scoped print from point of defer call to exit of a block with 0..n
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


