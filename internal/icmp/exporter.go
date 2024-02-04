/*
 * exporter.go --- ICMP exporter.
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

package icmp

import (
	"github.com/Asmodai/gohacks/logger"
	probing "github.com/prometheus-community/pro-bing"

	"context"
	"time"
)

var (
	icmpTimeout time.Duration = time.Second * 5
	icmpTtl     int           = 64
	icmpSize    int           = 24
	icmpCount   int           = 4
)

type Exporter struct {
	ctx     context.Context
	logger  logger.ILogger
	config  *Config
	metrics *Metrics
	calls   int
}

func NewExporter(ctx context.Context, logger logger.ILogger, config *Config) *Exporter {
	return &Exporter{
		ctx:     ctx,
		logger:  logger,
		config:  config,
		metrics: NewMetrics(),
		calls:   0,
	}
}

func (e *Exporter) ping(host string) (error, *probing.Statistics) {
	pinger, err := probing.NewPinger(host)
	if err != nil {
		e.logger.Fatal(
			"Could not create ICMP ping.",
			"err", err.Error(),
		)

		return err, nil
	}

	pinger.Count = icmpCount
	pinger.Size = icmpSize
	pinger.Timeout = icmpTimeout
	pinger.TTL = icmpTtl

	err = pinger.Run()
	if err != nil {
		e.logger.Warn(
			"Could not ping host.",
			"host", host,
			"err", err.Error(),
		)

		return err, nil
	}

	return nil, pinger.Statistics()
}

func (e *Exporter) Interval() int {
	return e.config.Interval
}

func (e *Exporter) Setup() error {
	for _, h := range e.config.Hosts {
		if c := e.metrics.HasHost(h); !c {
			m := e.metrics.GetHost(h)

			m.AddMetric("packet_loss", "Packet loss.", h)
			m.AddMetric("min_rtt", "Minimum RTT value. Nanoseconds.", h)
			m.AddMetric("avg_rtt", "Average RTT value. Nanoseconds.", h)
			m.AddMetric("max_rtt", "Maximum RTT value. Nanoseconds.", h)
			m.AddMetric("stddev_rtt", "Standard deviation of RTT value. Nanoseconds.", h)
		}
	}

	return nil
}

func (e *Exporter) Scrape() error {
	for _, h := range e.config.Hosts {
		err, res := e.ping(h)
		if err != nil {
			return err
		}
		if res != nil {
			m := e.metrics.GetHost(h)

			m.SetMetric("packet_loss", res.PacketLoss)
			m.SetMetric("min_rtt", float64(res.MinRtt))
			m.SetMetric("avg_rtt", float64(res.AvgRtt))
			m.SetMetric("max_rtt", float64(res.MaxRtt))
			m.SetMetric("stddev_rtt", float64(res.StdDevRtt))
		}
	}

	return nil
}

/* exporter.go ends here. */
