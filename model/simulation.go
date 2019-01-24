package model

type Simulation struct {
	FloorCount       int
	ElevatorCount    int
	ElevatorCapacity int
	PersonCount      int
}

func NewSimulation(floorCount int, elevatorCount int, elevatorCapacity int, personCapacity int) *Simulation {
	return &Simulation{
		FloorCount:       floorCount,
		ElevatorCount:    elevatorCount,
		ElevatorCapacity: elevatorCapacity,
		PersonCount:      personCapacity,
	}
}

func (s *Simulation) Run() {
	NewController(s.FloorCount).Run(s)
}
