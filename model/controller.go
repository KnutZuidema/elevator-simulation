package model

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type Controller struct {
	ControllingAlgorithms map[string]func()
	Floors                []chan *Person
	persons               []*Person
	elevators             []*Elevator
}

func NewController(floorCount int) *Controller {
	floors := make([]chan *Person, floorCount)
	for i := range floors {
		floors[i] = make(chan *Person)
	}
	return &Controller{
		ControllingAlgorithms: make(map[string]func()),
		Floors:                floors,
		persons:               []*Person{},
	}
}

func (c *Controller) Run(simulation *Simulation) {
	log.Print("Starting simulation")
	for i := 0; i < simulation.ElevatorCount; i++ {
		elevator := NewElevator(simulation.ElevatorCapacity)
		c.elevators = append(c.elevators, elevator)
		go elevator.Run(c)
	}
	var group sync.WaitGroup
	group.Add(simulation.PersonCount)
	go func() {
		for i := 0; i < simulation.PersonCount; i++ {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			person := NewRandomPerson(simulation.FloorCount)
			c.persons = append(c.persons, person)
			go person.Run(c, group)
		}
	}()
	go func() {

	}()
}
