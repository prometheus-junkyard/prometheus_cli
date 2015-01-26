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

// Interface for query results of various result types.
type QueryResponse interface {
	ToText() string
	ToCSV(delim rune) string
}

// Type for unserializing a generic Prometheus query response.
type StubQueryResponse struct {
	Type    string      `json:"type"`
	Value   interface{} `json:"value"`
	Version int         `json:"version"`
}

// Type for unserializing a scalar-typed Prometheus query response.
type ScalarQueryResponse struct {
	Value string `json:"value"`
}

// Type for unserializing a vector-typed Prometheus query response.
type VectorQueryResponse struct {
	Value []struct {
		Metric    model.Metric `json:"metric"`
		Value     string       `json:"value"`
		Timestamp float64      `json:"timestamp"`
	} `json:"value"`
}

// Type for unserializing a matrix-typed Prometheus query response.
type MatrixQueryResponse struct {
	Value []struct {
		Metric model.Metric `json:"metric"`
		Values [][]interface{}
	} `json:"value"`
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
		lines = append(lines, fmt.Sprintf("%s %s@%.3f\n", v.Metric, v.Value, v.Timestamp))
	}
	return strings.Join(lines, "")
}

func (r VectorQueryResponse) ToCSV(delim rune) string {
	rows := make([][]string, 0, len(r.Value))
	for _, v := range r.Value {
		rows = append(rows, []string{
			v.Metric.String(),
			v.Value,
			strconv.FormatFloat(v.Timestamp, 'f', -1, 64),
		})
	}
	return formatCSV(rows, delim)
}

func (r MatrixQueryResponse) ToText() string {
	lines := make([]string, 0, len(r.Value))
	for _, v := range r.Value {
		vals := make([]string, 0, len(v.Values))
		for _, s := range v.Values {
			vals = append(vals, fmt.Sprintf("%s@%.3f ", s[1], s[0]))
		}
		lines = append(lines, fmt.Sprintf("%s %s\n", v.Metric, strings.Join(vals, " ")))
	}
	return strings.Join(lines, "")
}

func (r MatrixQueryResponse) ToCSV(delim rune) string {
	rows := make([][]string, 0, len(r.Value))
	for _, v := range r.Value {
		vals := make([]string, 0, len(v.Values))
		for _, s := range v.Values {
			vals = append(vals, fmt.Sprintf("%s@%.3f", s[1], s[0]))
		}
		rows = append(rows, []string{
			v.Metric.String(),
			strings.Join(vals, " "),
		})
	}
	return formatCSV(rows, delim)
}
