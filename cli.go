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
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

var (
	server  = flag.String("server", "", "URL of the Prometheus server to query")
	timeout = flag.Duration("timeout", time.Minute, "Timeout to use when querying the Prometheus server")
)

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintln(os.Stderr, "")
	os.Exit(1)
}

func query(c *Client) {
	if flag.NArg() != 2 {
		flag.Usage()
		die("Please supply a query expression")
	}

	resp, err := c.Query(flag.Arg(1))
	if err != nil {
		die("Error querying server: %s", err)
	}

	switch resp.(type) {
	case *ScalarQueryResponse:
		fmt.Println(resp.(*ScalarQueryResponse).Value)
	case *VectorQueryResponse:
		for _, v := range resp.(*VectorQueryResponse).Value {
			fmt.Printf("%s %s@%d\n", v.Metric, v.Value, v.Timestamp)
		}
	case *MatrixQueryResponse:
		for _, v := range resp.(*MatrixQueryResponse).Value {
			fmt.Printf("%s ", v.Metric)
			for _, s := range v.Values {
				fmt.Printf("%s@%d ", s.Value, s.Timestamp)
			}
			fmt.Println("")
		}
	default:
		die("Unknown query response type")
	}
}

func queryRange(c *Client) {
	if flag.NArg() != 4 && flag.NArg() != 5 {
		flag.Usage()
		die("Wrong number of range query arguments")
	}
	end, err := strconv.ParseUint(flag.Arg(2), 10, 64)
	if err != nil {
		flag.Usage()
		die("Invalid end timestamp '%s'", flag.Arg(2))
	}
	rangeSec, err := strconv.ParseUint(flag.Arg(3), 10, 64)
	if err != nil {
		flag.Usage()
		die("Invalid query range '%s'", flag.Arg(3))
	}
	var step uint64
	if flag.NArg() == 5 {
		step, err = strconv.ParseUint(flag.Arg(4), 10, 64)
		if err != nil {
			flag.Usage()
			die("Invalid query step '%s'", flag.Arg(4))
		}
	} else {
		step = rangeSec / 250
	}
	if step < 1 {
		step = 1
	}

	resp, err := c.QueryRange(flag.Arg(1), end, rangeSec, step)
	if err != nil {
		die("Error querying server: %s", err)
	}

	for _, v := range resp.Value {
		fmt.Printf("%s ", v.Metric)
		for _, s := range v.Values {
			fmt.Printf("%s@%d ", s.Value, s.Timestamp)
		}
		fmt.Println("")
	}
}

func metrics(c *Client) {
	if flag.NArg() != 1 {
		flag.Usage()
		die("Too many arguments")
	}

	if metrics, err := c.Metrics(); err != nil {
		die("Error querying server: %s", err)
	} else {
		for _, m := range metrics {
			fmt.Println(m)
		}
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "\t%s [flags] query <expression>\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\t%s [flags] query_range <expression> <end_timestamp> <range_seconds> [<step_seconds>]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\t%s [flags] metrics\n", os.Args[0])
	fmt.Printf("\nFlags:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *server == "" {
		flag.Usage()
		die("Please provide a server name.")
	}
	if flag.NArg() < 1 {
		flag.Usage()
		die("Please provide a command.")
	}

	c := NewClient(*server, *timeout)
	switch flag.Arg(0) {
	case "query":
		query(c)
	case "query_range":
		queryRange(c)
	case "metrics":
		metrics(c)
	default:
		die("Unknown command '%s'", flag.Arg(0))
	}
}
