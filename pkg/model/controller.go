package model

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"
)

type ElevatorMap map[int]*Elevator

// NewElevatorMap returns a new map of elevators with the specified capacity
// elevators are set initially at floor 0
func NewElevatorMap(elevatorCount, elevatorCapacity int) (elevators ElevatorMap) {
	elevators = ElevatorMap{}
	for i := 0; i < elevatorCount; i++ {
		elevators[i] = NewElevator(i, elevatorCapacity)
	}
	return
}

type Floor chan *Person

type FloorMap map[int]Floor

// NewFloor returns a new map of floors with empty channels
func NewFloorMap(floorCount, elevatorCapacity int) (floors FloorMap) {
	floors = FloorMap{}
	for i := 0; i < floorCount; i++ {
		floors[i] = make(chan *Person, elevatorCapacity)
	}
	return
}

type Controller struct {
	Simulation *Simulation
	Floors     FloorMap
	Persons    *sync.Map
	Elevators  ElevatorMap
}

func NewController(simulation *Simulation) *Controller {
	return &Controller{
		Simulation: simulation,
		Floors:     NewFloorMap(simulation.FloorCount, simulation.ElevatorCapacity),
		Persons:    &sync.Map{},
		Elevators:  NewElevatorMap(simulation.ElevatorCount, simulation.ElevatorCapacity),
	}
}

func (c *Controller) Evaluate(file io.Writer) {
	if _, err := fmt.Fprintln(file, "Elevators"); err != nil {
		log.Fatal(err)
		return
	}
	for _, elevator := range c.Elevators {
		_, err := fmt.Fprintf(file,
			"  Elevator %02v:\n"+
				"    %v steps\n"+
				"    %v picked up\n"+
				"    %v floors traveled\n"+
				"    %v doors opened\n"+
				"    %v times idled\n", elevator.Id, elevator.StepsTaken, elevator.TotalPickedUp,
			elevator.FloorsTraveledCount, elevator.DoorsOpenedCount, elevator.IdleCount)
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

// Run runs the simulation until all persons reached their destination
func (c *Controller) Run(controlSimulation func(*Controller)) {
	for _, elevator := range c.Elevators {
		go elevator.Run(c)
	}
	var group sync.WaitGroup
	group.Add(c.Simulation.PersonCount)
	defer group.Wait()
	go c.generatePersons(&group)
	go controlSimulation(c)
}

// generatePersons adds new persons in random intervals to the controller
func (c *Controller) generatePersons(group *sync.WaitGroup) {
	for i := 0; i < c.Simulation.PersonCount; i++ {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Microsecond)
		person := NewRandomPerson(i, c.Simulation.FloorCount)
		c.Persons.Store(person.Id, person)
		go person.Run(c, group)
	}
}

// meanTravelTime returns the mean travel time of all persons
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

// meanWaitingTime returns the mean waiting time of all persons
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
