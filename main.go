package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	var pomo int
	flag.IntVar(&pomo, "length", 25, "pomodoro length in minutes")
	flag.Parse()
	done := make(chan struct{})

	fmt.Println(pomo)

	scanner := bufio.NewScanner(os.Stdin)

	// command line interface
	go func() {
		for {
			fmt.Printf("> ")
			scanner.Scan()
			text := scanner.Text()
			fmt.Println("Your text:", text)
			if text == "q" {
				done <- struct{}{}
			}
		}
	}()

	<-done
}
