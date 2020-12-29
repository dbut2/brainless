package brainless

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrientate(t *testing.T) {
	type args struct {
		o           Orientator
		orientation Orientation
	}
	tests := []struct {
		name string
		args args
		want Orientator
	}{
		{
			name: "Test0",
			args: args{
				o:           Shape(3),
				orientation: Orientation(0),
			},
			want: Shape(3),
		},
		{
			name: "Test1",
			args: args{
				o:           Shape(3),
				orientation: Orientation(1),
			},
			want: Shape(9),
		},
		{
			name: "Test2",
			args: args{
				o:           Shape(3),
				orientation: Orientation(2),
			},
			want: Shape(6),
		},
		{
			name: "Test3",
			args: args{
				o:           Shape(3),
				orientation: Orientation(3),
			},
			want: Shape(3),
		},
		{
			name: "Test4",
			args: args{
				o:           Shape(3),
				orientation: Orientation(4),
			},
			want: Shape(12),
		},
		{
			name: "Test5",
			args: args{
				o:           Shape(3),
				orientation: Orientation(5),
			},
			want: Shape(6),
		},
		{
			name: "Test6",
			args: args{
				o:           Shape(3),
				orientation: Orientation(6),
			},
			want: Shape(9),
		},
		{
			name: "Test7",
			args: args{
				o:           Shape(3),
				orientation: Orientation(7),
			},
			want: Shape(12),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Orientate(tt.args.o, tt.args.orientation))
		})
	}
}
