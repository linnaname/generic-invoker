package routine

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	MAX_RUN_SIZE = 10000000
)

var sum int64
var gobalTestWG = sync.WaitGroup{}

func TestPool(t *testing.T) {
	pool, err := New(10)
	assert.NoError(t, err)
	assert.NotNil(t, pool)
	assert.Equal(t, pool.State(), int64(RUNNING))

	wg := new(sync.WaitGroup)

	for i := 0; i < 30; i++ {
		wg.Add(1)

		task := &Task{
			Handler: func(v ...interface{}) {
				wg.Done()
				fmt.Println(v)
			},
		}
		task.Params = []interface{}{i, "linnana"}
		pool.Add(task)
	}
	assert.Equal(t, pool.State(), int64(RUNNING))
	wg.Wait()

	pool.Close()
	assert.Equal(t, pool.State(), int64(TERMINATED))
}

func TestNew(t *testing.T) {
	pool, err := New(10)
	assert.NoError(t, err)
	assert.NotNil(t, pool)
	assert.Equal(t, pool.State(), int64(RUNNING))
	assert.Equal(t, pool.Capacity(), uint64(10))
}

func TestAdd(t *testing.T) {
	pool, _ := New(10)
	assert.Equal(t, pool.RunningWorkerCnt(), uint64(0))
	pool.Add(&Task{
		Handler: func(v ...interface{}) {
			time.Sleep(1000)
		},
	})
	assert.Equal(t, pool.runningWorkerCnt, uint64(1))

}

func TestClose(t *testing.T) {
	pool, _ := New(10)
	wg := new(sync.WaitGroup)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		pool.Add(&Task{
			Handler: func(v ...interface{}) {
				wg.Done()
			},
		})
	}
	wg.Wait()
	pool.Close()

	assert.Equal(t, pool.State(), int64(TERMINATED))
	err := pool.Add(&Task{
		Handler: func(v ...interface{}) {},
	})

	assert.Error(t, err)
	fmt.Println(err)
}

func TestRunningWorkerCnt(t *testing.T) {
	pool, _ := New(10)
	wg := new(sync.WaitGroup)

	assert.Equal(t, pool.runningWorkerCnt, uint64(0))

	for i := 0; i < 100; i++ {
		wg.Add(1)
		pool.Add(&Task{
			Handler: func(v ...interface{}) {
				wg.Done()
			},
		})
	}
	assert.Equal(t, pool.runningWorkerCnt, uint64(10))
	wg.Wait()
	pool.Close()
	assert.Equal(t, pool.runningWorkerCnt, uint64(0))
}

func TestPanicHandler(t *testing.T) {
	pool, _ := New(20)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	pool.PanicHandler = func(v interface{}) {
		wg.Done()
		fmt.Printf("Handling paic:%s....", v)
	}
	pool.Add(&Task{
		Handler: func(v ...interface{}) {
			panic("I am panic")
		},
	})
	wg.Wait()
}

func task(v ...interface{}) {
	gobalTestWG.Done()
	for i := 0; i < 100; i++ {
		atomic.AddInt64(&sum, 1)
	}
}

func BenchmarkOrginalRoutine(b *testing.B) {
	for i := 0; i < MAX_RUN_SIZE; i++ {
		gobalTestWG.Add(1)
		go task()
	}
	gobalTestWG.Wait()
}

func BenchmarkPool(b *testing.B) {
	pool, err := New(20)
	if err != nil {
		b.Error(err)
	}

	task := &Task{
		Handler: task,
	}

	for i := 0; i < MAX_RUN_SIZE; i++ {
		gobalTestWG.Add(1)
		pool.Add(task)
	}
	gobalTestWG.Wait()
}
