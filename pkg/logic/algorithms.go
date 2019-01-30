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

func ControlSimulationMostRequestedDestination(controller *model.Controller) {
	ticker := time.NewTicker(time.Millisecond)
	passengerRequests := map[int][]int{}
	floorRequests := map[int][]bool{}
	for _, elevator := range controller.Elevators {
		passengerRequests[elevator.Id] = make([]int, len(controller.Floors))
		floorRequests[elevator.Id] = make([]bool, len(controller.Floors))
	}
	for range ticker.C {
		updatePersons(controller)
		checkFloorsMostRequestedDestinations(controller, floorRequests)
		signalElevatorsMostRequestedDestinations(controller, passengerRequests, floorRequests)
	}
}

func checkFloorsMostRequestedDestinations(controller *model.Controller, floorRequests map[int][]bool) {
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
				floorRequests[closest.Id][i] = true
			}
		}
	}
}

func signalElevatorsMostRequestedDestinations(controller *model.Controller, passengerRequests map[int][]int, floorRequests map[int][]bool) {
	for _, elevator := range controller.Elevators {
		// if there is a request from the current floor or passengers requested the current floor open doors
		// flush and reinitialize passenger request counter
		if passengerRequests[elevator.Id][elevator.CurrentFloor] > 0 || floorRequests[elevator.Id][elevator.CurrentFloor] {
			elevator.OpenDoorsSignal <- nil
			<-elevator.ContinueSignal
			floorRequests[elevator.Id][elevator.CurrentFloor] = false
			passengerRequests[elevator.Id] = make([]int, len(controller.Floors))
			for _, person := range elevator.Persons {
				passengerRequests[elevator.Id][person.DestinationFloor]++
			}
		} else {
			// if there are passengers find the most requested floor and move elevator towards it
			// if there are no passengers move elevator towards closest floor with active request
			// if there are no active requests do nothing
			nextFloor := elevator.CurrentFloor
			if len(elevator.Persons) > 0 {
				nextFloor = maxIndex(passengerRequests[elevator.Id])
			} else {
				for index, floor := range floorRequests[elevator.Id] {
					if floor {
						if nextFloor == elevator.CurrentFloor {
							nextFloor = index
						} else if absolute(elevator.CurrentFloor-index) < absolute(elevator.CurrentFloor-nextFloor) {
							nextFloor = index
						}
					}
				}
			}
			if nextFloor != elevator.CurrentFloor {
				if elevator.CurrentFloor > nextFloor {
					elevator.DescendSignal <- nil
				} else {
					elevator.AscendSignal <- nil
				}
			} else {
				elevator.IdleSignal <- nil
			}
			<-elevator.ContinueSignal
		}
	}
}

// return index of the maximum value in slice
func maxIndex(slice []int) int {
	result := 0
	for index, element := range slice {
		if element > slice[result] {
			result = index
		}
	}
	return result
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

// return absolute value of input
func absolute(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
