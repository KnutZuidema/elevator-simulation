package model

import (
	"log"
	"os"
	"time"
)

type Simulation struct {
	FloorCount       int
	ElevatorCount    int
	ElevatorCapacity int
	PersonCount      int
}

func NewSimulation(floorCount int, elevatorCount int, elevatorCapacity int, personCount int) *Simulation {
	return &Simulation{
		FloorCount:       floorCount,
		ElevatorCount:    elevatorCount,
		ElevatorCapacity: elevatorCapacity,
		PersonCount:      personCount,
	}
}

func (s *Simulation) Evaluate(filename string, controller *Controller) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatal(err)
		return
	}
	controller.Evaluate(file)
}

func (s *Simulation) Run() {
	controller := NewController(s)
	controller.Run(controlSimulationSimple)
	s.Evaluate("final_report.txt", controller)
}

func controlSimulationSimple(controller *Controller) {
	ticker := time.NewTicker(time.Millisecond)
	for range ticker.C {
		updatePersons(controller)
		checkFloors(controller)
		signalElevators(controller)
	}
}

func checkFloors(controller *Controller) {
	for i, floor := range controller.Floors {
		if len(floor) > 0 {
			var closest *Elevator
			for _, elevator := range controller.elevators {
				if elevator.Capacity > len(elevator.Persons) {
					if closest == nil {
						closest = elevator
					} else if absolute(elevator.CurrentFloor-i) < absolute(closest.CurrentFloor-i) {
						closest = elevator
					}
				}
			}
			if closest != nil {
				closest.Destinations[i] = true
			}
		}
	}
}

func signalElevators(controller *Controller) {
	for _, elevator := range controller.elevators {
		if elevator.Destinations[elevator.CurrentFloor] {
			elevator.OpenDoorsSignal <- nil
			<-elevator.ContinueSignal
			elevator.Destinations = map[int]bool{}
			for _, person := range elevator.Persons {
				elevator.Destinations[person.DestinationFloor] = true
			}
		} else {
			closestFloor := elevator.CurrentFloor
			for i := 1; elevator.CurrentFloor+i < controller.simulation.FloorCount || elevator.CurrentFloor-i >= 0; i++ {
				if elevator.Destinations[elevator.CurrentFloor+i] {
					closestFloor = elevator.CurrentFloor + i
					break
				}
				if elevator.Destinations[elevator.CurrentFloor-i] {
					closestFloor = elevator.CurrentFloor - i
					break
				}
			}
			if closestFloor < elevator.CurrentFloor {
				elevator.DescendSignal <- nil
			} else if closestFloor > elevator.CurrentFloor {
				elevator.AscendSignal <- nil
			} else {
				elevator.IdleSignal <- nil
			}
			<-elevator.ContinueSignal
		}
	}
}

func updatePersons(controller *Controller) {
	controller.persons.Range(func(_, value interface{}) bool {
		person := value.(*Person)
		if person.IsWaiting {
			person.WaitingTime++
		} else if person.IsTraveling {
			person.TravelTime++
		}
		return true
	})
}

func absolute(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
