package main

import (
	"flag"
	"fmt"
)

func main() {
	var pomo int
	flag.IntVar(&pomo, "length", 25, "pomodoro length in minutes")
	flag.Parse()

	fmt.Println(pomo)
}

func Timer(minutes int) {
}
