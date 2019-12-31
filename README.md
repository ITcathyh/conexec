## Introduction
[![Build Status](https://travis-ci.org/ITcathyh/conexec.svg?branch=master)](https://travis-ci.org/ITcathyh/conexec)
[![codecov](https://codecov.io/gh/ITcathyh/conexec/branch/master/graph/badge.svg)](https://codecov.io/gh/ITcathyh/conexec)
[![Go Report Card](https://goreportcard.com/badge/github.com/ITcathyh/conexec)](https://goreportcard.com/report/github.com/ITcathyh/conexec)
[![GoDoc](https://godoc.org/github.com/ITcathyh/conexec?status.svg)](https://godoc.org/github.com/ITcathyh/conexec)

conexec is a concurrent toolkit to help execute functions concurrently in an efficient and safe way.It supports specifying the overall timeout to avoid blocking.

## How to use
Generally it can be set as a singleton to save memory.Here is the example to use it.
### Normal Actuator
Actuator is a base struct to execute functions concurrently.
```
	opt := &Options{TimeOut:DurationPtr(time.Millisecond*50)}
	c := NewActuator(opt)
	
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
### Pooled Actuator
Pooled actuator uses the goroutine pool to execute functions.In some times it is a more efficient way.
```
	opt := &Options{TimeOut:DurationPtr(time.Millisecond*50)}
	c := NewPooledActuator(5, opt)
	
	err := c.Exec(...)
	
	if err != nil {
		// ...do sth
	}
```
### Simply exec using goroutine
```
	done := Exec(...)

	if !done {
		// ... do sth 
	}
```