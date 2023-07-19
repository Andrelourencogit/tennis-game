package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const P = 4 // Número de pontos

type Jogador struct {
	nome   string
	pontos int
}

var (
	bola          = make(chan bool)
	wg            sync.WaitGroup
	mutex         sync.Mutex
	jogoTerminado bool
)

func main() {
	jogador1 := Jogador{nome: "Jogador 1", pontos: 0}
	jogador2 := Jogador{nome: "Jogador 2", pontos: 0}

	fmt.Println("Início do jogo!")

	wg.Add(2)
	go jogar(&jogador1, &jogador2)
	go jogar(&jogador2, &jogador1)

	bola <- true // Inicia o jogo com o jogador 1 mandando a bola

	wg.Wait()
	fmt.Println("\nFim do jogo!")

	if jogador1.pontos == P {
		fmt.Printf("%s ganhou o jogo!\n", jogador1.nome)
	} else {
		fmt.Printf("%s ganhou o jogo!\n", jogador2.nome)
	}
}

func jogar(jogador, adversario *Jogador) {
	defer wg.Done()
	for !jogoTerminado {
		recebendoBola := <-bola
		if jogoTerminado {
			return
		}
		if recebendoBola {
			fmt.Printf("%s (%d) vs (%d) %s\n", jogador.nome, jogador.pontos, adversario.pontos, adversario.nome)
			fmt.Printf("%s está esperando para receber a bola.\n", jogador.nome)
			time.Sleep(2 * time.Second) // Pausa de 2 segundos
			if perdeBola() {
				mutex.Lock()
				adversario.pontos++
				fmt.Printf("Jogador 1 perdeu a bola! %s ganhou o ponto!\n", adversario.nome)
				mutex.Unlock()
			} else {
				fmt.Printf("Jogador 1 mandou a bola de volta para %s!\n", jogador.nome)
			}
		} else {
			fmt.Printf("%s (%d) vs (%d) %s\n", jogador.nome, jogador.pontos, adversario.pontos, adversario.nome)
			fmt.Printf("%s está mandando a bola para %s.\n", jogador.nome, adversario.nome)
			time.Sleep(2 * time.Second) // Pausa de 2 segundos
			if perdeBola() {
				mutex.Lock()
				adversario.pontos++
				fmt.Printf("%s perdeu a bola! %s ganhou o ponto!\n", jogador.nome, adversario.nome)
				mutex.Unlock()
			} else {
				fmt.Printf("%s mandou a bola de volta!\n", adversario.nome)
			}
		}
		bola <- !recebendoBola
		if jogador.pontos == P || adversario.pontos == P {
			jogoTerminado = true
		}
	}
}

func perdeBola() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 0 // 50% de chance de perder a bola
}
