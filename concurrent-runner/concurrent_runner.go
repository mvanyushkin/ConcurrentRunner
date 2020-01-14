package concurrent_runner

func RunConcurrent(tasks []func() error, concurrencyDegree uint32, errorsThreshold uint32) Result {
	runner := concurrentRunner{
		throttlingCh:       make(chan int, concurrencyDegree),
		performedCh:        make(chan int, 1),
		syncCh:             make(chan int, 1),
		maxErrorsThreshold: errorsThreshold,
		totalFired:         0,
		amountFired:        len(tasks),
	}

	runner.internalRun(tasks)
	return Result{
		ok:          !runner.hopelessRun,
		totalFired:  runner.totalFired,
		totalFailed: runner.errorsCount,
	}
}

type Result struct {
	ok          bool
	totalFired  uint32
	totalFailed uint32
}

type concurrentRunner struct {
	throttlingCh       chan int
	performedCh        chan int
	syncCh             chan int
	maxErrorsThreshold uint32
	errorsCount        uint32
	totalFired         uint32
	hopelessRun        bool
	amountFired        int
}

func (c *concurrentRunner) internalRun(tasks []func() error) {

	for index, task := range tasks {

		if c.hopelessRun {
			break
		}
		c.throttlingCh <- index
		go func(taskId int, task func() error) {
			if c.hopelessRun {
				_ = <-c.throttlingCh
				c.performedCh <- 0
				return
			}

			err := task()

			c.enterCriticalSection()
			c.totalFired++
			if err != nil {
				c.errorsCount++
				c.hopelessRun = c.errorsCount >= c.maxErrorsThreshold
			}
			c.exitCriticalSection()
			_ = <-c.throttlingCh
			c.performedCh <- 0
		}(index, task)
	}
	c.waitingForCompletion()
}

func (c *concurrentRunner) waitingForCompletion() {
	captured := 0
	for range c.performedCh {
		captured++
		if captured == int(c.totalFired) && c.hopelessRun {
			break
		}

		if captured == c.amountFired {
			break
		}
	}
}

func (c *concurrentRunner) enterCriticalSection() {
	c.syncCh <- 0
}

func (c *concurrentRunner) exitCriticalSection() {
	_ = <-c.syncCh
}
