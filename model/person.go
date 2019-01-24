package model

import (
	"math/rand"
	"sync"
)

type Person struct {
	CurrentFloor     int
	DestinationFloor int
	TravelTime       int
	IsUnderWay       bool
	WaitingTime      int
	IsWaiting        bool
	PickUpSignal     chan interface{}
	DropOffSignal    chan interface{}
}

func NewPerson(currentFloor int, destinationFloor int) *Person {
	return &Person{
		CurrentFloor:     currentFloor,
		DestinationFloor: destinationFloor,
		PickUpSignal:     make(chan interface{}),
		DropOffSignal:    make(chan interface{}),
	}
}

func NewRandomPerson(floorCount int) *Person {
	currentFloor := rand.Intn(floorCount)
	destinationFloor := rand.Intn(floorCount)
	return NewPerson(currentFloor, destinationFloor)
}

func (p *Person) Run(controller *Controller, group sync.WaitGroup) {
	defer group.Done()
	controller.Floors[p.CurrentFloor] <- p
	p.IsWaiting = true
	<-p.PickUpSignal
	p.IsWaiting = false
	p.IsUnderWay = true
	<-p.DropOffSignal
	p.IsUnderWay = false
}
