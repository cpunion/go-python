# go-python: a CPython wrapper for Go

Make Go and Python code inter-operable.

## Goal

- Provide automatically DecRef for Python objects.
- Wrap generic PyObject(s) to typed Python objects.
- Provide a way to define Python objects in Go.

## Usage

See the [examples](_demo).

### Hello World

```go
package main

import gp "github.com/cpunion/go-python"

func main() {
	gp.Initialize()
	plt := gp.ImportModule("matplotlib.pyplot")
	plt.Call("plot", gp.MakeTuple(5, 10), gp.MakeTuple(10, 15), gp.KwArgs{"color": "red"})
	plt.Call("show")
}
```

### Typed Python Objects

```go
package main

import gp "github.com/cpunion/go-python"

type plt struct {
	gp.Module
}

func Plt() plt {
	return plt{gp.ImportModule("matplotlib.pyplot")}
}

func (m plt) Plot(args ...any) gp.Object {
	return m.Call("plot", args...)
}

func (m plt) Show() {
	m.Call("show")
}

func main() {
	gp.Initialize()
	defer gp.Finalize()
	plt := Plt()
	plt.Plot(gp.MakeTuple(5, 10), gp.MakeTuple(10, 15), gp.KwArgs{"color": "red"})
	plt.Show()
}
```

### Define Python Objects

To be written.
