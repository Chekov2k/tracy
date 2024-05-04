package main

import (
	"time"
)

func work(fiber *Fiber, frame *Frame) {
	frame.Start()
	defer frame.End()
	zone := fiber.NewZone("work", ColorDeepSkyBlue, 2)
	defer zone.Close()
	zone.Text("bam")
	time.Sleep(time.Millisecond * 10)
	inner := fiber.NewZone("inner", ColorDarkRed, 2)
	defer inner.Close()
	inner.Text("yay")
}

func main() {
	fiber := NewFiber("work")
	defer fiber.Close()

	frame := NewFrame("frame", false)
	for {
		work(fiber, frame)
	}
}
