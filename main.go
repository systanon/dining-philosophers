package main

import (
	"fmt"
	"sync"
	"time"
)

const numPhilosophers = 5
const numMeals = 3

type ChopSticks struct {
	sync.Mutex
}

type Philosopher struct {
	id         int
	rightStick *ChopSticks
	leftStick  *ChopSticks
	host       chan struct{}
	milsCount  int
}

func (p *Philosopher) Eat(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < p.milsCount; i++ {
		p.host <- struct{}{}

		p.rightStick.Lock()
		p.leftStick.Lock()

		fmt.Printf("starting to eat %d\n", p.id)
		time.Sleep(time.Millisecond * 500)
		fmt.Printf("finishing eating %d\n", p.id)
		p.rightStick.Unlock()
		p.leftStick.Unlock()

		<-p.host
	}
}

func main() {

	chopSticks := make([]*ChopSticks, numPhilosophers)

	for i := 0; i < numPhilosophers; i++ {
		chopSticks[i] = &ChopSticks{}
	}
	host := make(chan struct{}, 2)
	philosophers := make([]*Philosopher, numPhilosophers)

	for i := 0; i < numPhilosophers; i++ {
		philosophers[i] = &Philosopher{
			id:         i + 1,
			leftStick:  chopSticks[i],
			rightStick: chopSticks[(i+1)%numPhilosophers],
			host:       host,
			milsCount:  numMeals,
		}
	}

	var wg sync.WaitGroup
	wg.Add(numPhilosophers)

	for _, p := range philosophers {
		go p.Eat(&wg)
	}
	wg.Wait()
}
