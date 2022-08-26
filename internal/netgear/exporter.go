/*
 * exporter.go --- netgear exporter.
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
	//"github.com/Asmodai/master-exporter/internal/exporter"
	"github.com/Asmodai/gohacks/logger"

	"github.com/yaamai/go-nsdp/nsdp"

	"context"
)

var (
	tlvs []nsdp.TLV = []nsdp.TLV{nsdp.HostName{}, nsdp.HostIPAddress{}, nsdp.PortLinkStatus{}, nsdp.PortStatistics{}}
)

type NsdpValues map[string]interface{}

type Exporter struct {
	ctx     context.Context
	logger  logger.ILogger
	config  *Config
	client  *nsdp.Client
	metrics *Metrics
	calls   int
}

func NewExporter(ctx context.Context, logger logger.ILogger, config *Config) *Exporter {
	nsdpClient, err := nsdp.NewDefaultClient()
	if err != nil {
		logger.Fatal(
			"Could not create NSDP client.",
			"err", err.Error(),
		)
	}

	return &Exporter{
		ctx:     ctx,
		logger:  logger,
		config:  config,
		client:  nsdpClient,
		metrics: NewMetrics(),
		calls:   0,
	}
}

func (e *Exporter) check(vals NsdpValues) {
	for _, k := range e.metrics.Keys() {
		e.metrics.GetSwitch(k).SetMetric("up", 0)

		for key, val := range vals {
			if key == "host_name" {
				if val.(*nsdp.HostName).String() == k {
					e.metrics.GetSwitch(k).SetMetric("up", 1)
					continue
				}
			}
		}
	}
}

func (e *Exporter) process(vals NsdpValues) {
	var hostname string
	var addr *nsdp.HostIPAddress
	var portstatus []nsdp.TLV = []nsdp.TLV{}
	var portstats []nsdp.TLV = []nsdp.TLV{}
	var created bool

	for key, val := range vals {
		switch key {
		case "host_name":
			hostname = val.(*nsdp.HostName).String()
			created = e.metrics.HasSwitch(hostname)

		case "ip":
			addr = val.(*nsdp.HostIPAddress)

		case "port_link_status":
			portstatus = val.([]nsdp.TLV)

		case "port_statistics":
			portstats = val.([]nsdp.TLV)
		}
	}

	if len(hostname) == 0 {
		e.logger.Warn(
			"Switch is not returning a hostname.",
			"addr", addr.String(),
		)

		return
	}

	if !created {
		e.metrics.AddSwitch(hostname)

		e.metrics.GetSwitch(hostname).AddMetric(
			"up",
			"Is the given switch online?",
			hostname,
		)
	}

	for idx := range portstatus {
		stat := portstatus[idx].(*nsdp.PortLinkStatus)

		if !created {
			e.metrics.GetSwitch(hostname).AddPortMetric(
				"speed",
				"Port speed.",
				hostname,
				idx,
			)
		}

		e.metrics.GetSwitch(hostname).SetPortMetric("speed", idx, uint64(stat.Speed))
	}

	for idx := range portstats {
		stat := portstats[idx].(*nsdp.PortStatistics)

		if !created {
			e.metrics.GetSwitch(hostname).AddPortMetric(
				"rx_total_bytes",
				"Total bytes received for a port.",
				hostname,
				idx,
			)
			e.metrics.GetSwitch(hostname).AddPortMetric(
				"tx_total_bytes",
				"Total bytes transmitted for a port.",
				hostname,
				idx,
			)
			e.metrics.GetSwitch(hostname).AddPortMetric(
				"packets",
				"Total packets on this port.",
				hostname,
				idx,
			)
			e.metrics.GetSwitch(hostname).AddPortMetric(
				"packets_bcast",
				"Total broadcast packets on this port.",
				hostname,
				idx,
			)
			e.metrics.GetSwitch(hostname).AddPortMetric(
				"packets_mcast",
				"Total multicast packets on this port.",
				hostname,
				idx,
			)
			e.metrics.GetSwitch(hostname).AddPortMetric(
				"crc_errors",
				"Total CRC errors on this port.",
				hostname,
				idx,
			)
		}

		e.metrics.GetSwitch(hostname).SetPortMetric("rx_total_bytes", idx, uint64(stat.Recv))
		e.metrics.GetSwitch(hostname).SetPortMetric("tx_total_bytes", idx, uint64(stat.Send))
		e.metrics.GetSwitch(hostname).SetPortMetric("packets", idx, uint64(stat.Pkt))
		e.metrics.GetSwitch(hostname).SetPortMetric("packets_bcast", idx, uint64(stat.Broadcast))
		e.metrics.GetSwitch(hostname).SetPortMetric("packets_mcast", idx, uint64(stat.Multicast))
		e.metrics.GetSwitch(hostname).SetPortMetric("crx_errors", idx, uint64(stat.Error))
	}
}

func (e *Exporter) poll() error {
	resp, err := e.client.Read(tlvs...)
	if err != nil {
		return err
	}

	tlvmap := NsdpValues{}
	for _, tlv := range resp.Body {
		tname := tlv.Tag().String()
		if _, ok := tlvmap[tname]; ok {
			array, ok := tlvmap[tname].([]nsdp.TLV)
			if !ok {
				prev := tlvmap[tname]
				tlvmap[tname] = []nsdp.TLV{prev.(nsdp.TLV), tlv}
			} else {
				tlvmap[tname] = append(array, tlv)
			}
		} else {
			tlvmap[tname] = tlv
		}
	}

	e.check(tlvmap)
	e.process(tlvmap)

	return nil
}

func (e *Exporter) Interval() int {
	return e.config.Interval
}

func (e *Exporter) Scrape() error {
	return e.poll()
}

/* exporter.go ends here. */
