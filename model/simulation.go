package model

import (
	"fmt"
	"log"
	"os"
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
	var iteration int
	for {
		iteration++
		updatePersons(controller)
		signalElevators(controller)
		checkFloors(controller)
		if iteration%10000 == 0 {
			controller.simulation.Evaluate(fmt.Sprintf("reports/%v-report.txt", iteration), controller)
		}
	}
}

func checkFloors(controller *Controller) {
floor:
	for i, floor := range controller.Floors {
		if len(floor) > 0 {
			for _, elevator := range controller.elevators {
				if elevator.CurrentFloor == elevator.DestinationFloor {
					elevator.DestinationFloor = i
					log.Printf("elevator %03v: new destination is %v (external signal)",
						elevator.Id, elevator.DestinationFloor)
					continue floor
				}
			}
		}
	}
}

func signalElevators(controller *Controller) {
	for _, elevator := range controller.elevators {
		if elevator.CurrentFloor < elevator.DestinationFloor {
			elevator.AscendSignal <- nil
		} else if elevator.CurrentFloor > elevator.DestinationFloor {
			elevator.DescendSignal <- nil
		} else {
			elevator.OpenDoorsSignal <- nil
		}
		<-elevator.ContinueSignal
	}
}

func updatePersons(controller *Controller) {
	for _, person := range controller.persons {
		if person.IsWaiting {
			person.WaitingTime++
		} else if person.IsTraveling {
			person.TravelTime++
		}
	}
}
