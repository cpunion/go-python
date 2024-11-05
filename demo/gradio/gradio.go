package main

import (
	"os"

	. "github.com/cpunion/go-python"
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

func UpdateExamples(country string) Object {
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
		// avoid gradio start subprocesses
		return
	}

	Initialize()
	defer Finalize()
	gr = ImportModule("gradio")
	fn := CreateFunc("update_examples", UpdateExamples,
		"(country, /)\n--\n\nUpdate examples based on country")
	// Would be (in the future):
	// fn := FuncOf(UpdateExamples)
	demo := With(gr.Call("Blocks"), func(v Object) {
		dropdown := gr.Call("Dropdown", KwArgs{
			"label":   "Country",
			"choices": []string{"USA", "Pakistan"},
			"value":   "USA",
		})
		textbox := gr.Call("Textbox")
		examples := gr.Call("Examples", [][]string{{"Chicago"}, {"Little Rock"}, {"San Francisco"}}, textbox)
		dataset := examples.Attr("dataset")
		dropdown.Call("change", fn, dropdown, dataset)
	})
	demo.Call("launch")
}
