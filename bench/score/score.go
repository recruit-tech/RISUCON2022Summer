package score

import (
	"sync"
)

var (
	mu        sync.RWMutex
	score     int64 = 0
	level     int64 = 1
	levelChan chan int64
)

func init() {
	levelChan = make(chan int64, 1)
}

const StepScoreByLevel = 500

func Increment() {
	mu.Lock()
	defer mu.Unlock()
	score += 1
	if score >= level*StepScoreByLevel {
		level += 1
		levelChan <- level
	}
}

func Sum() int64 {
	mu.RLock()
	defer mu.RUnlock()
	return score
}

func Level() int64 {
	mu.RLock()
	defer mu.RUnlock()
	return level
}

func LevelUp() chan int64 {
	return levelChan
}
