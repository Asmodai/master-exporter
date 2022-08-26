/*
 * exporter.go --- SabNZBd exporter.
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

package sabnzbd

import (
	//"github.com/Asmodai/master-exporter/internal/exporter"

	"github.com/Asmodai/gohacks/apiclient"
	"github.com/Asmodai/gohacks/logger"

	"encoding/json"
	"fmt"
	"time"
)

type Exporter struct {
	client  apiclient.IApiClient
	logger  logger.ILogger
	config  *Config
	data    *SabNZBd
	metrics *Metrics
	calls   int
}

func NewExporter(client apiclient.IApiClient, logger logger.ILogger, config *Config) *Exporter {
	return &Exporter{
		client:  client,
		logger:  logger,
		config:  config,
		data:    NewSabNZBd(),
		metrics: NewMetrics(),
		calls:   0,
	}
}

func (e *Exporter) Data() *SabNZBd {
	return e.data
}

func (e *Exporter) canCall() bool { return true }

func (e *Exporter) resetLimit() {
	now := time.Now()

	if now.Hour() == 0 && now.Minute() == 0 {
		if e.calls > 1 {
			e.calls = 0
		}
	}
}

func (e *Exporter) get(mode string) ([]byte, error) {
	e.resetLimit()

	if !e.canCall() {
		return []byte{}, fmt.Errorf("Daily call limit exceeded.")
	}

	params := &apiclient.Params{
		Url: fmt.Sprintf("%s&mode=%s", e.config.GetURL(), mode),
	}

	data, code, err := e.client.Get(params)
	if err != nil {
		return []byte{}, err
	}
	e.calls++

	e.metrics.SetCalls(float64(e.calls))
	e.metrics.SetLimit(float64(35000.0))

	switch code {
	case 200:
		return data, nil
	default:
		return []byte{}, fmt.Errorf("%d response: %s", code, err)
	}
}

func (e *Exporter) Interval() int {
	return e.config.Interval
}

func (e *Exporter) Scrape() error {
	queue, err := e.get(modeQueue)
	if err != nil {
		e.logger.Warn(
			"Scrape error.",
			"err", err.Error(),
			"exporter", "sabnzbd",
		)
	}

	server, err := e.get(modeServer)
	if err != nil {
		e.logger.Warn(
			"Scrape error.",
			"err", err.Error(),
			"exporter", "sabnzbd",
		)
	}

	err = json.Unmarshal(queue, e.data.Queue)
	if err != nil {
		return fmt.Errorf("JSON unmarshal: %s", err)
	}

	err = json.Unmarshal(server, e.data.Server)
	if err != nil {
		return fmt.Errorf("JSON unmarshal: %s", err)
	}

	e.metrics.SetSpeedLimit(e.data.Queue.Queue.SpeedLimit.Float64())
	e.metrics.SetSpeedLimitAbs(e.data.Queue.Queue.SpeedLimitAbs.Float64())
	e.metrics.SetKbs(e.data.Queue.Queue.KbPerSec.Float64())
	e.metrics.SetMbTotal(e.data.Queue.Queue.Mb.Float64() * 1.049e+6)
	e.metrics.SetMbLeft(e.data.Queue.Queue.MbLeft.Float64() * 1.049e+6)
	e.metrics.SetMbDone(e.data.Queue.Queue.MbDone() * 1.049e+6)
	e.metrics.SetSizeTotal(e.data.Queue.Queue.Size.Float64())
	e.metrics.SetSizeLeft(e.data.Queue.Queue.SizeLeft.Float64())
	e.metrics.SetTimeLeft(e.data.Queue.Queue.TimeLeft.Seconds())
	e.metrics.SetSlotCount(float64(e.data.Queue.Queue.NoOfSlots))
	e.metrics.SetSlotTotal(float64(e.data.Queue.Queue.NoOfSlotsTotal))

	e.metrics.SetXferTotal(float64(e.data.Server.Total))
	for k := range e.data.Server.Servers {
		e.metrics.SetServerXfer(k, float64(e.data.Server.Servers[k].Total))
	}

	return nil
}

/* exporter.go ends here. */
