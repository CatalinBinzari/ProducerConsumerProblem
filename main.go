package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var pizzaMade, pizzaFailed, total int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error // channel or errors
}

func (p *Producer) Close() error {
	// close the  quit channel
	ch := make(chan error)
	p.quit <- ch
	return <-ch // will be nill if we close the channel successfuly

}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++

	if pizzaNumber <= NumberOfPizzas {
		delay := rand.Intn(5) + 1
		fmt.Printf("received order #%d!\n", pizzaNumber)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		fmt.Printf("making pizza %d. it will take %d seconds\n", pizzaNumber, delay)
		// delay for a bit

		time.Sleep(time.Duration(delay) * time.Second)
		if rnd <= 2 {
			msg = fmt.Sprintf("we run out of ingredients #%d", pizzaNumber)
			pizzaFailed++
		} else if rnd <= 4 {
			msg = fmt.Sprintf("the cook quit #%d", pizzaNumber)
			pizzaFailed++
		} else {
			success = true
			msg = fmt.Sprintf("pizza is ready #%d", pizzaNumber)
			pizzaMade++
		}
		total++

		p := PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}

		return &p

	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}
}

func pizzaria(pizzaMaker *Producer) {
	// keep track of the which of pizzas made
	var i = 0

	// run forever or until we receive a quit notification

	//try to make pizza

	for {
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.pizzaNumber
			select {
			// try to make a pizza, we send smth to the datachannel
			case pizzaMaker.data <- *currentPizza:

			case quitChan := <-pizzaMaker.quit:
				//close channel
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		}
	}
}

// we have a pizzaria - make pizza - producer
// we have costumers - place orders - consumers
func main() {
	// seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// print out a message
	color.Cyan("Starting the pizza shop")

	// create producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	// run the producer in the background -> goroutine
	go pizzaria(pizzaJob)

	// create consumers
	for i := range pizzaJob.data {
		if i.pizzaNumber <= NumberOfPizzas {
			if i.success {
				color.Green(i.message)
				color.Green("Order #%d is ready", i.pizzaNumber)
			} else {
				color.Red(i.message)
				color.Green("costumer did not receive its order #%d", i.pizzaNumber)
			}
		} else {
			color.Cyan("done making pizzas")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("error closing the pizza job channel ", err)
			}
		}
	}
	// each time  the consumer places a order to the pizzaria

	// print out the ending message

	color.Cyan("-----------------")
	color.Cyan("Done for today")

	color.Cyan("Total pizzas made: %d, failed %d, attemts %d.", pizzaMade, pizzaFailed, total)

	switch {
	case pizzaFailed > 9:
		color.Red("It was an awful day")
	case pizzaFailed > 6:
		color.Red("It was not a very good day")
	case pizzaFailed >= 4:
		color.Yellow("It was a good day")
	case pizzaFailed >= 2:
		color.Yellow("It was a great day")
	}
}
