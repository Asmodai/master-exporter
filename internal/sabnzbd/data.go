/*
 * data.go --- SabNZBd data.
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

package sabnzbd

import (
	"fmt"
)

const (
	modeQueue  string = "queue"
	modeServer string = "server_stats"
)

type Server struct {
	Total uint64 `json:"total"`
}

func (s Server) Dump() {
	fmt.Printf("Server total:%d\n", s.Total)
}

type ServerStats struct {
	Total   uint64            `json:"total"`
	Servers map[string]Server `json:"servers"`
}

func NewServerStats() *ServerStats {
	return &ServerStats{
		Servers: map[string]Server{},
	}
}

type Queue struct {
	Version        string      `json:"version"`
	Paused         bool        `json:"paused"`
	PausedAll      bool        `json:"paused_all"`
	DiskSpace1     SabFloat    `json:"diskspace1"`
	DiskSpace2     SabFloat    `json:"diskspace2"`
	SpeedLimit     SabUint     `json:"speedlimit"`
	SpeedLimitAbs  SabFloat    `json:"speedlimit_abs"`
	CacheSize      SabSize     `json:"cache_size"`
	KbPerSec       SabFloat    `json:"kbpersec"`
	Speed          SabSize     `json:"speed"`
	MbLeft         SabFloat    `json:"mbleft"`
	Mb             SabFloat    `json:"mb"`
	SizeLeft       SabSize     `json:"sizeleft"`
	Size           SabSize     `json:"size"`
	TimeLeft       SabDuration `json:"timeleft"`
	NoOfSlots      uint64      `json:"noofslots"`
	NoOfSlotsTotal uint64      `json:"noofslots_total"`
}

func (s Queue) MbDone() float64 {
	return s.Mb.Float64() - s.MbLeft.Float64()
}

type QueueStats struct {
	Queue Queue `json:"queue"`
}

type SabNZBd struct {
	Queue  *QueueStats
	Server *ServerStats
}

func NewSabNZBd() *SabNZBd {
	return &SabNZBd{
		Queue:  &QueueStats{},
		Server: NewServerStats(),
	}
}

/* data.go ends here. */
