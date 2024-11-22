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
