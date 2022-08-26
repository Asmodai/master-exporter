/*
 * types.go --- Type hacks for this exporter.
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

package sabnzbd

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var (
	nonAlphaRE      = regexp.MustCompile(`[^a-zA-Z]+`)
	nonDigitRE      = regexp.MustCompile(`[^0-9]+`)
	nonDigitFloatRE = regexp.MustCompile(`[^0-9\+e\.]+`)
	timeRE          = regexp.MustCompile(`([0-9]+):([0-9]+):([0-9]+)`)
)

// ==================================================================
// {{{ Data type hacks:

// ------------------------------------------------------------------
// {{{ uint64 values:

type SabUint uint64

func (s *SabUint) UnmarshalJSON(b []byte) error {
	b = b[1 : len(b)-1]
	i, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return err
	}

	*s = SabUint(i)

	return nil
}

func (s SabUint) String() string   { return fmt.Sprintf("%d", uint64(s)) }
func (s SabUint) Uint64() uint64   { return uint64(s) }
func (s SabUint) Float64() float64 { return float64(uint64(s)) }
func (s SabUint) Int64() int64     { return int64(s) }

// }}}
// ------------------------------------------------------------------

// ------------------------------------------------------------------
// {{{ float64 values:

type SabFloat float64

func (s *SabFloat) UnmarshalJSON(b []byte) error {
	b = b[1 : len(b)-1]
	i, err := strconv.ParseFloat(string(b), 64)
	if err != nil {
		return err
	}

	*s = SabFloat(i)

	return nil
}

func (s SabFloat) String() string   { return fmt.Sprintf("%f", float64(s)) }
func (s SabFloat) Float64() float64 { return float64(s) }

// }}}
// ------------------------------------------------------------------

// ------------------------------------------------------------------
// {{{ Size values:

type SabSize float64

func (s *SabSize) UnmarshalJSON(b []byte) error {
	var i float64

	b = b[1 : len(b)-1]

	unit := nonAlphaRE.ReplaceAllString(string(b), "")
	num := nonDigitFloatRE.ReplaceAllString(string(b), "")

	i, err := strconv.ParseFloat(string(num), 64)
	if err != nil {
		return err
	}

	switch unit {
	case "M", "MB":
		i = i * 1.049e+6
	case "G", "GB":
		i = i * 1.074e+9
	case "T", "TB":
		i = i * 1.1e+12
	case "P", "PB":
		i = i * 1.126e+15
	}

	*s = SabSize(i)

	return nil
}

func (s SabSize) String() string   { return fmt.Sprintf("%f", float64(s)) }
func (s SabSize) Float64() float64 { return float64(s) }

// }}}
// ------------------------------------------------------------------

// ------------------------------------------------------------------
// {{{ Time duration value:

type SabDuration time.Duration

func (s *SabDuration) UnmarshalJSON(b []byte) error {
	b = b[1 : len(b)-1]

	dur := timeRE.ReplaceAllString(string(b), "${1}h${2}m${3}s")
	res, err := time.ParseDuration(dur)
	if err != nil {
		return err
	}

	*s = SabDuration(res)

	return nil
}

func (s SabDuration) String() string          { return time.Duration(s).String() }
func (s SabDuration) Duration() time.Duration { return time.Duration(s) }
func (s SabDuration) Float64() float64        { return float64(s) }
func (s SabDuration) Seconds() float64        { return float64(time.Duration(s).Seconds()) }

// }}}
// ------------------------------------------------------------------

// }}}
// ==================================================================

/* types.go ends here. */
