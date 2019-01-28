package main

import (
	"github.com/KnutZuidema/elevator-simulation/model"
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
	model.NewSimulation(100, 25, 10, 10000).Run()
}
