package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go input(ctx, cancel, &wg)
	wg.Wait()
	fmt.Println("Gracefull Shutdown")
}

func worker(ctx context.Context, workerNum int, threads int, results chan float64) {
	var i int
	var result float64

	run := true
	fmt.Printf("worker %d started\n", workerNum)
	for run {
		select {
		case <-ctx.Done():
			results <- result
			run = false
		default:
			i += threads

			result += math.Pow(-1, float64(i-workerNum)) / (2*float64(i-workerNum) + 1)
		}
	}
}

func sum(results chan float64, wg *sync.WaitGroup) {
	defer wg.Done()
	total := float64(0)

	for v := range results {
		total += v
	}

	fmt.Println("total:", 4*total)
}

// input ввод в отдельном потоке. Доступные команды threads X, stop, resume X, view X
func input(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%v\n", r)
			cancel()
		}
	}()
	defer wg.Done()

	var cmd string
	var n int
	var working bool
	var run = true

	results := make(chan float64, 10)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		select {
		case <-sigChan:
			cancel()
			working = false
			wg.Add(1)
			go sum(results, wg)
			close(results)
			run = false
		}
	}()

	for run {
		_, err := fmt.Scanf("%s %d\n", &cmd, &n)
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("input error:", err)
			}

			continue
		}

		switch cmd {
		case "t":
			if working {
				fmt.Println("Now the program is running, enter stop")
			}

			working = true

			for i := 1; i <= n; i++ {
				go worker(ctx, i, n, results)
			}
		case "stop":
			cancel()
			working = false
			wg.Add(1)
			go sum(results, wg)
			close(results)
			run = false
		}
	}
}
