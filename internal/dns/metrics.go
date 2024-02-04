/*
 * metrics.go --- DNS metrics.
 *
 * Copyright (c) 2024 Paul Ward <asmodai@gmail.com>
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

package dns

import (
	"github.com/prometheus/client_golang/prometheus"
)

type DnsMetrics struct {
	Metric map[string]prometheus.Gauge
}

func NewDnsMetrics() *DnsMetrics {
	return &DnsMetrics{
		Metric: map[string]prometheus.Gauge{},
	}
}

func (dm *DnsMetrics) AddMetric(name, help, pretty string) {
	if dm.Metric == nil {
		dm.Metric = map[string]prometheus.Gauge{}
	}

	if _, ok := dm.Metric[name]; !ok {
		dm.Metric[name] = prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "dns",
			Name:      name,
			Help:      help,
			ConstLabels: map[string]string{
				"host": pretty,
			},
		})
		_ = prometheus.Register(dm.Metric[name])
	}
}

func (dm *DnsMetrics) SetMetric(name string, value float64) {
	if _, ok := dm.Metric[name]; !ok {
		return
	}

	dm.Metric[name].Set(value)
}

// =================================================================

type Metrics struct {
	metrics map[string]*DnsMetrics
}

func NewMetrics() *Metrics {
	return &Metrics{
		metrics: map[string]*DnsMetrics{},
	}
}

func (m *Metrics) Keys() []string {
	keys := []string{}

	for k := range m.metrics {
		keys = append(keys, k)
	}

	return keys
}

func (m *Metrics) HasHost(key string) bool {
	_, ok := m.metrics[key]

	return ok
}

func (m *Metrics) AddHost(key string) {
	if _, ok := m.metrics[key]; ok {
		return
	}

	m.metrics[key] = NewDnsMetrics()
}

func (m *Metrics) GetHost(key string) *DnsMetrics {
	m.AddHost(key)

	return m.metrics[key]
}

/* metrics.go ends here. */
