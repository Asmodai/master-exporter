/*
 * metrics.go --- Sabnzbd metrics.
 *
 * Copyright (c) 2022-2024 Paul Ward <asmodai@gmail.com>
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

package sabnzbd

import (
	"github.com/Asmodai/master-exporter/internal/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func NewGauge(name string) prometheus.Gauge {
	return promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "sabnzbd",
		Name:      name,
		Help:      "SabNZBd data.",
	})
}

func NewServerGauge(name, server string) prometheus.Gauge {
	return promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "sabnzbd",
		Name:      name,
		Help:      "SabNZBd data.",
		ConstLabels: map[string]string{
			"server": server,
		},
	})
}

type Metrics struct {
	SpeedLimit    prometheus.Gauge
	SpeedLimitAbs prometheus.Gauge
	Speed         prometheus.Gauge
	Kbs           prometheus.Gauge
	MbTotal       prometheus.Gauge
	MbLeft        prometheus.Gauge
	MbDone        prometheus.Gauge
	SizeTotal     prometheus.Gauge
	SizeLeft      prometheus.Gauge
	TimeLeft      prometheus.Gauge
	SlotCount     prometheus.Gauge
	SlotTotal     prometheus.Gauge

	XferTotal  prometheus.Gauge
	ServerXfer map[string]prometheus.Gauge

	Calls prometheus.Gauge
	Limit prometheus.Gauge
}

func NewMetrics() *Metrics {
	return &Metrics{
		SpeedLimit:    NewGauge("download_speed_limit"),
		SpeedLimitAbs: NewGauge("download_speed_limit_abs"),
		Speed:         NewGauge("download_speed"),
		Kbs:           NewGauge("download_kb_per_sec"),
		MbTotal:       NewGauge("queue_mb_total"),
		MbLeft:        NewGauge("queue_mb_left"),
		MbDone:        NewGauge("queue_mb_done"),
		SizeTotal:     NewGauge("queue_size_total"),
		SizeLeft:      NewGauge("queue_size_left"),
		TimeLeft:      NewGauge("download_time_left"),
		SlotCount:     NewGauge("job_slots_count"),
		SlotTotal:     NewGauge("job_slots_total"),
		XferTotal:     NewGauge("server_xfer_total"),
		ServerXfer:    map[string]prometheus.Gauge{},
		Calls:         metrics.NewMetricsGauge("calls", "sabnzbd"),
		Limit:         metrics.NewMetricsGauge("limit", "sabnzbd"),
	}
}

func (m *Metrics) SetSpeedLimit(v float64)    { m.SpeedLimit.Set(v) }
func (m *Metrics) SetSpeedLimitAbs(v float64) { m.SpeedLimitAbs.Set(v) }
func (m *Metrics) SetSpeed(v float64)         { m.Speed.Set(v) }
func (m *Metrics) SetKbs(v float64)           { m.Kbs.Set(v) }
func (m *Metrics) SetMbTotal(v float64)       { m.MbTotal.Set(v) }
func (m *Metrics) SetMbLeft(v float64)        { m.MbLeft.Set(v) }
func (m *Metrics) SetMbDone(v float64)        { m.MbDone.Set(v) }
func (m *Metrics) SetSizeTotal(v float64)     { m.SizeTotal.Set(v) }
func (m *Metrics) SetSizeLeft(v float64)      { m.SizeLeft.Set(v) }
func (m *Metrics) SetTimeLeft(v float64)      { m.TimeLeft.Set(v) }
func (m *Metrics) SetSlotCount(v float64)     { m.SlotCount.Set(v) }
func (m *Metrics) SetSlotTotal(v float64)     { m.SlotTotal.Set(v) }
func (m *Metrics) SetXferTotal(v float64)     { m.XferTotal.Set(v) }

func (m *Metrics) SetServerXfer(name string, v float64) {
	if _, ok := m.ServerXfer[name]; !ok {
		m.ServerXfer[name] = NewServerGauge("xfer", name)
	}

	m.ServerXfer[name].Set(v)
}

func (m *Metrics) SetCalls(val float64) { m.Calls.Set(val) }
func (m *Metrics) SetLimit(val float64) { m.Limit.Set(val) }

/* metrics.go ends here. */
