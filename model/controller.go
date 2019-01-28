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
	simulation *Simulation
	Floors     []chan *Person
	persons    *sync.Map
	elevators  map[int]*Elevator
}

func NewController(simulation *Simulation) *Controller {
	floors := make([]chan *Person, simulation.FloorCount)
	for i := range floors {
		floors[i] = make(chan *Person, simulation.ElevatorCapacity)
	}
	return &Controller{
		simulation: simulation,
		Floors:     floors,
		persons:    &sync.Map{},
		elevators:  map[int]*Elevator{},
	}
}

func (c *Controller) Evaluate(file io.Writer) {
	if _, err := fmt.Fprintln(file, "Elevators"); err != nil {
		log.Fatal(err)
		return
	}
	for _, elevator := range c.elevators {
		_, err := fmt.Fprintf(file, "  Elevator %02v: %v steps\n", elevator.Id, elevator.StepsTaken)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	_, err := fmt.Fprintf(file, "Persons: mean waiting time %v, mean traveling time %v\n",
		meanWaitingTime(c.persons), meanTravelTime(c.persons))
	if err != nil {
		log.Fatal(err)
		return
	}
	c.persons.Range(func(_, value interface{}) bool {
		person := value.(*Person)
		_, err := fmt.Fprintf(file, "  Person %04v: waited %v steps, traveled %v steps\n",
			person.Id, person.WaitingTime, person.TravelTime)
		if err != nil {
			log.Fatal(err)
			return false
		}
		return true
	})
}

func (c *Controller) Run(controlSimulation func(*Controller)) {
	for i := 0; i < c.simulation.ElevatorCount; i++ {
		elevator := NewElevator(i, c.simulation.ElevatorCapacity)
		c.elevators[elevator.Id] = elevator
		go elevator.Run(c)
	}
	var group sync.WaitGroup
	group.Add(c.simulation.PersonCount)
	defer group.Wait()
	go c.generatePersons(&group)
	go controlSimulation(c)
}

func (c *Controller) generatePersons(group *sync.WaitGroup) {
	for i := 0; i < c.simulation.PersonCount; i++ {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Microsecond)
		person := NewRandomPerson(i, c.simulation.FloorCount)
		c.persons.Store(person.Id, person)
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
