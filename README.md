# matrix

[![Build Status](https://github.com/milosgajdos/matrix/workflows/CI/badge.svg)](https://github.com/milosgajdos/matrix/actions?query=workflow%3ACI)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/milosgajdos/matrix)
[![GoDoc](https://godoc.org/github.com/milosgajdos/matrix?status.svg)](https://godoc.org/github.com/milosgajdos/matrix)
[![License](https://img.shields.io/:license-apache-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/milosgajdos/matrix)](https://goreportcard.com/report/github.com/milosgajdos/matrix)

This Go package contains various function useful when working with [gonum](https://www.gonum.org) matrices which are not provided by `gonum` [mat](https://godoc.org/gonum.org/v1/gonum/mat) package.

# Get started

Get the package:
```
$ go get -u github.com/milosgajdos/matrix
```

Get dependencies:
```
$ make godep && make dep
```

Run the tests:
```
$ make test
```

# Contributing

**YES PLEASE!**

Please make sure you run the following command before you open a new PR:
```shell
$ make all
```
