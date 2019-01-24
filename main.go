package main

import (
	"elevator_simulation/model"
	"log"
	"os"
)

func main() {
	file, err := os.OpenFile("simulation.log", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	log.SetOutput(file)
	elevator := model.NewElevator()
}
