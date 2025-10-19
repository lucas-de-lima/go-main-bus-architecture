package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Event represents an item transported on the conveyors (e.g., iron plate)
type Event struct {
	ID       int
	Resource string
	Value    any
	Time     time.Time
}

// Conveyor represents a single belt (a channel)
type Conveyor chan Event

// MainBus represents a resource-specific main bus with multiple parallel conveyors
type MainBus struct {
	Resource  string
	Conveyors []Conveyor
}

// NewMainBus creates a new main bus for a given resource with N parallel conveyors and a buffer size
func NewMainBus(resource string, lines int, buffer int) *MainBus {
	if lines%2 != 0 {
		lines++ // ensure even number of conveyors, following Factorio convention
	}
	bus := &MainBus{Resource: resource}
	for i := 0; i < lines; i++ {
		bus.Conveyors = append(bus.Conveyors, make(Conveyor, buffer))
	}
	return bus
}

// Produce sends an event to a random conveyor on the main bus
func (bus *MainBus) Produce(ev Event) {
	if len(bus.Conveyors) == 0 {
		return
	}
	idx := rand.Intn(len(bus.Conveyors))
	bus.Conveyors[idx] <- ev
}

// Consume starts consuming a specific conveyor until it is closed
func (bus *MainBus) Consume(line int, wg *sync.WaitGroup) {
	defer wg.Done()
	for ev := range bus.Conveyors[line] {
		fmt.Printf("[Consumer-%s-L%d] ID:%d Value:%v Time:%s\n", bus.Resource, line, ev.ID, ev.Value, ev.Time.Format("15:04:05"))
	}
}

// Close closes all conveyors in the main bus
func (bus *MainBus) Close() {
	for _, c := range bus.Conveyors {
		close(c)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Create two independent resource main buses
	ironBus := NewMainBus("iron", 4, 20)
	copperBus := NewMainBus("copper", 2, 20)

	var wg sync.WaitGroup

	// Start consumers for each iron conveyor
	for i := range ironBus.Conveyors {
		wg.Add(1)
		go ironBus.Consume(i, &wg)
	}

	// Start consumers for each copper conveyor
	for i := range copperBus.Conveyors {
		wg.Add(1)
		go copperBus.Consume(i, &wg)
	}

	// Producers sending events
	for i := 0; i < 10; i++ {
		ironBus.Produce(Event{ID: i, Resource: "iron", Value: "Iron Plate", Time: time.Now()})
		copperBus.Produce(Event{ID: i, Resource: "copper", Value: "Copper Plate", Time: time.Now()})
		time.Sleep(100 * time.Millisecond)
	}

	// Finish
	time.Sleep(1 * time.Second)
	ironBus.Close()
	copperBus.Close()

	wg.Wait()
	fmt.Println("All main buses completed processing.")
}
