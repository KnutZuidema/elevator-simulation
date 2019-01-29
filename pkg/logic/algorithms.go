package logic

import (
	"github.com/KnutZuidema/elevator-simulation/pkg/model"
	"time"
)

func ControlSimulationSimple(controller *model.Controller) {
	ticker := time.NewTicker(time.Millisecond)
	for range ticker.C {
		updatePersons(controller)
		checkFloors(controller)
		signalElevators(controller)
	}
}

func checkFloors(controller *model.Controller) {
	for i, floor := range controller.Floors {
		if len(floor) > 0 {
			var closest *model.Elevator
			for _, elevator := range controller.Elevators {
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

func signalElevators(controller *model.Controller) {
	for _, elevator := range controller.Elevators {
		if elevator.Destinations[elevator.CurrentFloor] {
			elevator.OpenDoorsSignal <- nil
			<-elevator.ContinueSignal
			elevator.Destinations = map[int]bool{}
			for _, person := range elevator.Persons {
				elevator.Destinations[person.DestinationFloor] = true
			}
		} else {
			closestFloor := elevator.CurrentFloor
			for i := 1; elevator.CurrentFloor+i < controller.Simulation.FloorCount || elevator.CurrentFloor-i >= 0; i++ {
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

func updatePersons(controller *model.Controller) {
	controller.Persons.Range(func(_, value interface{}) bool {
		person := value.(*model.Person)
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
