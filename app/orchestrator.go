package app

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Orchestrator - функция, в котором производиться управление воркерами.
// Принимает единственный аргумент - канал типа Command. Управление состоянием воркеров
// осуществляется через структуры Stage и State.
//
// Stage отвечает за старт и остановку воркера, State за внутреннее состояние воркера.
//
// Реализация команд:
//
// 1. start n - запуск n воркеров
//
// 2. stop - остановка всех воркеров
//
// 3. pause x - остановка всех воркеров на х секунд, при этом произайдет остановка с сохранением состояния воркеров,
// по исчечению времени запуск
//
// 4. resume - продолжение вычеслений, отменяет ожидание в команде pause
//
// 5. view x - вывод значения из потока x в течении 10 секунд
func orchestrator(cmd chan Command, w *sync.WaitGroup) {
	defer w.Done()

	var wg sync.WaitGroup
	var threads int
	var pause bool
	var stage Stage
	var states []State
	results := make(chan State, runtime.NumCPU())
	resume := make(chan int, 1)

	for c := range cmd {
		switch {
		case c.Type == "start":
			if threads != 0 {
				fmt.Println("workers started")
				break
			}

			if pause {
				fmt.Println("expected resume")
				break
			}

			threads = c.Value
			for i := 1; i <= c.Value; i++ {
				wg.Add(1)
				go worker(&stage, c.Value, State{WorkerNum: i}, results, &wg)
			}
		case c.Type == "pause":

			if pause {
				fmt.Println("expected resume")
				break
			}

			stage.Set(1)
			wg.Wait()
			states = sum(results, threads)

			pause = true
			stage.Set(0)
			go func(resume chan int, c Command) {
				select {
				case <-time.After(time.Duration(c.Value) * time.Second):
					fmt.Println("info: resume compute after", c.Value, "seconds")
					for _, v := range states {
						wg.Add(1)
						go worker(&stage, threads, v, results, &wg)
					}
					pause = false
				case <-resume:
					fmt.Println("info: resume compute")
					for _, v := range states {
						wg.Add(1)
						go worker(&stage, threads, v, results, &wg)
					}
					pause = false
				}
			}(resume, c)
		case c.Type == "resume":
			resume <- 1
		case c.Type == "stop":

			if pause {
				fmt.Println("expected resume")
				break
			}

			stage.Set(1)
			wg.Wait()
			sum(results, threads)
			close(cmd)
			break
		case c.Type == "view":
			if pause {
				fmt.Println("expected resume")
				break
			}
			for i := 0; i < 10; i++ {
				stage.Set(1)
				wg.Wait()
				states = view(results, threads, c.Value)
				stage.Set(0)
				for _, v := range states {
					wg.Add(1)
					go worker(&stage, threads, v, results, &wg)
				}
				time.Sleep(time.Second)
			}
		}
	}
}
