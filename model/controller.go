package model

import (
	"errors"
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
	persons    []*Person
	elevators  []*Elevator
}

func NewController(simulation *Simulation) *Controller {
	floors := make([]chan *Person, simulation.FloorCount)
	for i := range floors {
		floors[i] = make(chan *Person, simulation.ElevatorCapacity)
	}
	return &Controller{
		simulation: simulation,
		Floors:     floors,
		persons:    []*Person{},
	}
}

func (c *Controller) Evaluate(file io.Writer) {
	if _, err := fmt.Fprintln(file, "Elevators"); err != nil {
		log.Fatal(err)
		return
	}
	for _, elevator := range c.elevators {
		_, err := fmt.Fprintf(file, "Elevator %v: %v steps\n", elevator.Id, elevator.StepsTaken)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	waitingPersons, _ := all(c.persons, func(i int) bool {
		return c.persons[i].WaitingTime > 0
	})
	var temp interface{}
	var meanWaitingTime float32
	temp = waitingPersons
	switch w := temp.(type) {
	case []*Person:
		meanWaitingTime, _ = mean(w, func(i int) int {
			return w[i].WaitingTime
		})
	default:

	}
	meanTravelingTime, _ := mean(c.persons, func(i int) int {
		return c.persons[i].TravelTime
	})
	_, err := fmt.Fprintf(file, "Persons: mean waiting time %v, mean traveling time %v\n",
		meanWaitingTime, meanTravelingTime)
	if err != nil {
		log.Fatal(err)
		return
	}
	for _, person := range c.persons {
		_, err := fmt.Fprintf(file, "Person %v: waited %v steps, travelled %v steps\n",
			person.Id, person.WaitingTime, person.TravelTime)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func (c *Controller) Run(controlSimulation func(*Controller)) {
	for i := 0; i < c.simulation.ElevatorCount; i++ {
		elevator := NewElevator(i, c.simulation.ElevatorCapacity)
		c.elevators = append(c.elevators, elevator)
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
		time.Sleep(time.Duration(10+rand.Intn(20)) * time.Millisecond)
		person := NewRandomPerson(i, c.simulation.FloorCount)
		c.persons = append(c.persons, person)
		go person.Run(c, group)
	}
}

func total(slice interface{}, key func(int) int) (total int, err error) {
	switch s := slice.(type) {
	case []interface{}:
		for i := range s {
			total += key(i)
		}
	default:
		err = errors.New("invalid type, must be type []interface{}")
	}
	return
}

func mean(slice interface{}, key func(int) int) (mean float32, err error) {
	total, err := total(slice, key)
	if err != nil {
		return
	}
	switch s := slice.(type) {
	case []interface{}:
		mean = float32(total) / float32(len(s))
	default:
		err = errors.New("invalid type, must be type []interface{}")
	}
	return
}

func all(slice interface{}, key func(int) bool) (result []interface{}, err error) {
	switch s := slice.(type) {
	case []interface{}:
		for i := range s {
			if key(i) {
				result = append(result, s[i])
			}
		}
	default:
		err = errors.New("invalid type, must be type []interface{}")
	}
	return
}
