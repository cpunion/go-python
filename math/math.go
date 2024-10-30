package math

import (
	gp "github.com/cpunion/go-python"
)

var math_ gp.Module

func math() gp.Module {
	if math_.Nil() {
		math_ = gp.ImportModule("math")
	}
	return math_
}

func Sqrt(x gp.Float) gp.Float {
	return math().CallMethod("sqrt", x.Obj()).AsFloat()
}
