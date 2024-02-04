/*
 * icmp.go --- ICMP exporter.
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
	"github.com/Asmodai/master-exporter/internal/exporter"
	"github.com/Asmodai/master-exporter/internal/icmp"
)

func (m *MasterExporter) initIcmp() {
	cnf := m.config.AppConfig.(*config.AppConfig)

	if !utils.Member(cnf.Enabled, "icmp") {
		return
	}

	exp := icmp.NewExporter(
		m.appl.Context(),
		m.appl.Logger(),
		cnf.Icmp,
	)

	if err := exp.Setup(); err != nil {
		panic(err.Error())
	}
	if err := exp.Scrape(); err != nil {
		panic(err.Error())
	}
	m.config.Logger.Info(
		"Initial scrape complete.",
		"exporter", "icmp",
	)

	params := exporter.NewParams(
		"icmp",
		exp,
		m.config.ProcessManager,
		m.config.Logger,
	)

	_, err := exporter.Spawn(params)
	if err != nil {
		panic(err.Error())
	}
}

/* icmp.go ends here. */
