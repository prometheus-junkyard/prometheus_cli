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

// Prometheus HTTP API client functionality.
//
// TODO(julius): This functionality should be moved to a separate
// library/repository once we have a good name for it (client_* is already used
// for the interface between metrics-exposing servers and Prometheus).

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	queryPath      = "/api/query"
	queryRangePath = "/api/query_range"
	metricsPath    = "/api/metrics"

	scalarType = "scalar"
	vectorType = "vector"
	matrixType = "matrix"
	errorType  = "error"
)

// Client is a client for executing queries against the Prometheus API.
type Client struct {
	Endpoint   string
	httpClient http.Client
}

// transport builds a new transport with the provided timeout.
func transport(netw, addr string, timeout time.Duration) (connection net.Conn, err error) {
	deadline := time.Now().Add(timeout)
	connection, err = net.DialTimeout(netw, addr, timeout)
	if err == nil {
		connection.SetDeadline(deadline)
	}
	return
}

// NewClient creates a new Client, given a server URL and timeout.
func NewClient(url string, timeout time.Duration) *Client {
	return &Client{
		Endpoint: url,
		httpClient: http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) { return transport(netw, addr, timeout) },
			},
		},
	}
}

// Query performs an instant expression query via the Prometheus API.
func (c *Client) Query(expr string) (QueryResponse, error) {
	u, err := url.Parse(c.Endpoint)
	if err != nil {
		return nil, err
	}

	u.Path = queryPath
	q := u.Query()

	q.Set("expr", expr)
	u.RawQuery = q.Encode()

	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var r StubQueryResponse
	if err := json.Unmarshal(buf, &r); err != nil {
		return nil, err
	}

	var typedResp QueryResponse
	switch r.Type {
	case errorType:
		return nil, fmt.Errorf("query error: %s", r.Value.(string))
	case scalarType:
		typedResp = &ScalarQueryResponse{}
	case vectorType:
		typedResp = &VectorQueryResponse{}
	case matrixType:
		typedResp = &MatrixQueryResponse{}
	default:
		return nil, fmt.Errorf("invalid response type %s", r.Type)
	}

	if err := json.Unmarshal(buf, typedResp); err != nil {
		return nil, err
	}
	return typedResp, err
}

// QueryRange performs an range expression query via the Prometheus API.
func (c *Client) QueryRange(expr string, end uint64, rangeSec uint64, step uint64) (*MatrixQueryResponse, error) {
	u, err := url.Parse(c.Endpoint)
	if err != nil {
		return nil, err
	}

	u.Path = queryRangePath
	q := u.Query()

	q.Set("expr", expr)
	q.Set("end", fmt.Sprintf("%d", end))
	q.Set("range", fmt.Sprintf("%d", rangeSec))
	q.Set("step", fmt.Sprintf("%d", step))
	u.RawQuery = q.Encode()

	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var r StubQueryResponse
	if err := json.Unmarshal(buf, &r); err != nil {
		return nil, err
	}

	switch r.Type {
	case errorType:
		return nil, fmt.Errorf("query error: %s", r.Value.(string))
	case matrixType:
		var typedResp MatrixQueryResponse
		if err := json.Unmarshal(buf, &typedResp); err != nil {
			return nil, err
		}
		return &typedResp, nil
	default:
		return nil, fmt.Errorf("invalid response type %s", r.Type)
	}
}

// Metrics retrieves the list of available metric names via the Prometheus API.
func (c *Client) Metrics() ([]string, error) {
	u, err := url.Parse(c.Endpoint)
	if err != nil {
		return nil, err
	}

	u.Path = metricsPath

	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var r []string
	if err := json.Unmarshal(buf, &r); err != nil {
		return nil, err
	}
	return r, nil
}
