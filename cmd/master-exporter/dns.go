/*
 * dns.go --- DNS resolver exporter.
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

package main

import (
	"github.com/Asmodai/gohacks/utils"

	"github.com/Asmodai/master-exporter/internal/config"
	"github.com/Asmodai/master-exporter/internal/dns"
	"github.com/Asmodai/master-exporter/internal/exporter"
)

func (m *MasterExporter) initDns() {
	cnf := m.config.AppConfig.(*config.AppConfig)

	if !utils.Member(cnf.Enabled, "dns") {
		return
	}

	exp := dns.NewExporter(
		m.appl.Context(),
		m.appl.Logger(),
		cnf.Dns,
	)

	if err := exp.Setup(); err != nil {
		panic(err.Error())
	}
	if err := exp.Scrape(); err != nil {
		panic(err.Error())
	}
	m.config.Logger.Info(
		"Initial scrape complete.",
		"exporter", "dns",
	)

	params := exporter.NewParams(
		"dns",
		exp,
		m.config.ProcessManager,
		m.config.Logger,
	)

	_, err := exporter.Spawn(params)
	if err != nil {
		panic(err.Error())
	}
}

/* dns.go ends here. */
