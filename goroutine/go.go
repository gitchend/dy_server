package goroutine

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	"github.com/panjf2000/ants/v2"
)

var _pool *ants.Pool

func init() {
	pool, err := ants.NewPool(math.MaxInt32, ants.WithExpiryDuration(time.Minute), ants.WithLogger(&antsLogger{}))
	if err != nil {
		panic(err)
	}
	_pool = pool
}

func Try(fn func(), catch func(ex interface{})) {
	defer func() {
		if r := recover(); r != nil {
			file, line := TraceFunc(fn)
			stack := string(debug.Stack())

			info := fmt.Sprintf("panic:%v\nfile:%v:%v\ntime:%v\nstack:%v\n", r, file, line, time.Now().Format(time.RFC3339), stack)
			os.Stderr.WriteString(info)

			if r == ExitErr {
			} else {
			}

			if catch != nil {
				catch(r)
			}
		}
	}()
	fn()
}

func GoLogic(fn func(), catch func(ex interface{})) {
	//go Try(func() { fn() }, catch)
	if err := _pool.Submit(func() { Try(func() { fn() }, catch) }); err != nil {
	} else {
	}
}

func TraceFunc(fn interface{}) (file string, line int) {
	if fn == nil {
		return "???", -1
	}
	rvalue := reflect.ValueOf(fn)
	if rvalue.Type().Kind() != reflect.Func {
		return "???", -1
	}
	pc := rvalue.Pointer()
	return runtime.FuncForPC(pc).FileLine(pc)
}

var ExitErr = errors.New("exit called")

func WaitExit() (waitFn func(), exitFn func()) {
	exit := make(chan os.Signal, 1)
	var closeOnce sync.Once
	signal.Notify(exit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	return func() { <-exit }, func() {
		closeOnce.Do(func() {
			close(exit)
		})
	}
}

type antsLogger struct{}

func (l *antsLogger) Printf(format string, args ...interface{}) {
}
