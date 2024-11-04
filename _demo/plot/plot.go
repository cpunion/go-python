package main

import gp "github.com/cpunion/go-python"

func main() {
	gp.Initialize()
	plt := gp.ImportModule("matplotlib.pyplot")
	plt.Call("plot", gp.MakeTuple(5, 10), gp.MakeTuple(10, 15), gp.KwArgs{"color": "red"})
	plt.Call("show")
}
