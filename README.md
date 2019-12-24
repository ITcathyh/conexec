## Introduction
[![Build Status](https://travis-ci.org/ITcathyh/conexec.svg?branch=master)](https://travis-ci.org/ITcathyh/conexec)
[![codecov](https://codecov.io/gh/ITcathyh/conexec/branch/master/graph/badge.svg)](https://codecov.io/gh/ITcathyh/conexec)
[![Go Report Card](https://goreportcard.com/badge/github.com/ITcathyh/conexec)](https://goreportcard.com/report/github.com/ITcathyh/conexec)
[![GoDoc](https://godoc.org/github.com/ITcathyh/conexec?status.svg)](https://godoc.org/github.com/ITcathyh/conexec)

conexec is a concurrent toolkit to help execute funcs concurrently in an efficient and safe way.It supports specifying the overall timeout to avoid blocking.

## How to use
Generally it can be set as a singleton to save memory.Here is the example to use it.
```
	c := NewActuator()
	c.WithTimeOut(time.Second) // set time out
	err := c.Exec(
		func() error {
			fmt.Println(1)
			time.Sleep(time.Second * 2)
			return nil
		},
		func() error {
			fmt.Println(2)
			return nil
		},
		func() error {
			time.Sleep(time.Second * 1)
			fmt.Println(3)
			return nil
		},
	)
	
	if err != nil {
		// ...do sth
	}
```