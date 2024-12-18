# go-python: Write Python in Go - The most intuitive Python wrapper for Golang

[![Build Status](https://github.com/gotray/go-python/actions/workflows/go.yml/badge.svg)](https://github.com/gotray/go-python/actions/workflows/go.yml)
[![codecov](https://codecov.io/github/gotray/go-python/graph/badge.svg?token=TnaFHV1E3y)](https://codecov.io/github/gotray/go-python)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/cpunion/go-python)
[![GitHub commits](https://badgen.net/github/commits/cpunion/go-python)](https://GitHub.com/Naereen/cpunion/go-python/commit/)
[![GitHub release](https://img.shields.io/github/v/tag/cpunion/go-python.svg?label=release)](https://github.com/gotray/go-python/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/gotray/go-python)](https://goreportcard.com/report/github.com/gotray/go-python)
[![Go Reference](https://pkg.go.dev/badge/github.com/gotray/go-python.svg)](https://pkg.go.dev/github.com/gotray/go-python)

## Goal

- Automatically DecRef for Python objects.
- Typed Python objects.
- Most intuitive and easy to use.

## Features

- [x] Call Python in Go.
  - [x] Basic types.
    - [x] Integers.
    - [x] Floats.
    - [x] Strings.
    - [x] Bytes.
    - [x] Booleans.
    - [x] Lists.
    - [x] Tuples.
    - [x] Dicts.
  - [x] Modules.
  - [ ] Functions
    - [x] Keyword arguments.
    - [x] Mapping arguments and return values.
    - [ ] Default arguments.
- [x] Call Go in Python.
  - [x] Export Go functions, struct types to Python.
  - [x] Auto export field types if it's a struct.
  - [x] Auto generate function signatures (used by gradio).
- [x] Basic tests on common used libraries.
  - [x] matplotlib
  - [x] gradio

## Plans

- [x] Python virtual environment (https://github.com/gotray/got).
- [ ] Preprocess reference counting.
- [ ] True multi-threading.
- [ ] Support [LLGo](https://github.com/goplus/llgo).

## Examples

See the [examples](demo).

### Hello World: Plot a line

<!-- embedme demo/plot/plot.go -->

```go
package main

import . "github.com/gotray/go-python"

func main() {
	Initialize()
	plt := ImportModule("matplotlib.pyplot")
	plt.Call("plot", MakeTuple(5, 10), MakeTuple(10, 15), KwArgs{"color": "red"})
	plt.Call("show")
}

```

### Typed Python Objects

<!-- embedme demo/plot2/plot2.go -->

```go
package main

import . "github.com/gotray/go-python"

type plt struct {
	Module
}

func Plt() plt {
	return plt{ImportModule("matplotlib.pyplot")}
}

func (m plt) Plot(args ...any) Object {
	return m.Call("plot", args...)
}

func (m plt) Show() {
	m.Call("show")
}

func main() {
	Initialize()
	defer Finalize()
	plt := Plt()
	plt.Plot([]int{5, 10}, []int{10, 15}, KwArgs{"color": "red"})
	plt.Show()
}

```

### Define Python Objects with Go

<!-- embedme demo/module/foo/foo.go -->

```go
package foo

import (
	"fmt"

	. "github.com/gotray/go-python"
)

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

func InitFooModule() Module {
	m := CreateModule("foo")
	// Add the function to the module
	m.AddMethod("add", Add, "(a, b) -> float\n--\n\nAdd two integers.")
	// Add the type to the module
	m.AddType(Point{}, (*Point).init, "Point", "Point objects")
	return m
}

```

Call foo module from Python and Go.

<!-- embedme demo/module/module.go -->

```go
package main

import (
	"fmt"

	. "github.com/gotray/go-python"
	"github.com/gotray/go-python/demo/module/foo"
)

func main() {
	Initialize()
	defer Finalize()
	fooMod := foo.InitFooModule()
	GetModuleDict().SetString("foo", fooMod)

	Main1(fooMod)
	Main2()
}

func Main1(fooMod Module) {
	sum := fooMod.Call("add", 1, 2).AsLong()
	fmt.Printf("Sum of 1 + 2: %d\n", sum.Int64())

	dict := fooMod.Dict()
	Point := dict.Get(MakeStr("Point")).AsFunc()

	point := Point.Call(3, 4)
	fmt.Printf("dir(point): %v\n", point.Dir())
	fmt.Printf("x: %v, y: %v\n", point.Attr("x"), point.Attr("y"))

	distance := point.Call("distance").AsFloat()
	fmt.Printf("Distance of 3 * 4: %f\n", distance.Float64())

	point.Call("move", 1, 2)
	fmt.Printf("x: %v, y: %v\n", point.Attr("x"), point.Attr("y"))

	distance = point.Call("distance").AsFloat()
	fmt.Printf("Distance of 4 * 6: %f\n", distance.Float64())

	point.Call("print")
}

func Main2() {
	fmt.Printf("=========== Main2 ==========\n")
	_ = RunString(`
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

```

### Call gradio

<!-- embedme demo/gradio/gradio.go -->

```go
package main

import (
	"fmt"
	"os"

	. "github.com/gotray/go-python"
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

var gr Module

func updateExamples(country string) Object {
	println("country:", country)
	if country == "USA" {
		return gr.Call("Dataset", KwArgs{
			"samples": [][]string{{"Chicago"}, {"Little Rock"}, {"San Francisco"}},
		})
	} else {
		return gr.Call("Dataset", KwArgs{
			"samples": [][]string{{"Islamabad"}, {"Karachi"}, {"Lahore"}},
		})
	}
}

func main() {
	if len(os.Args) > 2 {
		// start subprocesses
		fmt.Println("start subprocess:", os.Args)
		os.Exit(RunMain(os.Args))
	}

	Initialize()
	defer Finalize()
	gr = ImportModule("gradio")
	demo := With(gr.Call("Blocks"), func(v Object) {
		dropdown := gr.Call("Dropdown", KwArgs{
			"label":   "Country",
			"choices": []string{"USA", "Pakistan"},
			"value":   "USA",
		})
		textbox := gr.Call("Textbox")
		examples := gr.Call("Examples", [][]string{{"Chicago"}, {"Little Rock"}, {"San Francisco"}}, textbox)
		dataset := examples.Attr("dataset")
		dropdown.Call("change", updateExamples, dropdown, dataset)
	})
	demo.Call("launch")
}

```
