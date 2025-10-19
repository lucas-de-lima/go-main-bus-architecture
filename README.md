# Main Bus Architecture in Go (Inspired by Factorio)

[![Go Version](https://img.shields.io/badge/go-1.20+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg?style=for-the-badge)](LICENSE)
[![Factorio](https://img.shields.io/badge/inspired%20by-Factorio-orange?style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjQiIGhlaWdodD0iMjQiIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHBhdGggZD0iTTEyIDJMMTMuMDkgOC4yNkwyMCA5TDEzLjA5IDE1Ljc0TDEyIDIyTDEwLjkxIDE1Ljc0TDQgOUwxMC45MSA4LjI2TDEyIDJaIiBmaWxsPSIjRkY2QzAwIi8+Cjwvc3ZnPgo=)](https://factorio.com/)
[![Concurrency](https://img.shields.io/badge/concurrency-goroutines-blue?style=for-the-badge)](https://golang.org/doc/effective_go.html#concurrency)
[![Architecture](https://img.shields.io/badge/architecture-pattern-9cf?style=for-the-badge)](https://en.wikipedia.org/wiki/Software_architecture_pattern)

This repository demonstrates a concurrency-oriented data flow pattern in Go inspired by the Main Bus concept from the game Factorio.

In Factorio, a Main Bus is a corridor of parallel belts dedicated to a single resource (iron plates, copper plates, steel, etc.). Each resource has its own bus, and each bus is made of multiple parallel belts (preferably an even number) to ensure symmetry and expansion.

We translate that concept to software:
- Each belt becomes a Go channel.
- Each Main Bus is a set of parallel channels dedicated to a single resource type.
- Multiple independent Main Buses can coexist and scale separately.

## Core Ideas

- Predictable, scalable flows per resource
- Separation of concerns: one bus per resource
- Easy horizontal scaling by adding more conveyors (channels)
- Independent consumers per conveyor to avoid contention across buses

## Structure

- `Conveyor` is a `chan Event` (a single belt)
- `MainBus` groups multiple conveyors for a single resource
- Producers push events onto a random conveyor within the bus
- Consumers read from conveyors independently (one goroutine per conveyor)

## Code Overview

The core types are in `main.go`:

- `Event` carries an ID, resource name, value (payload), and timestamp
- `Conveyor` is an alias for `chan Event`
- `MainBus` holds the resource name and slice of conveyors
- `NewMainBus(resource, lines, buffer)` creates a bus with an even number of conveyors
- `Produce(ev)` pushes an event to a random conveyor
- `Consume(line, wg)` reads from a single conveyor until closed
- `Close()` shuts down all conveyors in the bus

## Example

Run the example:

```bash
go run .
```

What it does:
- Creates two independent buses: `iron` (4 conveyors) and `copper` (2 conveyors)
- Starts one consumer goroutine per conveyor
- Produces 10 events for each bus, distributing them across conveyors
- Closes the buses and waits for all consumers to finish

Example output (truncated):

```
[Consumer-iron-L0] ID:0 Value:Iron Plate Time:12:34:56
[Consumer-copper-L1] ID:1 Value:Copper Plate Time:12:34:56
...
All main buses completed processing.
```

## Design Notes

- Even number of conveyors is enforced (if an odd number is passed, one is added) to mirror Factorio conventions and maintain balanced parallelism.
- Random distribution simulates simple load-balancing across conveyors. You can replace it with other strategies (round-robin, hashing by event ID, etc.).
- Each bus is independent. Adding a new resource means creating a new `MainBus` with its own set of conveyors.

## Extending

- Add metrics for per-conveyor throughput and backpressure
- Implement different balancing strategies (round-robin, consistent hashing)
- Add backpressure signals or context cancellation for graceful shutdowns
- Introduce typed payloads using generics (Go 1.18+) if needed

## ASCII Diagram

```text
                        🧩 MAIN BUS ARCHITECTURE (inspired by Factorio)

   ┌─────────────────────────────┐          ┌─────────────────────────────┐
   │         PRODUCERS           │          │         CONSUMERS           │
   └─────────────────────────────┘          └─────────────────────────────┘

         │                                   │
         │                                   │
         ▼                                   ▼

═══════════════════════════════════════════════════════════════════════════════
 MAIN BUS COMPLEX (each bus dedicated to a single resource)
═══════════════════════════════════════════════════════════════════════════════

    ┌──────────────────────────────────────────────────────────────────────┐
    │                              IRON BUS                               │
    │──────────────────────────────────────────────────────────────────────│
    │                                                                      │
    │   ┌────────────┐        ┌────────────┐        ┌────────────┐         │
    │   │ Producer 1 │───────▶│ Conveyor 0 │───────▶│ Consumer 0 │         │
    │   └────────────┘        └────────────┘        └────────────┘         │
    │                                                                      │
    │   ┌────────────┐        ┌────────────┐        ┌────────────┐         │
    │   │ Producer 2 │───────▶│ Conveyor 1 │───────▶│ Consumer 1 │         │
    │   └────────────┘        └────────────┘        └────────────┘         │
    │                                                                      │
    │   ┌────────────┐        ┌────────────┐        ┌────────────┐         │
    │   │ Producer 3 │───────▶│ Conveyor 2 │───────▶│ Consumer 2 │         │
    │   └────────────┘        └────────────┘        └────────────┘         │
    │                                                                      │
    │   ┌────────────┐        ┌────────────┐        ┌────────────┐         │
    │   │ Producer 4 │───────▶│ Conveyor 3 │───────▶│ Consumer 3 │         │
    │   └────────────┘        └────────────┘        └────────────┘         │
    └──────────────────────────────────────────────────────────────────────┘
                     ▲               ▲               ▲               ▲
                     │               │               │               │
                     │               │               │               │
                     └────── Parallel channels (chan Event) ────────┘


    ┌──────────────────────────────────────────────────────────────────────┐
    │                             COPPER BUS                              │
    │──────────────────────────────────────────────────────────────────────│
    │                                                                      │
    │   ┌────────────┐        ┌────────────┐        ┌────────────┐         │
    │   │ Producer 1 │───────▶│ Conveyor 0 │───────▶│ Consumer 0 │         │
    │   └────────────┘        └────────────┘        └────────────┘         │
    │                                                                      │
    │   ┌────────────┐        ┌────────────┐        ┌────────────┐         │
    │   │ Producer 2 │───────▶│ Conveyor 1 │───────▶│ Consumer 1 │         │
    │   └────────────┘        └────────────┘        └────────────┘         │
    └──────────────────────────────────────────────────────────────────────┘


═══════════════════════════════════════════════════════════════════════════════
 SYSTEM BEHAVIOR
═══════════════════════════════════════════════════════════════════════════════

 - Each MainBus handles one type of resource (e.g., iron, copper, steel)
 - Each bus contains multiple conveyors (channels) for parallel throughput
 - Producers send Events to a random conveyor in their bus
 - Consumers independently drain events from their assigned conveyor
 - Buses operate concurrently and independently
 - Adding a new resource = adding a new MainBus (no refactor needed)

═══════════════════════════════════════════════════════════════════════════════
 LEGEND
═══════════════════════════════════════════════════════════════════════════════
 ▶  Flow direction of data (Event)
 │  Connection between components
 ─  Conveyor (chan Event)
 ═  Structural separation / conceptual grouping
```


## UML/ASCII Overview

```text
                         🧩 MAIN BUS SYSTEM - UML ASCII

                 ┌────────────────────────────┐
                 │          Event             │
                 │───────────────────────────│
                 │ + ID: int                 │
                 │ + Resource: string        │
                 │ + Value: any              │
                 │ + Time: time.Time         │
                 └────────────────────────────┘


                 ┌────────────────────────────┐
                 │        Conveyor            │
                 │───────────────────────────│
                 │ type Conveyor chan Event  │
                 └────────────────────────────┘


                 ┌────────────────────────────┐
                 │        MainBus             │
                 │───────────────────────────│
                 │ - Resource: string        │
                 │ - Conveyors: []Conveyor   │
                 │───────────────────────────│
                 │ + Produce(ev Event)        │
                 │ + Consume(line int, wg*)  │
                 │ + Close()                  │
                 └────────────────────────────┘
                          ▲
                          │ 1..*
                          │ Contains
                          │
                 ┌────────────────────────────┐
                 │        Conveyor[0..N]      │
                 └────────────────────────────┘


                 ┌────────────────────────────┐
                 │        Producer             │
                 │───────────────────────────│
                 │ + Produce(ev Event, bus*) │
                 └────────────────────────────┘
                          │
                          │ sends Event
                          ▼
                 ┌────────────────────────────┐
                 │        MainBus             │
                 └────────────────────────────┘
                          │
                          │  distributes Event to random Conveyor
                          ▼
                 ┌────────────────────────────┐
                 │        Consumer            │
                 │───────────────────────────│
                 │ + Consume(Conveyor)       │
                 └────────────────────────────┘
```

### Conceptual Flow

1. Producer creates Events and sends them to the MainBus.
2. MainBus distributes events across its Conveyors (parallel channels).
3. Each Consumer drains from a specific Conveyor.
4. Each resource has its own independent MainBus.
5. The system is modular, scalable, and concurrent, enabling expansion without refactors.

## Requirements

- Go 1.20+

## License

MIT
