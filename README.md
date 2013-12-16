# Prometheus Command Line Interface

A command line interface tool for querying the Prometheus server's HTTP API.

## Building

    go build

## Usage

    Usage:
      ./prometheus_cli [flags] query <expression>
      ./prometheus_cli [flags] query_range <expression> <end_timestamp> <range_seconds> [<step_seconds>]
      ./prometheus_cli [flags] metrics
    
    Flags:
      -csv=true: Whether to format output as CSV
      -csvDelimiter=";": Single-character delimiter to use in CSV output
      -server="": URL of the Prometheus server to query
      -timeout=1m0s: Timeout to use when querying the Prometheus server
