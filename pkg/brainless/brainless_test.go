package brainless

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrientatorRule(t *testing.T) {
	type args struct {
		cell    Orientator
		checker Orientator
		f       func(*Compass, int)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				cell:    Shape(1101),
				checker: T,
				f: func(o *Compass, r int) {
					assert.True(t, false)
				},
			},
		},
	}
	for _, tt := range tests {
		_ = tt
	}
}

type Shape int

const (
	P Shape = 1
	L Shape = 3
	I Shape = 5
	T Shape = 7
)

func (s Shape) Rotate(times int) Orientator {
	return s.rotate(times)
}

func (s Shape) rotate(times int) Shape {
	return Shape((s<<times + s<<times>>4) & 15)
}

func (s Shape) Flip(along Direction) Orientator {
	return s.flip(along)
}

func (s Shape) flip(along Direction) Shape {
	axis := int(along + along.Flip())
	return Shape(axis&int(s) + ^axis&int(s.rotate(2)))
}
