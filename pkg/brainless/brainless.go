package brainless

func OrientatorRule(cell, example Orientator, f func(Compass, int), checker func(Orientator, Orientator) bool) {
	c := NewCompass()
	for i := Orientation(0); i < 8; i += 2 {
		if checker(cell, example.Orientate(i)) {
			for j := Orientation(0); j < 8; j += 2 {
				f(c.Orientate(i).Orientate(j).(Compass), j.GetRotations())
			}
		}
	}
}

func AddRule(cell interface{}, checker interface{}, f func(Compass)) {
	o := NewCompass()
	for i := 1; i < 1<<4; i *= 2 {
		for j := 0; j < 2; j++ {
			if cell == checker {
				f(o)
			}
			o.Flip(Direction(i))
		}
		o.Rotate(1)
	}
}
