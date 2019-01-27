package main

import (
	"github.com/KnutZuidema/elevator-simulation/model"
	"log"
	"os"
)

const DEBUG = true

func main() {
	if !DEBUG {
		file, err := os.OpenFile("simulation.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			panic(err)
		}
		log.SetOutput(file)
	}
	model.NewSimulation(100, 25, 12, 1200).Run()
}
