package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	closingTime = 30 * time.Second // Adjust closing time as needed
	numSeats    = 5                // Number of seats in the waiting room
	numBarbers  = 2                // Number of barbers
)

type Barber struct {
	id        int
	waiting   bool
	customers chan *Customer
	done      chan struct{}
}

type Customer struct {
	id int
}

type BarberShop struct {
	waitingRoom chan *Customer
	barbers     []*Barber
	wg          sync.WaitGroup
}

func NewBarber(id int, customers chan *Customer) *Barber {
	return &Barber{
		id:        id,
		waiting:   false,
		customers: customers,
		done:      make(chan struct{}),
	}
}

func (b *Barber) run(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case customer := <-b.customers:
			b.waiting = false
			fmt.Printf("Barber %d is cutting hair for Customer %d\n", b.id, customer.id)
			time.Sleep(time.Second * 2) // Simulate haircut time
			fmt.Printf("Barber %d finished cutting hair for Customer %d\n", b.id, customer.id)
			b.waiting = true
		case <-b.done:
			return
		}
	}
}

func NewBarberShop() *BarberShop {
	customers := make(chan *Customer, numSeats)
	barbers := make([]*Barber, numBarbers)
	for i := 0; i < numBarbers; i++ {
		barbers[i] = NewBarber(i, customers)
	}
	return &BarberShop{
		waitingRoom: customers,
		barbers:     barbers,
	}
}

func (bs *BarberShop) Open() {
	var wg sync.WaitGroup
	wg.Add(numBarbers)
	for _, barber := range bs.barbers {
		go barber.run(&wg)
	}

	go func() {
		wg.Wait()
		for _, barber := range bs.barbers {
			close(barber.done)
		}
	}()

	startTime := time.Now()
	for {
		if time.Since(startTime) > closingTime {
			break
		}
		time.Sleep(time.Second)                        // Simulate time passing
		customer := &Customer{id: time.Now().Second()} // Generate a customer with a unique ID
		select {
		case bs.waitingRoom <- customer:
			fmt.Printf("Customer %d has entered the shop\n", customer.id)
		default:
			fmt.Printf("Customer %d left because all chairs are occupied\n", customer.id)
		}
	}

	close(bs.waitingRoom)
	fmt.Println("Shop is closed, waiting for all customers to finish")
	bs.wg.Wait()
	fmt.Println("All customers have left, barbers are going home")
}

func main() {
	shop := NewBarberShop()
	shop.Open()
}
