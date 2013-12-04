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
	"github.com/prometheus/client_golang/model"
)

// Type for unserializing a generic Prometheus query response.
type QueryResponse struct {
	Type  string      `json:"Type"`
	Value interface{} `json:"Value"`
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
