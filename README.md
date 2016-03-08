# Worker

Simple wrapper for long-running applications with some helpful features:

* All events delegated to executing code
* Wrapper registers OS signals listener for following events:
  * `SIGTERM` and `SIGINT` - used to gracefully shutdown application if it possible
  * `SIGHUP`
  * `SIGUSR1` - prints current heap allocations and various application info
  * `SIGUSR2` - runs garbage collection
* Wrapper provides logging API and parses incoming arguments to determine verbosity
  * `-q` or `--quiet` - no output
  * `-v` or `--verbose` - verbose output
  
  
# Usage

```
package main

import 'github.com/gorigin/worker'


func main() {
    worker.StartCtx(workFunc, true)
}

func workFunc(ctx worker.Context) error {
    ctx.Info("Starting application")
    
    running := true
    
    // Registering event listeners
    go func() {
        for evt := range ctx.Env().Events {
            if evt == worker.EventShutdown {
                // Shutdown event
                running = false
            }
        }
    } ()
    
    // Main job
    for running {
        // BL goes here
    }
    
    ctx.Info("Application done")
    return nil
}
```

# Events

All events are plain integers

* `worker.EventShutdown` - used for both `SIGTERM` and `SIGINT`
* `worker.EventReload` - stands for `SIGHUP`
* `worker.EventInfo` - stands for `SIGUSR1`
* `worker.EventGarbageCollect` - stands for `SIGUSR2`

# Available start initializers 

## StartCtx

`func StartCtx(f func(ctx Context) error, blocking bool) error`

* Takes function to execute as first argument.
* If `blocking` set to `true` locks current goroutine until provided function done (recommended)


## Start

`Start(f func(chan int) error, blocking bool) error`

* `chan int`, supplied to first argument function, is events channel
* If `blocking` set to `true` locks current goroutine until provided function done (recommended)


## Exec

`func Exec(f func() error) error`
 
Runs provided function without any event delivery


## Loop

`func Loop(f func() error, blocking bool) error`

Runs provided function in loop until shutdown event received