package main

import "fmt"

import "flag"
import "strconv"
import "os"
import "runtime/pprof"

func rec_sum(slice []uint64, c chan uint64) {
	if len(slice) < 10 {
		sum := uint64(0)
		for i := 0; i < len(slice); i++ {
			sum += slice[i]
		}
		c <- sum
	} else {
		c0 := make(chan uint64)
		c1 := make(chan uint64)
		pivot := len(slice) / 2
		go rec_sum(slice[0:pivot], c0)
		go rec_sum(slice[pivot:], c1)
		s0 := <-c0
		s1 := <-c1
		c <- (s0 + s1)
	}
}

func fork_join_add(slice []uint64) uint64 {
	c := make(chan uint64)
	go rec_sum(slice, c)
	return <-c
}

var profile = flag.String("prof", "", "write profile to file")

func main() {
	flag.Parse()
	args := flag.Args()
	if *profile != ""  {
        f, err := os.Create(*profile)
        if err != nil {
            fmt.Printf("Could not open %s\n", *profile)
            os.Exit(1)
        }
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }
	if len(args) == 0 {
		fmt.Printf("Not enough arguments")
		os.Exit(1)
	}
	elements, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("Could not convert %s to int", args[0])
		os.Exit(1)
	}
	array := make([]uint64, elements)
	for i := 0; i < len(array); i++ {
		array[i] = uint64(i)
	}
	fmt.Printf("Generated %d ints\n", len(array))
	fmt.Printf("Sum is %d", fork_join_add(array))
}
