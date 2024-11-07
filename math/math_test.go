package math

import (
	"testing"
)

func Test_minmax(t *testing.T) {
	type args struct {
		a int
		b int
	}
	type want struct {
		min int
		max int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"a<b",
			args{
				-100, 100,
			},
			want{
				-100,
				100,
			},
		},
		{
			"a==b",
			args{
				100, 100,
			},
			want{
				100,
				100,
			},
		},
		{
			"a>b",
			args{
				100, -100,
			},
			want{
				-100,
				100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.args.a, tt.args.b); got != tt.want.min {
				t.Errorf("Min(%v) = %v, want %v", tt.args, got, tt.want.min)
			}
			if got := Max(tt.args.a, tt.args.b); got != tt.want.max {
				t.Errorf("Max(%v) = %v, want %v", tt.args, got, tt.want.max)
			}
		})
	}
}

func Test_Clamp(t *testing.T) {
	type args struct {
		v int
		a int
		b int
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"v<a<b",
			args{
				-100, 100, 200,
			},
			100,
		},
		{
			"a<v<b",
			args{
				100, -100, 200,
			},
			100,
		},
		{
			"a<b<v",
			args{
				200, -100, 100,
			},
			100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := Clamp(tt.args.v, tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Clamp(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}
