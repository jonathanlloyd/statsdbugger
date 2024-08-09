package statsd

import (
	"reflect"
	"testing"
)

func TestUnmarshalCounter(t *testing.T) {
	data := []byte("pugs.cuddled:1|c")
	metrics := []Metric{}
	err := Unmarshal(data, &metrics)
	if err != nil {
		t.Fatalf("Unmarshal failed: %s", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("Expected 1 metric, got %d", len(metrics))
	}

	metric := metrics[0]
	counter, ok := metric.(Counter)
	if !ok {
		t.Fatalf("Expected Counter, got %T", metric)
	}

	if counter.CName != "pugs.cuddled" {
		t.Fatalf("Expected pugs.cuddled, got %s", counter.CName)
	}
	if counter.CValue != 1 {
		t.Fatalf("Expected 1, got %f", counter.CValue)
	}
}

func TestUnmarshalSampledCounter(t *testing.T) {
	data := []byte("pugs.cuddled:1|c|@0.5")
	metrics := []Metric{}
	err := Unmarshal(data, &metrics)
	if err != nil {
		t.Fatalf("Unmarshal failed: %s", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("Expected 1 metric, got %d", len(metrics))
	}

	metric := metrics[0]
	counter, ok := metric.(Counter)
	if !ok {
		t.Fatalf("Expected Counter, got %T", metric)
	}

	if counter.CSampleRate != 0.5 {
		t.Fatalf("Expected 0.5, got %f", counter.CSampleRate)
	}
}

func TestUnmarshalCounterWithTags(t *testing.T) {
	data := []byte("pugs.cuddled:1|c|#env:prod,service:puginator")
	metrics := []Metric{}
	err := Unmarshal(data, &metrics)
	if err != nil {
		t.Fatalf("Unmarshal failed: %s", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("Expected 1 metric, got %d", len(metrics))
	}

	metric := metrics[0]

	expectedTags := map[string]string{
		"env":     "prod",
		"service": "puginator",
	}
	if !reflect.DeepEqual(metric.Tags(), expectedTags) {
		t.Fatalf("Expected %v, got %v", expectedTags, metric.Tags())
	}
}

func TestUnmarshalSampledCounterWithTags(t *testing.T) {
	data := []byte("pugs.cuddled:1|c|@0.5|#env:prod,service:puginator")
	metrics := []Metric{}
	err := Unmarshal(data, &metrics)
	if err != nil {
		t.Fatalf("Unmarshal failed: %s", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("Expected 1 metric, got %d", len(metrics))
	}

	metric := metrics[0]
	counter, ok := metric.(Counter)
	if !ok {
		t.Fatalf("Expected Counter, got %T", metric)
	}

	if counter.CSampleRate != 0.5 {
		t.Fatalf("Expected 0.5, got %f", counter.CSampleRate)
	}

	expectedTags := map[string]string{
		"env":     "prod",
		"service": "puginator",
	}
	if !reflect.DeepEqual(counter.CTags, expectedTags) {
		t.Fatalf("Expected %v, got %v", expectedTags, metric.Tags())
	}
}

func TestUnmarshalMultiple(t *testing.T) {
	data := []byte("pugs.cuddled:1|c\npugs.cuddled:1|c")
	metrics := []Metric{}
	err := Unmarshal(data, &metrics)
	if err != nil {
		t.Fatalf("Unmarshal failed: %s", err)
	}

	if len(metrics) != 2 {
		t.Fatalf("Expected 2 metrics, got %d", len(metrics))
	}
}

func TestUnmarshalGauge(t *testing.T) {
	data := []byte("pugs.cuddled:1|g")
	metrics := []Metric{}
	err := Unmarshal(data, &metrics)
	if err != nil {
		t.Fatalf("Unmarshal failed: %s", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("Expected 1 metric, got %d", len(metrics))
	}

	metric := metrics[0]
	gauge, ok := metric.(Gauge)
	if !ok {
		t.Fatalf("Expected Gauge, got %T", metric)
	}

	if gauge.GName != "pugs.cuddled" {
		t.Fatalf("Expected pugs.cuddled, got %s", gauge.GName)
	}
	if gauge.GValue != 1 {
		t.Fatalf("Expected 1, got %f", gauge.GValue)
	}
}
