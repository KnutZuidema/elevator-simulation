package model

type Elevator struct {
	Capacity      int
	Persons       []*Person
	CurrentFloor  int
	StepsTaken    int
	PickUpSignal  chan chan *Person
	DropOffSignal chan int
	AscendSignal  chan interface{}
	DescendSignal chan interface{}
}

func NewElevator(capacity int) *Elevator {
	return &Elevator{
		Capacity:      capacity,
		PickUpSignal:  make(chan chan *Person),
		DropOffSignal: make(chan int),
		AscendSignal:  make(chan interface{}),
		DescendSignal: make(chan interface{}),
	}
}

func (e *Elevator) Run(controller *Controller) {
	for {
		select {
		case floor := <-e.PickUpSignal:
			e.pickUp(floor)
		case floor := <-e.DropOffSignal:
			e.dropOff(floor)
		case <-e.AscendSignal:
			e.CurrentFloor++
		case <-e.DescendSignal:
			e.CurrentFloor--
		}
		e.StepsTaken++
	}
}

func (e *Elevator) pickUp(floor chan *Person) {
	for e.Capacity-len(e.Persons) > 0 {
		select {
		case person := <-floor:
			close(person.PickUpSignal)
			e.Persons = append(e.Persons, person)
		default:
			break
		}
	}
}

func (e *Elevator) dropOff(floor int) {
	for i, person := range e.Persons {
		if person.DestinationFloor == e.CurrentFloor {
			close(person.DropOffSignal)
			e.Persons = append(e.Persons[:i], e.Persons[i+1:]...)
		}
	}
}
