/*
 * data.go --- Data.
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
	"time"
)

type ReportCoord struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lon"`
}

type ReportConditions struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type ReportMain struct {
	Temperature    float32 `json:"temp"`
	FeelsLike      float32 `json:"feels_like"`
	TemperatureMin float32 `json:"temp_min"`
	TemperatureMax float32 `json:"temp_max"`
	AirPressure    int     `json:"pressure"`
	SeaLevel       int     `json:"sea_level"`
	GroundLevel    int     `json:"grnd_level"`
	Humidity       int     `json:"humidity"`
}

type ReportWind struct {
	Speed     float32 `json:"speed"`
	Direction int     `json:"deg"`
	Gust      float32 `json:"gust"`
}

type ReportClouds struct {
	Coverage int `json:"all"`
}

type ReportPrecipitation struct {
	Volume1h float32 `json:"1h"`
	Volume3h float32 `json:"3h"`
}

type ReportSystem struct {
	Type    int    `json:"type"`
	Id      int    `json:"id"`
	Message string `json:"message"`
	Country string `json:"country"`
	Sunrise int    `json:"sunrise"`
	Sunset  int    `json:"sunset"`
}

type OpenWeatherMap struct {
	Timestamp int64  `json:"dt"`
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Base      string `json:"base"`

	Coords     ReportCoord         `json:"coord"`
	Conditions []ReportConditions  `json:"weather"`
	Main       ReportMain          `json:"main"`
	Visibility int                 `json:"visibility"`
	Wind       ReportWind          `json:"wind"`
	Clouds     ReportClouds        `json:"clouds"`
	Rain       ReportPrecipitation `json:"rain"`
	Snow       ReportPrecipitation `json:"snow"`
	System     ReportSystem        `json:"sys"`
}

func NewOpenWeatherMap() *OpenWeatherMap {
	return &OpenWeatherMap{
		Conditions: []ReportConditions{},
	}
}

func (o *OpenWeatherMap) ToTime() time.Time {
	return time.Unix(o.Timestamp, 0)
}

/* data.go ends here. */
