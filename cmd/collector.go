package main

import (
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/guptarohit/asciigraph"
	"github.com/jonathanlloyd/statsdbugger/pkg/statsd"
)

var colors = []asciigraph.AnsiColor{
	asciigraph.Red,
	asciigraph.Green,
	asciigraph.Blue,
	asciigraph.Cyan,
	asciigraph.Magenta,
	asciigraph.Yellow,
}

type Aggregator struct {
	Mu     sync.Mutex
	Gauges map[string]float64
	Tags   map[string]map[string]string
}

var a Aggregator = Aggregator{
	Gauges: make(map[string]float64),
	Tags:   make(map[string]map[string]string),
}

func listen(udpAddr *net.UDPAddr, out chan statsd.Metric) {
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}

		var metrics []statsd.Metric
		err = statsd.Unmarshal(buf[:n], &metrics)
		if err != nil {
			panic(err)
		}
		for _, metric := range metrics {
			if strings.HasPrefix(metric.Name(), "datadog.") {
				continue
			}
			out <- metric
		}
	}
}

var data = make(map[string][100]float64)

func plot() {
	for {
		a.Mu.Lock()

		for name, value := range a.Gauges {
			newData := [100]float64{}
			newData[0] = value
			for i := 0; i < 99; i++ {
				newData[i+1] = data[name][i]
			}
			data[name] = newData
		}

		keys := []string{}
		for k, _ := range data {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		plotData := make([][]float64, 0)
		for _, key := range keys {
			values := data[key]
			series := make([]float64, 100)
			copy(series, values[:])
			plotData = append(plotData, series)
		}

		plotColors := colors[:len(plotData)]
		plotLegend := make([]string, 0)
		for _, name := range keys {
			nameWithTags := fmt.Sprintf(
				"%s{%+v}",
				name,
				a.Tags[name],
			)

			plotLegend = append(plotLegend, nameWithTags)
		}

		graph := asciigraph.PlotMany(
			plotData,
			asciigraph.Precision(0),
			asciigraph.SeriesColors(plotColors...),
			asciigraph.SeriesLegends(plotLegend...),
			asciigraph.Caption("statsd bugger"),
			asciigraph.Width(100),
			asciigraph.LowerBound(0),
			asciigraph.UpperBound(40),
		)

		fmt.Print("\033[H\033[2J")
		fmt.Println(keys)
		fmt.Println(graph)

		a.Mu.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:8125")
	if err != nil {
		panic(err)
	}

	metrics := make(chan statsd.Metric)
	go listen(udpAddr, metrics)
	go plot()

	for metric := range metrics {
		switch m := metric.(type) {
		case statsd.Gauge:
			a.Mu.Lock()
			a.Gauges[m.GName] = m.GValue
			a.Tags[m.GName] = m.GTags
			a.Mu.Unlock()
		}
	}
}
