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
