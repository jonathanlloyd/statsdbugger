package statsd

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func Unmarshal(data []byte, metrics *[]Metric) error {
	if !utf8.Valid(data) {
		return fmt.Errorf("data must be valid UTF-8")
	}
	dataUtf8 := string(data)

	parsedMetrics := []Metric{}
	lines := strings.Split(dataUtf8, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		metric, err := parseMetric(line)
		if err != nil {
			return err
		}

		parsedMetrics = append(parsedMetrics, metric)
	}

	*metrics = parsedMetrics
	return nil
}

// parseMetric parses a single metric line and returns a Metric object.
// The line must be in the format:
// <name>:<value>|<type>[|@<sample_rate>][|#<tag1>,<tag2>]
func parseMetric(line string) (Metric, error) {
	parts := strings.Split(line, "|")
	if len(parts) < 2 {
		return nil, fmt.Errorf(
			"malformed line, must contain <name>:<value>|<type>",
		)
	}

	nameParts := strings.Split(parts[0], ":")
	if len(nameParts) != 2 {
		return nil, fmt.Errorf("malformed line, must contain <name>:<value>")
	}

	name := nameParts[0]

	valueStr := nameParts[1]
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return nil, fmt.Errorf(
			"malformed line, value must be int|float",
		)
	}

	metricType := parts[1]
	_ = metricType

	var sampleRate float64 = 1
	var sampleRatePresent bool = false
	if len(parts) > 2 && strings.HasPrefix(parts[2], "@") {
		sampleRatePresent = true
		sampleRateStr := parts[2][1:]
		sampleRate, err = strconv.ParseFloat(sampleRateStr, 64)
		if err != nil {
			return nil, fmt.Errorf(
				"malformed line, sample rate must be int|float",
			)
		}
	}

	var tags map[string]string = map[string]string{}
	var tagIndex int
	if sampleRatePresent {
		tagIndex = 3
	} else {
		tagIndex = 2
	}
	if len(parts) > tagIndex && strings.HasPrefix(parts[tagIndex], "#") {
		tagsStr := parts[tagIndex][1:]
		kvPairs := strings.Split(tagsStr, ",")
		for _, kvPair := range kvPairs {
			kv := strings.Split(kvPair, ":")
			if len(kv) != 2 {
				return nil, fmt.Errorf(
					"malformed line, tags must be in the format key:value",
				)
			}
			tags[kv[0]] = kv[1]
		}
	}

	switch metricType {
	case "c":
		return Counter{
			CName:       name,
			CValue:      value,
			CSampleRate: sampleRate,
			CTags:       tags,
		}, nil
	case "g":
		return Gauge{
			GName:  name,
			GValue: value,
			GTags:  tags,
		}, nil
	default:
		return nil, fmt.Errorf("unknown metric type: %s", metricType)
	}
}
