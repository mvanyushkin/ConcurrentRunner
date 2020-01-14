package main

import (
	concurrent_runner "github.com/mvanyushkin/ConcurrentRunner/concurrent-runner"
	"time"
)

func main() {

	tasks := make([]func() error, 10)
	for k := 0; k < 10; k++ {
		tasks[k] = func() error {
			println("I'm running")
			time.Sleep(time.Second * 3)
			return nil
		}
	}

	concurrent_runner.RunConcurrent(tasks, 3, 4)
}
