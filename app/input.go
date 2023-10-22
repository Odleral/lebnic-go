package app

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

// Input ввод в отдельном потоке. Доступные команды threads X, stop, resume X, view X
func Input() {
	var run = true
	var wg sync.WaitGroup
	cmdChan := make(chan Command, 1)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		select {
		case <-sigChan:
			run = false
			cmdChan <- Command{Type: "stop"}
		}
	}()

	wg.Add(1)
	go orchestrator(cmdChan, &wg)

loop:
	for run {
		var cmd []string

		for run {
			var arg string
			if _, err := fmt.Scan(&arg); err != nil {
				if errors.Is(err, io.EOF) {
					break loop
				}
			}

			if arg == "stop" || arg == "resume" {
				cmd = append(cmd, arg)
				break
			} else {
				if len(cmd) == 0 {
					cmd = append(cmd, arg)
				} else if _, err := strconv.Atoi(arg); err == nil && len(cmd) > 0 {
					cmd = append(cmd, arg)
					break
				} else {
					fmt.Println("input error")
					goto loop
				}
			}
		}

		if cmd[0] == "stop" {
			run = false
			cmdChan <- Command{Type: "stop"}
		} else if cmd[0] == "resume" {
			cmdChan <- Command{Type: "resume"}
			continue
		} else {
			num, err := strconv.Atoi(cmd[1])
			if err != nil {
				fmt.Println("expected number")
				goto loop
			}

			switch cmd[0] {
			case "start":
				cmdChan <- Command{Type: "start", Value: num}
				continue
			case "pause":
				cmdChan <- Command{Type: "pause", Value: num}
				continue
			case "view":
				cmdChan <- Command{Type: "view", Value: num}
				continue
			default:
				continue
			}
		}
	}

	wg.Wait()
}
