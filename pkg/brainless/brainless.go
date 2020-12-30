package brainless

func OrientatorRule(cell, example Orientator, f func(Compass, int), checker func(a, b Orientator) bool) {
	c := NewCompass()
	for i := Orientation(0); i < 8; i += 2 {
		if checker(cell, example.Orientate(i)) {
			for j := Orientation(0); j < 8; j += 2 {
				f(c.Orientate(i).Orientate(j).(Compass), j.GetRotations())
			}
		}
	}
}

func FlipperRule(flip Flipper, f func(Compass)) {
	c := NewCompass()
	for i := Orientation(0); i < 8; i++ {
		for j := 0; j < 2; j++ {
			f(c.Orientate(i).(Compass))
			flip.Flip()
		}
	}
}
