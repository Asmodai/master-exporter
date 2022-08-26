/*
 * metrics.go --- Prometheus metrics.
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

package openweathermap

import (
	"github.com/Asmodai/master-exporter/internal/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	Temp          prometheus.Gauge
	TempFeelsLike prometheus.Gauge
	TempMax       prometheus.Gauge
	TempMin       prometheus.Gauge
	AirPressure   prometheus.Gauge
	Humidity      prometheus.Gauge
	RainLevel     prometheus.Gauge
	SnowLevel     prometheus.Gauge
	WindSpeed     prometheus.Gauge
	WindGust      prometheus.Gauge
	WindDirection prometheus.Gauge
	Visibility    prometheus.Gauge
	CloudCover    prometheus.Gauge

	Calls prometheus.Gauge
	Limit prometheus.Gauge
}

func NewGauge(name string, site string) prometheus.Gauge {
	return promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "weather",
		Name:      name,
		Help:      "Weather data.",
		ConstLabels: map[string]string{
			"location": site,
		},
	})
}

func NewMetrics(site string) *Metrics {
	return &Metrics{
		Temp:          NewGauge("temp", site),
		TempFeelsLike: NewGauge("temp_feels_like", site),
		TempMax:       NewGauge("temp_max", site),
		TempMin:       NewGauge("temp_min", site),
		AirPressure:   NewGauge("air_pressure", site),
		Humidity:      NewGauge("humidity", site),
		RainLevel:     NewGauge("rain_level", site),
		SnowLevel:     NewGauge("snow_level", site),
		WindSpeed:     NewGauge("wind_speed", site),
		WindGust:      NewGauge("wind_gust", site),
		WindDirection: NewGauge("wind_direction", site),
		Visibility:    NewGauge("visibility", site),
		CloudCover:    NewGauge("cloud_cover", site),

		Calls: metrics.NewMetricsGauge("calls", "weather"),
		Limit: metrics.NewMetricsGauge("limit", "weather"),
	}
}

func (m *Metrics) SetTemp(val float64)          { m.Temp.Set(val) }
func (m *Metrics) SetTempFeelsLike(val float64) { m.TempFeelsLike.Set(val) }
func (m *Metrics) SetTempMax(val float64)       { m.TempMax.Set(val) }
func (m *Metrics) SetTempMin(val float64)       { m.TempMin.Set(val) }
func (m *Metrics) SetHumidity(val float64)      { m.Humidity.Set(val) }
func (m *Metrics) SetAirPressure(val float64)   { m.AirPressure.Set(val) }
func (m *Metrics) SetRainLevel(val float64)     { m.RainLevel.Set(val) }
func (m *Metrics) SetSnowLevel(val float64)     { m.SnowLevel.Set(val) }
func (m *Metrics) SetWindSpeed(val float64)     { m.WindSpeed.Set(val) }
func (m *Metrics) SetWindGust(val float64)      { m.WindGust.Set(val) }
func (m *Metrics) SetWindDirection(val float64) { m.WindDirection.Set(val) }
func (m *Metrics) SetVisibility(val float64)    { m.Visibility.Set(val) }
func (m *Metrics) SetCloudCover(val float64)    { m.CloudCover.Set(val) }

func (m *Metrics) SetCalls(val float64) { m.Calls.Set(val) }
func (m *Metrics) SetLimit(val float64) { m.Limit.Set(val) }

/* metrics.go ends here. */
