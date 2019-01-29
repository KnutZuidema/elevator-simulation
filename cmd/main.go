package main

import (
	"github.com/KnutZuidema/elevator-simulation/pkg/logic"
	"github.com/KnutZuidema/elevator-simulation/pkg/model"
	"log"
	"os"
)

const DEBUG = true

//noinspection ALL
func main() {
	if !DEBUG {
		file, err := os.OpenFile("simulation.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			panic(err)
		}
		log.SetOutput(file)
	}
	simulation := model.NewSimulation(100, 25, 10, 10000)
	simulation.ControllingAlgorithms["simple"] = logic.ControlSimulationSimple
	simulation.Run()
}
