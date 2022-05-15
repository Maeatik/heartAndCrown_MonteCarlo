package main

import (
	"flag"
	"fmt"
	runtime "runtime"
)

// main function
func main() {

// command line args

	numCPU := runtime.NumCPU()
	runtime := flag.Int("runtime", 1, "run time in minutes")

	cores := flag.Int("cores", numCPU, "number of cores")

	flag.Parse()

// compute percent winnings using Monte Carlo
	pctPlayer, pctBank, totaltrials := MonteCA(*cores, *runtime)

	fmt.Printf("Total trials using %d cores for %d minutes: %d\n", *cores, *runtime, totaltrials)
	fmt.Printf("Percentage for player: %5.2f\n", pctPlayer)
	fmt.Printf("Percentage for bank: %5.2f\n", pctBank)
}