package routine

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"
)

/**
state of pool not the state of worker or routine
*/
const (
	//pool can accept new task
	RUNNING = 1
	//task maybae still running
	STOPED = 2
	//task all stoped and channel was close
	TERMINATED = 3
)

type RoutinePool struct {
	sync.Mutex

	capacity         uint64
	runningWorkerCnt uint64
	state            int64
	taskBufferChan   chan *Task
	PanicHandler     func(interface{})
}

/**
new pool
just init datastruct of pool
*/
func New(capacity uint64) (*RoutinePool, error) {
	if capacity <= 0 {
		return nil, errors.New("illegal argument:capacity invalid")
	}

	return &RoutinePool{
		capacity:       capacity,
		state:          RUNNING,
		taskBufferChan: make(chan *Task, capacity), //buffer channel simulation to block queue
	}, nil
}

/**
add task to pool and run
*/
func (pl *RoutinePool) Add(task *Task) error {
	if pl.state > RUNNING {
		return errors.New("routine pool already closed")
	}

	pl.Lock()
	if pl.RunningWorkerCnt() < pl.Capacity() {
		pl.run()
	}

	//make sure pool still running before submit task
	if pl.state == RUNNING {
		pl.taskBufferChan <- task
	}
	pl.Unlock()

	return nil
}

/**
close pool
1.change state of pool to stoped
2.close all worker
3.close channel
4.change state of pool to terminated
*/
func (pl *RoutinePool) Close() {
	//can't close more than once
	if pl.state >= STOPED {
		return
	}

	pl.setState(STOPED)

	//before close channel we need to make sure all task is finished
	for len(pl.taskBufferChan) > 0 {
	}
	pl.closeChan()
	pl.resetRunningWorkerCnt()
	pl.setState(TERMINATED)
}

/**
real work
TODO add and excute
*/
func (pl *RoutinePool) run() {
	pl.incrementRunningWorkerCnt()

	go func() {
		defer func() {
			pl.decrementRunningWorkerCnt()
			//when have PanicHandler recover and invoke it
			if r := recover(); r != nil {
				if pl.PanicHandler != nil {
					pl.PanicHandler(r)
				} else {
					log.Printf("Something wrong: %s\n", r)
				}
			}
		}()

		for {
			select {
			case task, ok := <-pl.taskBufferChan:
				if !ok {
					return
				}
				task.Handler(task.Params...)
			}
		}
	}()
}

/**
channel can't close more than once
make sure thread-safe
*/
func (pl *RoutinePool) closeChan() {
	pl.Lock()
	defer pl.Unlock()
	close(pl.taskBufferChan)
}

/**
to guarantee thread-safe by use atomic
*/
func (p *RoutinePool) incrementRunningWorkerCnt() {
	atomic.AddUint64(&p.runningWorkerCnt, 1)
}

func (p *RoutinePool) decrementRunningWorkerCnt() {
	atomic.AddUint64(&p.runningWorkerCnt, ^uint64(0))
}

func (p *RoutinePool) resetRunningWorkerCnt() {
	atomic.StoreUint64(&p.runningWorkerCnt, uint64(0))
}

/**
running worker count of pool
*/
func (pl *RoutinePool) RunningWorkerCnt() uint64 {
	return atomic.LoadUint64(&pl.runningWorkerCnt)
}

/**
capactiy of pool
*/
func (pl *RoutinePool) Capacity() uint64 {
	return pl.capacity
}

/**
get state of pool
*/
func (pl *RoutinePool) State() int64 {
	pl.Lock()
	defer pl.Unlock()
	return pl.state
}

func (pl *RoutinePool) setState(state int64) {
	pl.Lock()
	defer pl.Unlock()
	pl.state = state
}
