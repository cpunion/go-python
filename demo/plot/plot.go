package main

import . "github.com/cpunion/go-python"

func main() {
	Initialize()
	plt := ImportModule("matplotlib.pyplot")
	plt.Call("plot", MakeTuple(5, 10), MakeTuple(10, 15), KwArgs{"color": "red"})
	plt.Call("show")
}
