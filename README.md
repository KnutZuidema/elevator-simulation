# elevator-simulation

Simulate a building with multiple elevators and floors carrying passengers to their destinations.  
Used as a golang learning exercise.

## What I used

* basic language constructs:
  * structs, channels, interfaces, go routines etc.
* `math/rand` package for generation of new passengers
* `time` package for `time.Ticker` functionality (simulate travel time of elevators)
* `sync` package for coordination of program flow
  * `sync.WaitGroup` to determine when program should terminate (i.e. all passengers have reached their destination)
  * `sync.Map` to enable concurrent read/write access to passengers
* `log` package for simple logging
* `os` package for creating and writing to files