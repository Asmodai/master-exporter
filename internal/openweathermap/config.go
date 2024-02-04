/*
 * config.go --- OpenWeatherMap config.
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

package openweathermap

import (
	"fmt"
	"sync"
)

type Config struct {
	initUrl sync.Once

	BaseUrl  string  `json:"base_url"`
	Version  float32 `json:"version"`
	Endpoint string  `json:"endpoint"`
	Location string  `json:"location"`
	Key      string  `json:"api_key"`
	Units    string  `json:"units"`
	Limit    int     `json:"limit"`
	Interval int     `json:"interval"`

	urlCache string
}

func NewDefaultConfig() *Config {
	return &Config{
		BaseUrl:  "",
		Version:  0,
		Endpoint: "",
		Location: "",
		Key:      "",
		Units:    "metric",
		Limit:    1000,
		Interval: 120,
	}
}

func Validate(cnf *Config) {
	if cnf.Limit == 0 {
		cnf.Limit = 1000
	}
}

func (c *Config) GetURL() string {
	c.initUrl.Do(func() {
		c.urlCache = fmt.Sprintf(
			"%s/data/%.1f/%s?q=%s&units=%s&appid=%s",
			c.BaseUrl,
			c.Version,
			c.Endpoint,
			c.Location,
			c.Units,
			c.Key,
		)
	})

	return c.urlCache
}

/* config.go ends here. */
