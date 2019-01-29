package model

import "log"

type Elevator struct {
	Id              int
	TotalPickedUp   int
	Capacity        int
	Persons         map[int]*Person
	CurrentFloor    int
	Destinations    map[int]bool
	StepsTaken      int
	OpenDoorsSignal chan interface{}
	AscendSignal    chan interface{}
	DescendSignal   chan interface{}
	IdleSignal      chan interface{}
	ContinueSignal  chan interface{}
}

func NewElevator(id int, capacity int) *Elevator {
	return &Elevator{
		Id:              id,
		Capacity:        capacity,
		Persons:         map[int]*Person{},
		Destinations:    map[int]bool{},
		OpenDoorsSignal: make(chan interface{}),
		AscendSignal:    make(chan interface{}),
		DescendSignal:   make(chan interface{}),
		IdleSignal:      make(chan interface{}),
		ContinueSignal:  make(chan interface{}),
	}
}

func (e *Elevator) Run(controller *Controller) {
	for {
		select {
		case <-e.OpenDoorsSignal:
			droppedOff := e.dropOff(e.CurrentFloor)
			log.Printf("elevator %04v: dropped off %v persons on floor %v", e.Id, droppedOff, e.CurrentFloor)
			pickedUp := e.pickUp(controller.Floors[e.CurrentFloor])
			log.Printf("elevator %04v: picked up %v persons on floor %v", e.Id, pickedUp, e.CurrentFloor)
			e.StepsTaken++
		case <-e.AscendSignal:
			// log.Printf("elevator %04v: ascending from floor %v", e.Id, e.CurrentFloor)
			e.CurrentFloor++
			e.StepsTaken++
		case <-e.DescendSignal:
			// log.Printf("elevator %04v: descending from floor %v", e.Id, e.CurrentFloor)
			e.CurrentFloor--
			e.StepsTaken++
		case <-e.IdleSignal:
		}
		e.ContinueSignal <- nil
	}
}

func (e *Elevator) pickUp(floor chan *Person) (count int) {
loop:
	for e.Capacity-len(e.Persons) > 0 {
		select {
		case person := <-floor:
			count++
			close(person.PickUpSignal)
			e.Persons[person.Id] = person
		default:
			break loop
		}
	}
	e.TotalPickedUp += count
	return
}

func (e *Elevator) dropOff(floor int) (count int) {
	for _, person := range e.Persons {
		if person.DestinationFloor == e.CurrentFloor {
			count++
			close(person.DropOffSignal)
			delete(e.Persons, person.Id)
		}
	}
	return
}
