package os

import "sync"

var wg sync.WaitGroup

func init() {

}

func AddThread() {
	wg.Add(1)
}

func DeleteThread() {
	wg.Done()
}

func Wait() {
	wg.Wait()
}