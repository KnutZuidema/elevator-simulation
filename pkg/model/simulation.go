package model

import (
	"fmt"
	"log"
	"os"
)

type Simulation struct {
	FloorCount            int
	ElevatorCount         int
	ElevatorCapacity      int
	PersonCount           int
	ControllingAlgorithms map[string]func(*Controller)
}

func NewSimulation(floorCount int, elevatorCount int, elevatorCapacity int, personCount int) *Simulation {
	return &Simulation{
		FloorCount:            floorCount,
		ElevatorCount:         elevatorCount,
		ElevatorCapacity:      elevatorCapacity,
		PersonCount:           personCount,
		ControllingAlgorithms: map[string]func(*Controller){},
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
	for name, algorithm := range s.ControllingAlgorithms {
		controller := NewController(s)
		controller.Run(algorithm)
		s.Evaluate(fmt.Sprintf("%v_report.txt", name), controller)
	}
}
