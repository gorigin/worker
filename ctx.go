package worker

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Context holds environment information and some utility logging / metrics funcs
type Context interface {
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Error(string, ...interface{})
	Env() *Environment
}

// BuildCtx builds context object using environment settings
func BuildCtx() Context {
	c := &ctx{
		env: &Environment{
			Pid:       os.Getpid(),
			StartedAt: time.Now(),
			Args:      []string{},
			Options:   map[string]string{},
			Events:    make(chan int, 1),
		},
	}

	// Reading host
	if host, err := os.Hostname(); err == nil {
		c.env.Host = host
	}

	// Building arguments and options
	for _, v := range os.Args[1:] {
		if len(v) > 2 && v[0:2] == "--" {
			// This is option
			if i := strings.Index(v, "="); i == 0 {
				// No parameter
				c.env.Options[v[2:]] = ""
			} else {
				// Has parameter
				chunks := strings.Split(v[2:], "=")
				c.env.Options[chunks[0]] = chunks[1]
			}
		} else if len(v) > 1 && v[0:1] == "-" {
			// This is option
			if i := strings.Index(v, "="); i == 0 {
				// No parameter
				c.env.Options[v[1:]] = ""
			} else {
				// Has parameter
				chunks := strings.Split(v[1:], "=")
				c.env.Options[chunks[0]] = chunks[1]
			}
		} else {
			// This is argument
			c.env.Args = append(c.env.Args, v)
		}
	}

	// Checking verbosity
	if c.env.Options.HasOneOf("q", "quiet") {
		c.commonlog = func(string, ...interface{}) {}
	} else {
		if c.env.Options.HasOneOf("v", "verbose", "vv", "vvv") {
			c.debuglog = commonLog
		} else {
			c.commonlog = func(string, ...interface{}) {}
		}
		c.commonlog = commonLog
	}
	c.errorlog = commonLog

	// Updating environment
	go func(e *Environment) {
		for {
			e.update()
			time.Sleep(time.Second)
		}
	}(c.env)

	// Printing startup information
	c.Info("Starting worker")
	c.Info("Running on PID %d", c.Env().Pid)

	return c
}

type ctx struct {
	env                           *Environment
	debuglog, commonlog, errorlog func(string, ...interface{})
}

func (this ctx) Env() *Environment {
	return this.env
}

func (this ctx) Debug(p string, a ...interface{}) {
	this.errorlog(p, a...)
}

func (this ctx) Info(p string, a ...interface{}) {
	this.commonlog(p, a...)
}

func (this ctx) Error(p string, a ...interface{}) {
	this.errorlog(p, a...)
}

func commonLog(p string, a ...interface{}) {
	now := time.Now()
	fmt.Printf(
		"%02d %02d:%02d:%02d "+p+"\n",
		append([]interface{}{
			now.YearDay(),
			now.Hour(), now.Minute(), now.Second(),
		}, a...)...,
	)
}
