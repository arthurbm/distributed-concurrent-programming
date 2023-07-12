package main

import (
	"fmt"
	"sync"
	"time"
)

type Panela struct {
	porcoes int
	mtx     *sync.Mutex
	cv      *sync.Cond
}

func NovaPanela(n int) *Panela {
	mutex := &sync.Mutex{}
	return &Panela{
		porcoes: n,
		mtx:     mutex,
		cv:      sync.NewCond(mutex),
	}
}

func (p *Panela) PegarSopa() {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	for p.porcoes == 0 {
		p.cv.Signal() // acorda o cozinheiro
		p.cv.Wait()   // aguarda o cozinheiro encher a panela
	}
	p.porcoes--
}

func (p *Panela) EncherPanela(n int) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	for p.porcoes != 0 {
		p.cv.Wait() // aguarda a panela esvaziar
	}
	fmt.Println("Panela vazia, enchendo com novas porções...")
	p.porcoes = n
	time.Sleep(time.Second * 2) // para simular o tempo de cozinhar
	p.cv.Broadcast()             // acorda as pessoas
}

func pessoa(p *Panela, id int) {
	for {
		p.PegarSopa()
		fmt.Printf("Pessoa %d pegou sopa, restam %d porções\n", id, p.porcoes)
		time.Sleep(time.Second * 1) // para simular o tempo de comer
	}
}

func cozinheiro(p *Panela, n int) {
	for {
		p.EncherPanela(n)
	}
}

func main() {
	n := 5 // número de porções na panela
	p := NovaPanela(n)

	// cria threads para as pessoas
	for i := 1; i <= 10; i++ {
		go pessoa(p, i)
	}

	// cria thread para o cozinheiro
	go cozinheiro(p, n)

	// espera para que as goroutines tenham tempo de executar
	time.Sleep(time.Minute * 1)
}
