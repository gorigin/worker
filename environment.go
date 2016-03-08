package worker

import (
	"runtime"
	"sync"
	"time"
)

// Environment struct holds information about all worker environment
// including command-line arguments, memory stats and event channel
type Environment struct {
	Pid       int
	Host      string
	StartedAt time.Time
	Args      Args
	Options   Options
	Events    chan int

	// Stat
	updateLock sync.Mutex
	Goroutines int
	Memory     runtime.MemStats
}

// Uptime method returns current uptime
func (this *Environment) Uptime() time.Duration {
	return time.Now().Sub(this.StartedAt)
}

// update updates environment stats information
func (this *Environment) update() {
	this.updateLock.Lock()
	defer this.updateLock.Unlock()

	this.Goroutines = runtime.NumGoroutine()
	runtime.ReadMemStats(&this.Memory)
}

// printToCtx prints environment information to context
func (this *Environment) printToCtx(ctx Context) {
	ctx.Info(
		"Meta: [Pid: %d][Goroutines: %d][Uptime: %.0f s][Host: %s]",
		this.Pid,
		this.Goroutines,
		this.Uptime().Seconds(),
		this.Host,
	)
	ctx.Info(
		"Heap: [Alloc: %d][Sys: %d][Idle: %d][InUse: %d][Objects: %d]",
		this.Memory.HeapAlloc,
		this.Memory.HeapSys,
		this.Memory.HeapIdle,
		this.Memory.HeapInuse,
		this.Memory.HeapObjects,
	)
	ctx.Info(
		"  GC: [Num: %d][Next: %d][Last: %s][Total pause: %s]",
		this.Memory.NumGC,
		this.Memory.NextGC,
		time.Unix(0, int64(this.Memory.LastGC)),
		time.Duration(this.Memory.PauseTotalNs),
	)
}
