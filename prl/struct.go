package prl

import (
	"sync"
	"time"
)

// Task is a single unit of tasks
type Task struct{}

// Run simple task
func (t Task) Run() {
	var wg sync.WaitGroup
	wg.Add(1)
	t.run(&wg)
	wg.Wait()
}

func (t Task) run(wg *sync.WaitGroup) {
	it := uint64(0)
	go func() {
		now := time.Now()
		logger.Printf("Performing work at %s", now.Format(time.RFC850))
	rt:
		for timeout := time.After(time.Second * 3); ; {
			//logger.Println("New cycle")
			select {
			case <-timeout:
				logger.Println("Got channel")
				break rt
			default:
			}
			logger.Printf("Iteration %d here", it)
			it++
			time.Sleep(1)
		}
		logger.Printf("go routine stopped after %d iteration at %s", it, time.Now().Format(time.RFC850))
		logger.Printf("time spent: %s", time.Since(now))
		wg.Done()
	}()
}
