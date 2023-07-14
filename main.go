package main

import (
	"fmt"
	"sync"
	"time"
)

type Runner struct {
	id    uint
	stick chan uint
	wait  *sync.WaitGroup
	next  uint
}

var NUM_RUNNERS uint

func main() {
	//Defina aqui a quantidade de corredores desejada
	//Ao mudar os parâmetros, talvez seja necessário dar KILL no terminal e abrir um novo, com as alterações feitas, pois está sendo usada uma variável global.
	NUM_RUNNERS = 4
	stick := make(chan uint, 1) // Esse channel precisa ter um buffer para evitar deadlock
	runnerSlice := setRunners(stick)
	wg := sync.WaitGroup{}
	wg.Add(len(runnerSlice))
	//Note que está sendo adicionado no wait group a quantidade de corredores criados no setRunners.

	for _, runner := range runnerSlice {
		go run(&wg, runner)
	}
	//O início da corrida é aqui.
	stick <- 1
	fmt.Printf("The runner: [%d] start the track!\n", 1)
	wg.Wait()
	fmt.Println("The track has been completed!")
}

// Função responsável por iniciar a corrida e fazer a checagem da passagem de bastão.
// Estou passando ponteiros para que sejam alterados os objetos diretamente, ao invés de criar cópia.
func run(wg *sync.WaitGroup, runner *Runner) {
	defer wg.Done()
	temp := <-runner.stick
	// Essa condição é necessária para evitar que um corredor inicie tentando "pegar" um bastão de outro corredor
	if runner.id != temp {
		runner.stick <- temp
		return
	}

	time.Sleep(1 * time.Second)
	if runner.id < NUM_RUNNERS {
		fmt.Printf("The runner: [%d] finished.\n", runner.id)
		time.Sleep(1 * time.Second)
		fmt.Printf("The runner: [%d] has started!\n", runner.next)
		runner.stick <- runner.next
	}
	time.Sleep(2 * time.Second)
}

// Função responsável por definir os corredores em um slice. Note que podem ser criados quantos corredores quiser, alterando o parâmetro dentro do make.
func setRunners(stick chan uint) []*Runner {
	runnerSlice := make([]*Runner, NUM_RUNNERS)
	for i := 0; i < len(runnerSlice); i++ {
		runnerSlice[i] = &Runner{
			id:    uint(i + 1),
			stick: stick,
			next:  uint(i + 2),
		}
	}
	return runnerSlice
}
