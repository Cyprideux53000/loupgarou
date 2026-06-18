package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

// func printNumbers() {
// 	for i := range 5 {
// 		fmt.Println(i)
// 	}
// }
// func printLetters() {
// 	for char := 'A'; char <= 'E'; char++ {
// 		fmt.Println(string(char))
// 	}
// }
// func main() {
// 	go printNumbers()
// 	go printLetters()
// 	time.Sleep(time.Second)
// }

//	func takeCoffees(ch chan<- int, coffeNumber int) {
//		ch <- coffeNumber * 2
//	}
//
//	func deliverCoffees(ch <-chan int) {
//		delivred := <-ch
//		fmt.Println("Delivered coffee: ", delivred)
//	}
//
//	func main() {
//		ch := make(chan int)
//		go takeCoffees(ch, 2)
//		deliverCoffees(ch)
//	}
func printMessage(message string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(message)
}
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <repetition>")
		os.Exit(1)
	}
	repetition, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Error: repetition must be a number")
		os.Exit(1)
	}
	var wg sync.WaitGroup

	for i := 0; i < repetition; i++ {
		wg.Add(2)
		go printMessage("Hello", &wg)
		go printMessage("World", &wg)

	}
	wg.Wait()
}
