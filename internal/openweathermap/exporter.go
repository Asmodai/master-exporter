/*
 * exporter.go --- openweathermap exporter.
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
	data    *OpenWeatherMap
	metrics *Metrics
	calls   int
}

func NewExporter(client apiclient.IApiClient, logger logger.ILogger, config *Config) *Exporter {
	return &Exporter{
		client:  client,
		logger:  logger,
		config:  config,
		data:    NewOpenWeatherMap(),
		metrics: NewMetrics(config.Location),
		calls:   0,
	}
}

func (e *Exporter) Data() *OpenWeatherMap {
	return e.data
}

func (e *Exporter) canCall() bool { return e.calls < e.config.Limit }

func (e *Exporter) resetLimit() {
	now := time.Now()

	if now.Hour() == 0 && now.Minute() == 0 {
		if e.calls > 1 {
			e.calls = 0
		}
	}
}

func (e *Exporter) get() error {
	e.resetLimit()

	if !e.canCall() {
		return fmt.Errorf("Daily call limit exceeded.")
	}

	params := &apiclient.Params{
		Url: e.config.GetURL(),
	}

	data, code, err := e.client.Get(params)
	if err != nil {
		return err
	}
	e.calls++

	e.metrics.SetCalls(float64(e.calls))
	e.metrics.SetLimit(float64(e.config.Limit))

	switch code {
	case 200:
		{
			err := json.Unmarshal(data, e.data)
			if err != nil {
				return fmt.Errorf("JSON unmarshal: %s", err)
			}
		}

	default:
		return fmt.Errorf("%d response: %s", code, err)
	}

	return nil
}

func (e *Exporter) Interval() int {
	return e.config.Interval
}

func (e *Exporter) Scrape() error {
	if err := e.get(); err != nil {
		e.logger.Warn(
			"Scrape error.",
			"err", err.Error(),
			"exporter", "openweathermap",
		)
	}

	e.metrics.SetTemp(float64(e.data.Main.Temperature))
	e.metrics.SetTempFeelsLike(float64(e.data.Main.FeelsLike))
	e.metrics.SetTempMin(float64(e.data.Main.TemperatureMin))
	e.metrics.SetTempMax(float64(e.data.Main.TemperatureMax))
	e.metrics.SetAirPressure(float64(e.data.Main.AirPressure))
	e.metrics.SetHumidity(float64(e.data.Main.Humidity))
	e.metrics.SetWindSpeed(float64(e.data.Wind.Speed))
	e.metrics.SetWindGust(float64(e.data.Wind.Gust))
	e.metrics.SetWindDirection(float64(e.data.Wind.Direction))
	e.metrics.SetRainLevel(float64(e.data.Rain.Volume1h))
	e.metrics.SetSnowLevel(float64(e.data.Snow.Volume1h))
	e.metrics.SetVisibility(float64(e.data.Visibility))
	e.metrics.SetCloudCover(float64(e.data.Clouds.Coverage))

	return nil
}

/* exporter.go ends here. */
