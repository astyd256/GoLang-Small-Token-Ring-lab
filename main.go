package main

import (
	"fmt"
	"sync"
)

type token struct {
	data      string
	recipient int
	ttl       int
}

var wg sync.WaitGroup

func node(thread_index int, ch_resv chan token, ch_send chan token) {
	for { // waiting for any signal
		message := <-ch_resv
		if message.data == "qqq" {
			ch_send <- message
			break
		}
		fmt.Println("Thread ", thread_index, "| Recieved message.")
		if message.recipient == thread_index {
			fmt.Println("Thread ", thread_index, "|", message.data)
		} else {
			if message.ttl > 0 {
				message.ttl -= 1
				ch_send <- message
			} else {
				fmt.Println("Thread ", thread_index, "| Message expired.")
			}
		}
	}
	wg.Done()
}

func main() {
	fmt.Println("MAIN THREAD| Enter number of threads")
	var threads_number int
	fmt.Scanln(&threads_number)

	// creating channels, 0 channel is reserved for talking with main thread
	var channels []chan token = make([]chan token, threads_number)
	var quit = make(chan int)

	for i := 0; i < threads_number; i++ { // creating threads
		channels[i] = make(chan token)
		wg.Add(1)
	}

	for i := 0; i < threads_number; i++ {
		if i == threads_number-1 {
			go node(i, channels[i], channels[i]) // sending for last channel same channel bcs he never send msg to anyone
		} else {
			go node(i, channels[i], channels[i+1])
		}
	}

	for threads_number > 0 {

		var new_token token

		fmt.Println("MAIN THREAD| Enter message data or write qqq to exit")
		fmt.Scanln(&new_token.data)
		if new_token.data == "qqq" {
			quit <- 0
			break
		}
		fmt.Println("MAIN THREAD| Enter recipient id (starting from 0)")
		fmt.Scanln(&new_token.recipient)
		for {
			if new_token.recipient > len(channels)-1 {
				fmt.Println("MAIN THREAD| This is now a valid recipient id! Try again! Avaliable ids: from 0 to ", len(channels)-1)
				fmt.Scanln(&new_token.recipient)
			} else {
				break
			}
		}
		fmt.Println("MAIN THREAD| Enter message timeout")
		fmt.Scanln(&new_token.ttl)

		channels[0] <- new_token

	}
	wg.Wait()
	for i := 1; i < threads_number; i++ { // creating threads
		close(channels[i])
	}
	close(quit)
}
