/*
 * main.go --- Master exporter command.
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
	"github.com/Asmodai/gohacks/apiclient"
	"github.com/Asmodai/gohacks/app"
	"github.com/Asmodai/gohacks/logger"
	"github.com/Asmodai/gohacks/process"
	"github.com/Asmodai/gohacks/semver"

	"github.com/Asmodai/master-exporter/internal/config"
)

var (
	GitVers string = "0:<local>"
)

type MasterExporter struct {
	config *app.Config
	appl   *app.Application
	apic   apiclient.IApiClient
}

func NewMasterExporter() *MasterExporter {
	version, err := semver.MakeSemVer(GitVers)
	if err != nil {
		panic(err.Error())
	}

	c := &app.Config{
		Name:           "Master-Exporter",
		Version:        version,
		Logger:         logger.NewLogger(),
		ProcessManager: process.NewManager(),
		AppConfig:      &config.AppConfig{},
	}

	a := app.NewApplication(c)
	c.ProcessManager.SetContext(a.Context())
	a.Init()

	if err := c.AppConfig.(*config.AppConfig).Init(); err != nil {
		panic(err.Error())
	}

	apic := apiclient.NewClient(
		c.AppConfig.(*config.AppConfig).GetAPIClient(),
		c.Logger,
	)

	return &MasterExporter{
		config: c,
		appl:   a,
		apic:   apic,
	}
}

func (m *MasterExporter) Main() {
	m.initOMW()
	m.initSab()
	m.initNG()
	m.initIcmp()
	m.initDns()
	m.initPrometheus()

	m.appl.Run()
}

func main() {
	me := NewMasterExporter()

	me.Main()
}

/* main.go ends here. */
