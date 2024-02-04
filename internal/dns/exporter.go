/*
 * exporter.go --- DNS exporter.
 *
 * Copyright (c) 2020-2024 Paul Ward <asmodai@gmail.com>
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
	"github.com/Asmodai/gohacks/logger"

	"context"
	"net"
	"time"
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

func (e *Exporter) lookup(host string) (error, time.Duration) {
	var r *net.Resolver

	ctx, cancel := context.WithTimeout(context.TODO(), time.Millisecond*5)
	defer cancel()

	start := time.Now()
	if _, err := r.LookupHost(ctx, host); err != nil {
		e.logger.Warn(
			"Could not resolve host.",
			"host", host,
			"err", err.Error(),
		)

		return err, time.Hour * 24
	}

	return nil, time.Since(start)
}

func (e *Exporter) Interval() int {
	return e.config.Interval
}

func (e *Exporter) Setup() error {
	for _, h := range e.config.Hosts {
		if c := e.metrics.HasHost(h); !c {
			m := e.metrics.GetHost(h)

			m.AddMetric("response_time", "DNS query response time. Nanoseconds.", h)
		}
	}

	return nil
}

func (e *Exporter) Scrape() error {
	for _, h := range e.config.Hosts {
		err, res := e.lookup(h)
		if err != nil {
			return err
		}

		m := e.metrics.GetHost(h)
		m.SetMetric("response_time", float64(res))
	}

	return nil
}

/* exporter.go ends here. */
