package gp

type KwArgs map[string]any

func splitArgs(args ...any) (Tuple, KwArgs) {
	if len(args) > 0 {
		last := args[len(args)-1]
		if kw, ok := last.(KwArgs); ok {
			return MakeTuple(args[:len(args)-1]...), kw
		}
	}
	return MakeTuple(args...), nil
}
