package app

import (
	"math"
	"sync"
)

// Worker - выполняет вычесления
func worker(stage *Stage, maxWorkers int, state State, results chan State, wg *sync.WaitGroup) {
	defer wg.Done()

	var i = state.I
	var result = state.Result
	var workerNum = state.WorkerNum

	run := true

	for run {
		switch stage.Get() {
		case 0:
			i += maxWorkers
			result += math.Pow(-1, float64(i-workerNum)) / (2*float64(i-workerNum) + 1)
		case 1:
			results <- State{i, result, workerNum}
			run = false
		}
	}
}
