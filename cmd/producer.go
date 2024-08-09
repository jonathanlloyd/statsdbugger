package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

func newValue(oldValue int) int {
	n := (oldValue + rand.Intn(10) - 5)
	if n < 0 {
		n = 0
	}
	if n > 40 {
		n = 40
	}
	return n
}

func main() {
	statsd, err := statsd.New("127.0.0.1:8125")
	if err != nil {
		log.Fatal(err)
	}
	var pugCuddleCount int = 1
	var barkCount int = 2
	var snoreCount int = 3

	for true {
		pugCuddleCount = newValue(pugCuddleCount)
		barkCount = newValue(barkCount)
		snoreCount = newValue(snoreCount)

		fmt.Println("Sending metrics...")
		statsd.Gauge(
			"pug.cuddles",
			float64(pugCuddleCount),
			[]string{"environment:dev"},
			1,
		)
		statsd.Gauge(
			"pug.barks",
			float64(barkCount),
			[]string{"environment:dev"},
			1,
		)
		statsd.Gauge(
			"pug.snores",
			float64(snoreCount),
			[]string{"environment:dev"},
			1,
		)
		time.Sleep(1 * time.Second)
	}
}
