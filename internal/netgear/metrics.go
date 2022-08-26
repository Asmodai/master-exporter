/*
 * metrics.go --- Netgear metrics.
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person
 * obtaining a copy of this software and associated documentation files
 * (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge,
 * publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
 * BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package netgear

import (
	"github.com/prometheus/client_golang/prometheus"

	"fmt"
)

type SwitchMetrics struct {
	Metric map[string]prometheus.Gauge
}

func NewSwitchMetrics() *SwitchMetrics {
	return &SwitchMetrics{
		Metric: map[string]prometheus.Gauge{},
	}
}

func (sm *SwitchMetrics) AddMetric(name, help, pretty string) {
	if sm.Metric == nil {
		sm.Metric = map[string]prometheus.Gauge{}
	}

	if _, ok := sm.Metric[name]; !ok {
		fmt.Printf("Adding metric %s{switch=\"%s\"}\n", name, pretty)
		sm.Metric[name] = prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "netgear",
			Name:      name,
			Help:      help,
			ConstLabels: map[string]string{
				"switch": pretty,
			},
		})
		_ = prometheus.Register(sm.Metric[name])
	}
}

func (sm *SwitchMetrics) AddPortMetric(name, help, pretty string, port int) {
	if sm.Metric == nil {
		sm.Metric = map[string]prometheus.Gauge{}
	}

	sport := fmt.Sprintf("%02d", port+1)

	if _, ok := sm.Metric[name+sport]; !ok {
		sm.Metric[name+sport] = prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "netgear",
			Name:      name,
			Help:      help,
			ConstLabels: map[string]string{
				"port":   sport,
				"switch": pretty,
			},
		})
		_ = prometheus.Register(sm.Metric[name+sport])
	}
}

func (sm *SwitchMetrics) SetMetric(name string, value uint64) {
	if _, ok := sm.Metric[name]; !ok {
		return
	}

	sm.Metric[name].Set(float64(value))
}

func (sm *SwitchMetrics) SetPortMetric(name string, port int, value uint64) {
	sport := fmt.Sprintf("%02d", port+1)

	if _, ok := sm.Metric[name+sport]; !ok {
		return
	}

	sm.Metric[name+sport].Set(float64(value))
}

// =================================================================

type Metrics struct {
	metrics map[string]*SwitchMetrics
}

func NewMetrics() *Metrics {
	return &Metrics{
		metrics: map[string]*SwitchMetrics{},
	}
}

func (m *Metrics) Keys() []string {
	keys := []string{}

	for k := range m.metrics {
		keys = append(keys, k)
	}

	return keys
}

func (m *Metrics) HasSwitch(key string) bool {
	_, ok := m.metrics[key]

	return ok
}

func (m *Metrics) AddSwitch(key string) {
	if _, ok := m.metrics[key]; ok {
		return
	}

	m.metrics[key] = NewSwitchMetrics()
}

func (m *Metrics) GetSwitch(key string) *SwitchMetrics {
	m.AddSwitch(key)

	return m.metrics[key]
}

/* metrics.go ends here. */
