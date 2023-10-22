package app

import "fmt"

// sum - вывод резултата вычеслений
func sum(results chan State, threads int) []State {
	total := float64(0)
	var states []State

	if len(results) == threads {
		for i := 0; i < threads; i++ {
			if val, ok := <-results; ok {
				total += val.Result
				states = append(states, val)
			}
		}
	}

	fmt.Println("total:", 4*total)

	return states
}

func view(results chan State, threads int, wNum int) []State {
	var states []State

	if len(results) == threads {
		for i := 0; i < threads; i++ {
			if val, ok := <-results; ok {
				states = append(states, val)
				if val.WorkerNum == wNum {
					fmt.Printf("worker %d result: %2.12f\n", wNum, val.Result)
				}
			}
		}
	}
	return states
}
