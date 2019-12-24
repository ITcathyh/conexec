## Introduction
conexec is a concurrent toolkit to help execute funcs concurrently in an efficient and safe way.It supports specifying the overall timeout to avoid blocking.

## How to use
Generally it can be set as a singleton to save memory.
```
	c := NewActuator()
	c.WithTimeOut(time.Second)
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