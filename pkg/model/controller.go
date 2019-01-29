package model

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Controller struct {
	Simulation *Simulation
	Floors     []chan *Person
	Persons    *sync.Map
	Elevators  map[int]*Elevator
}

func NewController(simulation *Simulation) *Controller {
	floors := make([]chan *Person, simulation.FloorCount)
	for i := range floors {
		floors[i] = make(chan *Person, simulation.ElevatorCapacity)
	}
	return &Controller{
		Simulation: simulation,
		Floors:     floors,
		Persons:    &sync.Map{},
		Elevators:  map[int]*Elevator{},
	}
}

func (c *Controller) Evaluate(file io.Writer) {
	if _, err := fmt.Fprintln(file, "Elevators"); err != nil {
		log.Fatal(err)
		return
	}
	for _, elevator := range c.Elevators {
		_, err := fmt.Fprintf(file, "  Elevator %02v:\n    %v steps\n    %v picked up\n",
			elevator.Id, elevator.StepsTaken, elevator.TotalPickedUp)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	_, err := fmt.Fprintf(file, "Persons: mean waiting time %v, mean traveling time %v\n",
		meanWaitingTime(c.Persons), meanTravelTime(c.Persons))
	if err != nil {
		log.Fatal(err)
		return
	}
	c.Persons.Range(func(_, value interface{}) bool {
		person := value.(*Person)
		_, err := fmt.Fprintf(file, "  Person %04v:\n    waited %v steps\n    traveled %v steps\n",
			person.Id, person.WaitingTime, person.TravelTime)
		if err != nil {
			log.Fatal(err)
			return false
		}
		return true
	})
}

func (c *Controller) Run(controlSimulation func(*Controller)) {
	for i := 0; i < c.Simulation.ElevatorCount; i++ {
		elevator := NewElevator(i, c.Simulation.ElevatorCapacity)
		c.Elevators[elevator.Id] = elevator
		go elevator.Run(c)
	}
	var group sync.WaitGroup
	group.Add(c.Simulation.PersonCount)
	defer group.Wait()
	go c.generatePersons(&group)
	go controlSimulation(c)
}

func (c *Controller) generatePersons(group *sync.WaitGroup) {
	for i := 0; i < c.Simulation.PersonCount; i++ {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Microsecond)
		person := NewRandomPerson(i, c.Simulation.FloorCount)
		c.Persons.Store(person.Id, person)
		go person.Run(c, group)
	}
}

func meanTravelTime(persons *sync.Map) float32 {
	totalTravelTime := 0
	personCount := 0
	persons.Range(func(_, value interface{}) bool {
		person := value.(*Person)
		totalTravelTime += person.TravelTime
		personCount++
		return true
	})
	return float32(totalTravelTime) / float32(personCount)
}

func meanWaitingTime(persons *sync.Map) float32 {
	totalWaitingTime := 0
	personCount := 0
	persons.Range(func(_, value interface{}) bool {
		person := value.(*Person)
		totalWaitingTime += person.WaitingTime
		personCount++
		return true
	})
	return float32(totalWaitingTime) / float32(personCount)
}
