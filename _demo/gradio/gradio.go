package main

import (
	"os"

	gp "github.com/cpunion/go-python"
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

var gr gp.Module

func UpdateExamples(country string) gp.Object {
	println("country:", country)
	if country == "USA" {
		return gr.CallKeywords("Dataset")(gp.MakeDict(map[any]any{
			"samples": [][]string{{"Chicago"}, {"Little Rock"}, {"San Francisco"}},
		}))
	} else {
		return gr.CallKeywords("Dataset")(gp.MakeDict(map[any]any{
			"samples": [][]string{{"Islamabad"}, {"Karachi"}, {"Lahore"}},
		}))
	}
}

func main() {
	if len(os.Args) > 2 {
		// avoid gradio start subprocesses
		return
	}

	gp.Initialize()
	gr = gp.ImportModule("gradio")
	fn := gp.FuncOf(UpdateExamples,
		"update_examples(country, /)\n--\n\nUpdate examples based on country")
	// fn := gp.FuncOf(UpdateExamples)
	blocks := gr.Call("Blocks")
	demo := gp.With(blocks, func(v gp.Object) {
		dropdown := gr.CallKeywords("Dropdown")(gp.MakeDict(map[any]any{
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
