package model

import (
	"log"
	"math/rand"
	"sync"
)

type Person struct {
	Id               int
	StartingFloor    int
	DestinationFloor int
	TravelTime       int
	IsTraveling      bool
	WaitingTime      int
	IsWaiting        bool
	PickUpSignal     chan interface{}
	DropOffSignal    chan interface{}
}

func NewPerson(id int, currentFloor int, destinationFloor int) *Person {
	return &Person{
		Id:               id,
		StartingFloor:    currentFloor,
		DestinationFloor: destinationFloor,
		PickUpSignal:     make(chan interface{}),
		DropOffSignal:    make(chan interface{}),
	}
}

func NewRandomPerson(id int, floorCount int) *Person {
	startingFloor := rand.Intn(floorCount)
	destinationFloor := rand.Intn(floorCount)
	if destinationFloor == startingFloor {
		destinationFloor = (destinationFloor + 1) % floorCount
	}
	return NewPerson(id, startingFloor, destinationFloor)
}

func (p *Person) Run(controller *Controller, group *sync.WaitGroup) {
	defer group.Done()
	controller.Floors[p.StartingFloor] <- p
	p.IsWaiting = true
	<-p.PickUpSignal
	log.Printf("person   %03v: traveling from floor %v to floor %v", p.Id, p.StartingFloor, p.DestinationFloor)
	p.IsWaiting = false
	p.IsTraveling = true
	<-p.DropOffSignal
	log.Printf("person   %03v: arrived on floor %v", p.Id, p.DestinationFloor)
	p.IsTraveling = false
}
