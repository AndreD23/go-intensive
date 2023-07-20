package main

import (
	"fmt"
	"time"
)

func processando() {
	for i := 0; i < 10; i++ {
		fmt.Println(i)
		time.Sleep(time.Second)
	}

}

// T1
func main() {
	//go processando() // T2
	//go processando() // T3
	//processando()

	canal := make(chan int)

	go func() {
		// canal <- 1 //T2

		for i := 0; i < 10; i++ {
			canal <- i
			fmt.Println("Jogou no canal ", i)
		}
	}()

	go func() {
		// canal <- 1 //T2

		for i := 0; i < 10; i++ {
			canal <- i
			fmt.Println("Jogou no canal ", i)
		}
	}()

	//go fmt.Println(<-canal) // esvazia o canal

	//for x := range canal {
	//	fmt.Println(x)
	//	fmt.Println("Recebeu do canal ", x)
	//	time.Sleep(time.Second)
	//}

	go worker(canal, 1)
	worker(canal, 2)
}

func worker(canal chan int, workerID int) {
	for {
		fmt.Println("Recebeu do canal ", <-canal, "no worker", workerID)
		time.Sleep(time.Second)
	}
}
