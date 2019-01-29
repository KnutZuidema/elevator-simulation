package model

import (
	"github.com/KnutZuidema/elevator-simulation/pkg/logic"
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
	controller.Run(logic.ControlSimulationSimple)
	s.Evaluate("final_report.txt", controller)
}
