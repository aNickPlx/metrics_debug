package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	requests    = "metrics_test_requests"
	requests10  = "metrics_test_requests_ten"
	requests100 = "metrics_test_requests_hundred"
	requestsSin = "metrics_test_requests_sin"

	period    = 100
	amplitude = 100
)

func main() {
	metrics := metrics{}

	http.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		metrics.updateRequestMetrics()
		io.WriteString(w, metrics.getTemplate())
	})

	http.ListenAndServe(":8080", nil)
}

type metrics struct {
	// this mutex protects all the fields in this struct
	mu          sync.RWMutex
	lastQueried int64
	requests    int
	requests10  int
	requests100 int
	sin         int
}

func (m *metrics) updateRequestMetrics() {
	m.mu.Lock()
	defer m.mu.Unlock()

	timestamp := time.Now().Unix()
	if m.lastQueried == 0 {
		m.lastQueried = timestamp
	}

	d := int(timestamp - m.lastQueried)
	m.requests += d
	m.requests10 += (d * 10)
	m.requests100 += (d * 100)
	m.sin += int(getSinTimestamp(timestamp, d))

	m.lastQueried = timestamp
}

func getSinTimestamp(timestamp int64, delta int) float64 {
	s := math.Sin(float64(timestamp) / float64(period))
	requests := float64(delta * amplitude / 2)
	return math.Round((1 + s) * requests)
}

const tpl = `# TYPE %s counter
%s %d
`

func (m *metrics) getTemplate() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var b strings.Builder

	b.WriteString(fmt.Sprintf(tpl, requests, requests, m.requests))
	b.WriteString(fmt.Sprintf(tpl, requests10, requests10, m.requests10))
	b.WriteString(fmt.Sprintf(tpl, requests100, requests100, m.requests100))
	b.WriteString(fmt.Sprintf(tpl, requestsSin, requestsSin, m.sin))

	return b.String()
}
