package main

//http://golangtutorials.blogspot.jp/2011/06/goroutines.html

import (
	"fmt"
	"time"
)

func simulateEvent(name string, timeInSecs int) {
	// sleep for a while to simulate time consumed by event
	fmt.Println("Started ", name, ": Should take", timeInSecs, "seconds.")
	time.Sleep(time.Duration(timeInSecs) * time.Second) //warn:these style convert
	fmt.Println("Finished ", name, "my time is ", time.Duration(timeInSecs)*time.Second)
}

func main() {
	go simulateEvent("100m sprint", 10) //start 100m sprint, it should take 10 seconds
	go simulateEvent("Long jump", 6)    //start long jump, it should take 6 seconds
	go simulateEvent("High jump", 3)    //start high jump, it should take 3 seconds

	//so that the program doesn't exit here, we make the program wait here for a while
	time.Sleep(12 * time.Second)
}
