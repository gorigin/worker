package worker

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

const (
	EventShutdown       = int(syscall.SIGTERM)
	EventReload         = int(syscall.SIGHUP)
	EventGarbageCollect = int(syscall.SIGUSR2)
	EventInfo           = int(syscall.SIGUSR1)
)

// StartCtx starts worker using provided function
// If second argument is true, does not return anything until work function done
func StartCtx(f func(ctx Context) error, blocking bool) error {
	// Building context
	ctx := BuildCtx()

	// Registering event listeners
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for sig := range c {
			switch sig {
			case os.Interrupt, syscall.SIGTERM:
				ctx.Info("Received %s, starting shutdown sequence", sig)
				ctx.Env().printToCtx(ctx)
				ctx.Env().Events <- EventShutdown
			case syscall.SIGHUP:
				ctx.Info("Received %s, sending reload event", sig)
				ctx.Env().Events <- EventReload
			case syscall.SIGUSR1:
				ctx.Info("Received %s, starting info sequence", sig)
				ctx.Env().printToCtx(ctx)
				ctx.Env().Events <- EventInfo
			case syscall.SIGUSR2:
				ctx.Info("Received %s, starting gc sequence", sig)
				runtime.GC()
				ctx.Env().Events <- EventGarbageCollect
			}
		}
		signal.Stop(c)
	}()

	if blocking {
		ctx.Info("Starting worker in blocking mode")
		if err := f(ctx); err != nil {
			ctx.Error("Worker done with error %s", err)
			return err
		} else {
			ctx.Info("Worker done without errors")
			return nil
		}
	}

	go f(ctx)
	return nil
}

// Start runs worker using provided function, it will be supplied with events channel on start
// If second argument is true, does not return anything until work function done
func Start(f func(chan int) error, blocking bool) error {
	return StartCtx(func(ctx Context) error {
		return f(ctx.Env().Events)
	}, blocking)
}

// Exec runs worker in blocking mode using provided function
func Exec(f func() error) error {
	return StartCtx(func(ctx Context) error {
		return f()
	}, true)
}

// Loop runs worker using provided function
// This function will be executed forever in loop until it returns error or system shutdown event received
// If second argument is true, does not return anything until work function done
func Loop(f func() error, blocking bool) error {
	return StartCtx(func(ctx Context) error {

		running := true

		go func() {
			for evt := range ctx.Env().Events {
				if evt == EventShutdown {

				}
			}
		}()

		for running {
			err := f()
			if err != nil {
				return err
			}
		}

		return nil
	}, blocking)
}
