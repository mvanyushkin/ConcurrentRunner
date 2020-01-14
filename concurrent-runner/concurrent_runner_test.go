package concurrent_runner

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func TestWhenAllTasksAreOk(t *testing.T) {
	amount := 123345
	tasks := make([]func() error, amount)
	for k := 0; k < amount; k++ {
		tasks[k] = presumableCorrectFunc
	}

	goroutinsCountBefore := runtime.NumGoroutine()
	concurrencyDegree := 663
	result := RunConcurrent(tasks, uint32(concurrencyDegree), 1)
	assert.True(t, result.ok)
	assert.Equal(t, amount, int(result.totalFired))
	assert.Equal(t, 0, int(result.totalFailed))

	// Waiting for gc
	time.Sleep(time.Second)
	assert.Equal(t, goroutinsCountBefore, runtime.NumGoroutine())
}

func TestWhenFirstTasksAreGoingToFail(t *testing.T) {
	amount := 6
	failingFirst := 2
	tasks := make([]func() error, amount)
	for k := 0; k < amount; k++ {
		if k < failingFirst {
			tasks[k] = presumableFailingFunc
		} else {
			tasks[k] = presumableCorrectFunc
		}
	}

	goroutinsCountBefore := runtime.NumGoroutine() + 1
	var concurrencyDegree uint32 = 2
	var errorsThreshold = uint32(failingFirst) - 1

	result := RunConcurrent(tasks, concurrencyDegree, errorsThreshold)
	assert.False(t, result.ok)
	assert.GreaterOrEqual(t, int(concurrencyDegree)+failingFirst, int(result.totalFired))

	// Waiting for gc
	time.Sleep(time.Second)
	assert.Equal(t, goroutinsCountBefore, runtime.NumGoroutine())
}

func presumableCorrectFunc() error {
	time.Sleep(time.Microsecond * time.Duration(rand.Intn(100)))
	return nil
}

func presumableFailingFunc() error {
	time.Sleep(time.Microsecond * time.Duration(rand.Intn(100)))
	return errors.New("Ooops")
}
