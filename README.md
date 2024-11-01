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

See [autoderef/foo](_demo/autoderef/foo).

```go
type Point struct {
	X float64
	Y float64
}

func (p *Point) init(x, y float64) {
	p.X = x
	p.Y = y
}

func (p *Point) Print() {
	fmt.Printf("Point(%f, %f)\n", p.X, p.Y)
}

func (p *Point) Distance() float64 {
	return p.X * p.Y
}

// Move method for Point
func (p *Point) Move(dx, dy float64) {
	p.X += dx
	p.Y += dy
}

func Add(a, b int) int {
	return a + b
}

func InitFooModule() gp.Module {
	m := gp.CreateModule("foo")
	// Add the function to the module
	m.AddMethod("add", Add, "(a, b) -> float\n--\n\nAdd two integers.")
	// Add the type to the module
	gp.AddType[Point](m, (*Point).init, "Point", "Point objects")
	return m
}
```

Call foo module from Python and Go.

```go
package main

import (
	"fmt"
	"runtime"

	gp "github.com/cpunion/go-python"
	"github.com/cpunion/go-python/_demo/autoderef/foo"
	pymath "github.com/cpunion/go-python/math"
)

func Main1() {
	gp.RunString(`
import foo
point = foo.Point(3, 4)
print("dir(point):", dir(point))
print("x:", point.x)
print("y:", point.y)

print("distance:", point.distance())

point.move(1, 2)
print("x:", point.x)
print("y:", point.y)
print("distance:", point.distance())

point.print()
	`)
}

func Main2(fooMod gp.Module) {
	sum := fooMod.Call("add", gp.MakeLong(1), gp.MakeLong(2)).AsLong()
	fmt.Printf("Sum of 1 + 2: %d\n", sum.Int64())

	dict := fooMod.Dict()
	Point := dict.Get(gp.MakeStr("Point")).AsFunc()

	point := Point.Call(gp.MakeLong(3), gp.MakeLong(4))
	fmt.Printf("dir(point): %v\n", point.Dir())
	fmt.Printf("x: %v, y: %v\n", point.GetAttr("x"), point.GetAttr("y"))

	distance := point.Call("distance").AsFloat()
	fmt.Printf("Distance of 3 * 4: %f\n", distance.Float64())

	point.Call("move", gp.MakeFloat(1), gp.MakeFloat(2))
	fmt.Printf("x: %v, y: %v\n", point.GetAttr("x"), point.GetAttr("y"))

	distance = point.Call("distance").AsFloat()
	fmt.Printf("Distance of 4 * 6: %f\n", distance.Float64())
	point.Call("print")
}

func main() {
	gp.Initialize()
	defer gp.Finalize()
	fooMod := foo.InitFooModule()
	gp.GetModuleDict().Set(gp.MakeStr("foo").Object, fooMod.Object)

	Main1()
	Main2(fooMod)
}
```

### Call gradio

See [gradio](_demo/gradio).

```go
package main

import (
	"os"

	gp "github.com/cpunion/go-python"
)

/*
import gradio as gr

def update_examples(country):
		print("country:", country)
    if country == "USA":
        return gr.Dataset(samples=[["Chicago"], ["Little Rock"], ["San Francisco"]])
    else:
        return gr.Dataset(samples=[["Islamabad"], ["Karachi"], ["Lahore"]])

with gr.Blocks() as demo:
    dropdown = gr.Dropdown(label="Country", choices=["USA", "Pakistan"], value="USA")
    textbox = gr.Textbox()
    examples = gr.Examples([["Chicago"], ["Little Rock"], ["San Francisco"]], textbox)
    dropdown.change(update_examples, dropdown, examples.dataset)

demo.launch()
*/

var gr gp.Module

func UpdateExamples(country string) gp.Object {
	println("country:", country)
	if country == "USA" {
		return gr.Call("Dataset", gp.KwArgs{
			"samples": [][]string{{"Chicago"}, {"Little Rock"}, {"San Francisco"}},
		})
	} else {
		return gr.Call("Dataset", gp.KwArgs{
			"samples": [][]string{{"Islamabad"}, {"Karachi"}, {"Lahore"}},
		})
	}
}

func main() {
	if len(os.Args) > 2 {
		// avoid gradio start subprocesses
		return
	}

	gp.Initialize()
	gr = gp.ImportModule("gradio")
	fn := gp.CreateFunc(UpdateExamples,
		"(country, /)\n--\n\nUpdate examples based on country")
	// Would be (in the future):
	// fn := gp.FuncOf(UpdateExamples)
	demo := gp.With(gr.Call("Blocks"), func(v gp.Object) {
		dropdown := gr.Call("Dropdown", gp.KwArgs{
			"label":   "Country",
			"choices": []string{"USA", "Pakistan"},
			"value":   "USA",
		})
		textbox := gr.Call("Textbox")
		examples := gr.Call("Examples", [][]string{{"Chicago"}, {"Little Rock"}, {"San Francisco"}}, textbox)
		dataset := examples.GetAttr("dataset")
		dropdown.Call("change", fn, dropdown, dataset)
	})
	demo.Call("launch")
}
```

### Call matplotlib

See [plot](_demo/plot).

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
