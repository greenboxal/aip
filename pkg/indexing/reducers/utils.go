package reducers

func Grow[S ~[]E, E any](s S, n int) S {
	if n < 0 {
		panic("cannot be negative")
	}

	t := len(s) + n

	if t < cap(s) {
		s = s[:t]
	} else {
		newArr := make([]E, t, t*3/2+1)
		copy(newArr, s)
		s = newArr[:t]
	}

	return s
}
