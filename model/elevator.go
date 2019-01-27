package model

import "log"

type Elevator struct {
	Id               int
	Capacity         int
	Persons          []*Person
	CurrentFloor     int
	DestinationFloor int
	StepsTaken       int
	OpenDoorsSignal  chan interface{}
	AscendSignal     chan interface{}
	DescendSignal    chan interface{}
	ContinueSignal   chan interface{}
}

func NewElevator(id int, capacity int) *Elevator {
	return &Elevator{
		Id:              id,
		Capacity:        capacity,
		OpenDoorsSignal: make(chan interface{}),
		AscendSignal:    make(chan interface{}),
		DescendSignal:   make(chan interface{}),
		ContinueSignal:  make(chan interface{}),
	}
}

func (e *Elevator) Run(controller *Controller) {
	for {
		select {
		case <-e.OpenDoorsSignal:
			e.pickUp(controller.Floors[e.CurrentFloor])
			e.dropOff(e.CurrentFloor)
			if len(e.Persons) > 0 && e.CurrentFloor != e.DestinationFloor {
				e.DestinationFloor = e.Persons[0].DestinationFloor
				log.Printf("elevator %03v: new destination is floor %v (internal signal)", e.Id, e.DestinationFloor)
			}
		case <-e.AscendSignal:
			log.Printf("elevator %03v: ascending from floor %v", e.Id, e.CurrentFloor)
			e.CurrentFloor++
		case <-e.DescendSignal:
			log.Printf("elevator %03v: descending from floor %v", e.Id, e.CurrentFloor)
			e.CurrentFloor--
		}
		e.StepsTaken++
		e.ContinueSignal <- nil
	}
}

func (e *Elevator) pickUp(floor chan *Person) {
loop:
	for e.Capacity-len(e.Persons) > 0 {
		select {
		case person := <-floor:
			close(person.PickUpSignal)
			e.Persons = append(e.Persons, person)
		default:
			break loop
		}
	}
}

func (e *Elevator) dropOff(floor int) {
	var newPersons []*Person
	for _, person := range e.Persons {
		if person.DestinationFloor == e.CurrentFloor {
			close(person.DropOffSignal)
		} else {
			newPersons = append(newPersons, person)
		}
	}
	e.Persons = newPersons
}
