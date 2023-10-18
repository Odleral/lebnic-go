package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	Stop = iota
	Pause
	Resume
	View
)

type Command struct {
	t int
	v int
}

func main() {
	if err := run(); err != nil {
		fmt.Println("[main.run]", err)
	}
}

func run() error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered", r)
		}
	}()
	var threads int
	var totalAns float64
	var wg sync.WaitGroup

	cmd := make(chan Command, 1)
	ctx, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	fmt.Print("Threads count: ")
	_, err := fmt.Scan(&threads)
	if err != nil {
		fmt.Println("[main.run] fmt.Scan", err)
	}

	wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()

		num := make(chan float64, 1000)
		ans := make(chan float64, 1000)

		for i := 0; i < threads; i++ {
			wg.Add(1)
			go worker(ctx, num, ans, wg)
		}

		go generator(ctx, num, cmd)
		go sum(ans, &totalAns)
		go enter(ctx, cmd)

	}(ctx, &wg)

	wg.Wait()

	return nil
}

func worker(ctx context.Context, num chan float64, ans chan float64, wg *sync.WaitGroup) {
	defer wg.Done()

	for v := range num {
		select {
		case <-ctx.Done():
			return
		default:
			ans <- math.Pow(-1, v) / (2*v + 1)
		}

	}
}

func generator(ctx context.Context, num chan float64, cmd chan Command) {
	i := 0
	for {
		i++
		select {
		case g := <-cmd:
			switch g.t {
			case Pause:
				time.Sleep(time.Duration(g.v) * time.Second)
			case Stop:
				break
			}
		case <-ctx.Done():
			break
		default:
			num <- float64(i)
		}
	}
}

func sum(ans chan float64, totalAns *float64) {
	for v := range ans {
		*totalAns += v
	}
}

func enter(ctx context.Context, cmd chan Command) {
	for {
		var scan string

		fmt.Print("Command: ")
		_, err := fmt.Scan(&scan)
		if err != nil {
			fmt.Println("[main.run] fmt.Scan", err)
			return
		}

		scanSplit := strings.Split(scan, " ")

		switch scanSplit[0] {
		case "stop":
			fmt.Println("Stop")
			cmd <- Command{t: Stop}
			fmt.Println()
		case "pause":
			fmt.Println("Pause", scanSplit[1])
			v, err := strconv.Atoi(scanSplit[1])
			if err != nil {
				fmt.Println("Error", err)
				break
			}
			cmd <- Command{t: Pause, v: v}
		case "resume":
			fmt.Println("Resume", scanSplit[1])
			v, err := strconv.Atoi(scanSplit[1])
			if err != nil {
				fmt.Println("Error", err)
				break
			}
			cmd <- Command{t: Resume, v: v}
		}
	}
}
