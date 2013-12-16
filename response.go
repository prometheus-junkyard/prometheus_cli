// Copyright 2013 Prometheus Team
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/model"
)

// Type for unserializing a generic Prometheus query response.
type StubQueryResponse struct {
	Type  string      `json:"Type"`
	Value interface{} `json:"Value"`
}

// Interface for query results of various result types.
type QueryResponse interface {
	ToText() string
	ToCSV(delim rune) string
}

// Type for unserializing a scalar-typed Prometheus query response.
type ScalarQueryResponse struct {
	Value string `json:"Value"`
}

// Type for unserializing a vector-typed Prometheus query response.
type VectorQueryResponse struct {
	Value []struct {
		Metric    model.Metric `json:"Metric"`
		Value     string       `json:"Value"`
		Timestamp int64        `json:"Timestamp"`
	} `json:"Value"`
}

// Type for unserializing a matrix-typed Prometheus query response.
type MatrixQueryResponse struct {
	Value []struct {
		Metric model.Metric `json:"Metric"`
		Values []struct {
			Value     string `json:"Value"`
			Timestamp int64  `json:"Timestamp"`
		}
	} `json:"Value"`
}

func (r ScalarQueryResponse) ToText() string {
	return fmt.Sprint(r.Value)
}

func formatCSV(rows [][]string, delim rune) string {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	w.Comma = delim
	for _, row := range rows {
		w.Write(row)
		if err := w.Error(); err != nil {
			panic("error formatting CSV")
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		panic("error dumping CSV")
	}
	return buf.String()
}

func (r ScalarQueryResponse) ToCSV(delim rune) string {
	return formatCSV([][]string{{r.Value}}, delim)
}

func (r VectorQueryResponse) ToText() string {
	lines := make([]string, 0, len(r.Value))
	for _, v := range r.Value {
		lines = append(lines, fmt.Sprintf("%s %s@%d", v.Metric, v.Value, v.Timestamp))
	}
	return strings.Join(lines, "\n")
}

func (r VectorQueryResponse) ToCSV(delim rune) string {
	rows := make([][]string, 0, len(r.Value))
	for _, v := range r.Value {
		rows = append(rows, []string{
			v.Metric.String(),
			v.Value,
			strconv.FormatInt(v.Timestamp, 10),
		})
	}
	return formatCSV(rows, delim)
}

func (r MatrixQueryResponse) ToText() string {
	lines := make([]string, 0, len(r.Value))
	for _, v := range r.Value {
		vals := make([]string, 0, len(v.Values))
		for _, s := range v.Values {
			vals = append(vals, fmt.Sprintf("%s@%d ", s.Value, s.Timestamp))
		}
		lines = append(lines, fmt.Sprintf("%s %s", v.Metric, strings.Join(vals, " ")))
	}
	return strings.Join(lines, "\n")
}

func (r MatrixQueryResponse) ToCSV(delim rune) string {
	rows := make([][]string, 0, len(r.Value))
	for _, v := range r.Value {
		vals := make([]string, 0, len(v.Values))
		for _, s := range v.Values {
			vals = append(vals, fmt.Sprintf("%s@%d", s.Value, s.Timestamp))
		}
		rows = append(rows, []string{
			v.Metric.String(),
			strings.Join(vals, " "),
		})
	}
	return formatCSV(rows, delim)
}
