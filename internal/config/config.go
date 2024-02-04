/*
 * config.go --- Exporter configuration.
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

package config

import (
	"github.com/Asmodai/gohacks/apiclient"

	"github.com/Asmodai/master-exporter/internal/dns"
	"github.com/Asmodai/master-exporter/internal/icmp"
	"github.com/Asmodai/master-exporter/internal/netgear"
	"github.com/Asmodai/master-exporter/internal/openweathermap"
	"github.com/Asmodai/master-exporter/internal/sabnzbd"
)

type AppConfig struct {
	BasePort int      `json:"base_port"`
	Enabled  []string `json:"enabled"`

	ApiClient *apiclient.Config `json:"api_client"`

	OpenWeatherMap *openweathermap.Config `json:"openweathermap"`
	SabNZBd        *sabnzbd.Config        `json:"sabnzbd"`
	Netgear        *netgear.Config        `json:"netgear"`
	Icmp           *icmp.Config           `json:"icmp"`
	Dns            *dns.Config            `json:"dns"`
}

func (c *AppConfig) Init() error {
	if c.BasePort == 0 {
		c.BasePort = 9500
	}

	if c.ApiClient == nil {
		c.ApiClient = apiclient.NewDefaultConfig()
	}

	openweathermap.Validate(c.OpenWeatherMap)
	sabnzbd.Validate(c.SabNZBd)
	netgear.Validate(c.Netgear)
	icmp.Validate(c.Icmp)
	dns.Validate(c.Dns)

	return nil
}

func (c *AppConfig) GetAPIClient() *apiclient.Config {
	return c.ApiClient
}

/* config.go ends here. */
