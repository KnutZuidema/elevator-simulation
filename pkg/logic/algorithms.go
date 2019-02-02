package logic

import (
	"github.com/KnutZuidema/elevator-simulation/pkg/model"
	"time"
)

type RequestMap map[int][]bool
type DestinationMap map[int][]int

// newRequestMap returns a new map of requests where each request is not active
func newRequestMap(elevatorCount, floorCount int) (requests RequestMap) {
	requests = RequestMap{}
	for i := 0; i < elevatorCount; i++ {
		requests[i] = make([]bool, floorCount)
	}
	return
}

// newRequestMap returns a new map of requests where each request is not active
func newDestinationMap(elevatorCount, elevatorCapacity int) (destinations DestinationMap) {
	destinations = DestinationMap{}
	for i := 0; i < elevatorCount; i++ {
		destinations[i] = make([]int, elevatorCapacity)
	}
	return
}

// Simulation controlling algorithm
// move towards closest requested floor
func ControlSimulationSimple(controller *model.Controller) {
	destinations := newRequestMap(len(controller.Elevators), len(controller.Floors))
	ticker := time.NewTicker(time.Millisecond)
	for range ticker.C {
		updatePersons(controller)
		checkFloors(controller, destinations)
		signalElevators(controller.Elevators, destinations)
	}
}

func ControlSimulationMostRequestedDestination(controller *model.Controller) {
	floorRequests := newRequestMap(len(controller.Elevators), len(controller.Floors))
	passengerRequests := newDestinationMap(len(controller.Elevators), controller.Simulation.ElevatorCapacity)
	ticker := time.NewTicker(time.Millisecond)
	for range ticker.C {
		updatePersons(controller)
		checkFloors(controller, floorRequests)
		signalElevatorsMostRequestedDestinations(controller, passengerRequests, floorRequests)
	}
}

// checkFloors requests the closest elevator with remaining capacity for each floor with active requests
func checkFloors(controller *model.Controller, floorRequests RequestMap) {
	for index := range getFloorsWithRequests(controller.Floors) {
		closest := getClosestElevator(getElevatorsWithRemainingCapacity(controller.Elevators), index)
		if closest != nil {
			floorRequests[closest.Id][index] = true
		}
	}
}

// getFloorsWithRequests returns all floors with active requests
func getFloorsWithRequests(floors model.FloorMap) (result model.FloorMap) {
	result = model.FloorMap{}
	for index, floor := range floors {
		if len(floor) > 0 {
			result[index] = floor
		}
	}
	return
}

// getElevatorsWithRemainingCapacity returns all elevators with remaining capacity
func getElevatorsWithRemainingCapacity(elevators model.ElevatorMap) (available model.ElevatorMap) {
	available = model.ElevatorMap{}
	for index, elevator := range elevators {
		if elevator.Capacity > len(elevator.Persons) {
			available[index] = elevator
		}
	}
	return
}

// getClosestElevator returns the elevator which is closest to the specified floor
// if elevators is an empty slice returns nil
func getClosestElevator(elevators model.ElevatorMap, floor int) (closest *model.Elevator) {
	closest = nil
	for _, elevator := range elevators {
		if closest == nil {
			closest = elevator
		} else if absolute(elevator.CurrentFloor-floor) < absolute(closest.CurrentFloor-floor) {
			closest = elevator
		}
	}
	return
}

// getClosestRequest returns the floor id of the closest active request
// if there are no requests returns -1
func getClosestRequest(requests RequestMap, elevator *model.Elevator) (closest int) {
	closest = -1
	for index, request := range requests[elevator.Id] {
		if request {
			if closest == -1 {
				closest = index
			} else if absolute(elevator.CurrentFloor-index) < absolute(closest-index) {
				closest = index
			}
		}
	}
	return
}

func signalElevatorsMostRequestedDestinations(controller *model.Controller, passengerRequests DestinationMap, floorRequests RequestMap) {
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

// maxIndex returns the index of the maximum value in slice
func maxIndex(slice []int) int {
	result := 0
	for index, element := range slice {
		if element > slice[result] {
			result = index
		}
	}
	return result
}

func signalElevators(elevators map[int]*model.Elevator, requests RequestMap) {
	for _, elevator := range elevators {
		if requests[elevator.Id][elevator.CurrentFloor] {
			elevator.OpenDoorsSignal <- nil
			<-elevator.ContinueSignal
			requests[elevator.Id] = make([]bool, len(requests[elevator.Id]))
			for _, person := range elevator.Persons {
				requests[elevator.Id][person.DestinationFloor] = true
			}
		} else {
			closestRequest := getClosestRequest(requests, elevator)
			if closestRequest == -1 {
				elevator.IdleSignal <- nil
			} else if closestRequest < elevator.CurrentFloor {
				elevator.DescendSignal <- nil
			} else if closestRequest > elevator.CurrentFloor {
				elevator.AscendSignal <- nil
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
