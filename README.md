# Main Bus Architecture in Go (Inspired by Factorio)

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

## Requirements

- Go 1.20+

## License

MIT
