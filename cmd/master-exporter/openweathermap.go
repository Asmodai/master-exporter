/*
 * openweathermap.go --- OpenWeatherMap exporter.
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

package main

import (
	"github.com/Asmodai/gohacks/utils"

	"github.com/Asmodai/master-exporter/internal/config"
	"github.com/Asmodai/master-exporter/internal/exporter"
	"github.com/Asmodai/master-exporter/internal/openweathermap"
)

func (m *MasterExporter) initOMW() {
	cnf := m.config.AppConfig.(*config.AppConfig)

	if !utils.Member(cnf.Enabled, "openweathermap") {
		return
	}

	exp := openweathermap.NewExporter(
		m.apic,
		m.config.Logger,
		cnf.OpenWeatherMap,
	)

	// Scrape now
	if err := exp.Scrape(); err != nil {
		panic(err.Error())
	}
	m.config.Logger.Info(
		"Initial scrape complete.",
		"exporter", "openweathermap",
	)

	params := exporter.NewParams(
		"openweathermap",
		exp,
		m.config.ProcessManager,
		m.config.Logger,
	)

	_, err := exporter.Spawn(params)
	if err != nil {
		panic(err.Error())
	}
}

/* openweathermap.go ends here. */
