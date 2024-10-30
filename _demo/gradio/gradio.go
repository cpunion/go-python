package main

/*
#cgo pkg-config: python-3.12-embed
#include <Python.h>

extern PyObject* UpdateExamples2(PyObject* self, PyObject* args);
*/
import "C"

import (
	"os"

	"github.com/cpunion/go-python"
)

/*
import gradio as gr

def update_examples(country):
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

var gr python.Module

func UpdateExamples(country string) python.Object {
	println("country:", country)
	if country == "USA" {
		return gr.CallKeywords("Dataset")(python.MakeDict(map[any]any{
			"samples": [][]string{{"Chicago"}, {"Little Rock"}, {"San Francisco"}},
		}))
	} else {
		return gr.CallKeywords("Dataset")(python.MakeDict(map[any]any{
			"samples": [][]string{{"Islamabad"}, {"Karachi"}, {"Lahore"}},
		}))
	}
}

func main() {
	if len(os.Args) > 2 {
		// avoid gradio start subprocesses
		return
	}

	python.Initialize()
	gr = python.ImportModule("gradio")
	fn := python.FuncOf(UpdateExamples,
		"update_examples(country, /)\n--\n\nUpdate examples based on country")
	// fn := python.FuncOf1("update_examples", unsafe.Pointer(C.UpdateExamples2),
	// 	"update_examples(country, /)\n--\n\nUpdate examples based on country")
	// fn := python.FuncOf(UpdateExamples)
	blocks := gr.Call("Blocks")
	demo := python.With(blocks, func(v python.Object) {
		dropdown := gr.CallKeywords("Dropdown")(python.MakeDict(map[any]any{
			"label":   "Country",
			"choices": []string{"USA", "Pakistan"},
			"value":   "USA",
		}))
		textbox := gr.Call("Textbox")
		examples := gr.Call("Examples", [][]string{{"Chicago"}, {"Little Rock"}, {"San Francisco"}}, textbox)
		dataset := examples.GetAttr("dataset")
		dropdown.CallMethod("change", fn, dropdown, dataset)
	})
	demo.CallMethod("launch")
}
