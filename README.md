# enid

- A very simple unique id generator.
- Methods to parse existing enid ids.
- Methods to convert a enid id into several other data types and back.
- JSON Marshal/Unmarshal functions to easily use enid ids within a JSON API.
- Monotonic Clock calculations protect from clock drift.
- Optional use entropy.

[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/thinkgos/enid?tab=doc)
[![codecov](https://codecov.io/gh/thinkgos/enid/branch/main/graph/badge.svg)](https://codecov.io/gh/thinkgos/enid)
[![Tests](https://github.com/thinkgos/enid/actions/workflows/ci.yml/badge.svg)](https://github.com/thinkgos/enid/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/thinkgos/enid)](https://goreportcard.com/report/github.com/thinkgos/enid)
[![License](https://img.shields.io/github/license/thinkgos/enid)](https://github.com/thinkgos/enid/raw/main/LICENSE)
[![Tag](https://img.shields.io/github/v/tag/thinkgos/enid)](https://github.com/thinkgos/enid/tags)

## Specification

```text
+---------------------------------------------------------------------------------------+
| 1 Bit Unused | 43 Bit Timestamp |   12 Bit Sequence Id  | 8 Bit NodeId Or rand number |
+---------------------------------------------------------------------------------------+
```

By default, the id format follows:

- The id as a whole is a 63 bit integer stored in an int64
- 43 bits are used to store a timestamp with millisecond precision, using a custom epoch. default Dec 01 2024 05:06:07 UTC in milliseconds. The lifetime almost 278 year
- 12 bits are used to store a sequence number - a range from 0 through 4095.
- 8 bits are used to store a node id or rand number - a range from 0 through 255.

## Getting Started

### Installation

Use go get.

```sh
    go get github.com/thinkgos/enid
```

Then import the package into your own code.

```sh
    import "github.com/thinkgos/enid"
```

### Example

```go

```

## Performance

To benchmark the generator on your system run the following command inside the snowflake package directory.

```sh
go test -run=^$ -bench=.
```

## Reference

- [snowflake](https://github.com/bwmarrin/snowflake)
- [sonyflake](https://github.com/sony/sonyflake)
