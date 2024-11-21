package math

import (
	gp "github.com/gotray/go-python"
)

var math_ gp.Module

func math() gp.Module {
	if math_.Nil() {
		math_ = gp.ImportModule("math")
	}
	return math_
}

func Sqrt(x gp.Float) gp.Float {
	return math().Call("sqrt", x).AsFloat()
}
