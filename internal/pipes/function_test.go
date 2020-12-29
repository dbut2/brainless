package pipes

import (
	"testing"

	"github.com/puzzle-solvers/brainless/pkg/brainless"
	"github.com/stretchr/testify/assert"
)

func TestShapeChecker(t *testing.T) {
	type args struct {
		a brainless.Orientator
		b brainless.Orientator
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "TT",
			args: args{
				a: Shape(7),
				b: Shape(7),
			},
			want: true,
		},
		{
			name: "TL",
			args: args{
				a: Shape(7),
				b: Shape(3),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ShapeChecker(tt.args.a, tt.args.b))
		})
	}
}

func TestShape_Flip(t *testing.T) {
	type args struct {
		along brainless.Direction
	}
	tests := []struct {
		name string
		s    Shape
		args args
		want brainless.Orientator
	}{
		{
			name: "Test1",
			s:    Shape(3),
			args: args{
				along: brainless.East,
			},
			want: Shape(9),
		},
		{
			name: "Test2",
			s:    Shape(3),
			args: args{
				along: brainless.North,
			},
			want: Shape(6),
		},
		{
			name: "Test3",
			s:    Shape(3),
			args: args{
				along: brainless.West,
			},
			want: Shape(9),
		},
		{
			name: "Test4",
			s:    Shape(3),
			args: args{
				along: brainless.South,
			},
			want: Shape(6),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.s.Flip(tt.args.along))
		})
	}
}
