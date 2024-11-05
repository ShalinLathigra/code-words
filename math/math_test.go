package math

import (
	"reflect"
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
			if got := min(tt.args.a, tt.args.b); got != tt.want.min {
				t.Errorf("min(%v) = %v, want %v", tt.args, got, tt.want.min)
			}
			if got := max(tt.args.a, tt.args.b); got != tt.want.max {
				t.Errorf("max(%v) = %v, want %v", tt.args, got, tt.want.max)
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

			if got := clamp(tt.args.v, tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("clamp(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

func Test_Contains(t *testing.T) {

	testRect := Rect{
		Vec2{-100, -100}, Vec2{201, 201},
	}
	type args struct {
		r Rect
		p Vec2
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"p < minR",
			args{
				testRect,
				Vec2{-150, -150},
			},
			false,
		},
		{
			"p > maxR",
			args{
				testRect,
				Vec2{150, 150},
			},
			false,
		},
		{
			"p <= maxR",
			args{
				testRect,
				Vec2{100, 100},
			},
			true,
		},
		{
			"x,y mix",
			args{
				testRect,
				Vec2{-50, 150},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.args.r, tt.args.p); got != tt.want {
				t.Errorf("contains(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

func Test_Overlap(t *testing.T) {

	type args struct {
		r Rect
		o Rect
	}
	type want struct {
		inter Rect
		ok    bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"perfect",
			args{
				Rect{
					Vec2{-10, -10}, Vec2{20, 20},
				},
				Rect{
					Vec2{-10, -10}, Vec2{20, 20},
				},
			},
			want{
				Rect{
					Vec2{-10, -10},
					Vec2{20, 20},
				},
				true,
			},
		},
		{
			"rect oversize",
			args{
				Rect{
					Vec2{-15, -15}, Vec2{30, 30},
				},
				Rect{
					Vec2{-10, -10}, Vec2{20, 20},
				},
			},
			want{
				Rect{
					Vec2{-10, -10},
					Vec2{20, 20},
				},
				true,
			},
		},
		{
			"rect undersize",
			args{
				Rect{
					Vec2{-5, -5}, Vec2{10, 10},
				},
				Rect{
					Vec2{-10, -10}, Vec2{20, 20},
				},
			},
			want{
				Rect{
					Vec2{-5, -5},
					Vec2{10, 10},
				},
				true,
			},
		},
		{
			"topleft",
			args{
				Rect{
					Vec2{-10, -10}, Vec2{20, 20},
				},
				Rect{
					Vec2{-15, -15}, Vec2{10, 10},
				},
			},
			want{
				Rect{
					Vec2{-10, -10},
					Vec2{5, 5},
				},
				true,
			},
		},
		{
			"botright",
			args{
				Rect{
					Vec2{-10, -10}, Vec2{20, 20},
				},
				Rect{
					Vec2{5, 5}, Vec2{10, 10},
				},
			},
			want{
				Rect{
					Vec2{5, 5},
					Vec2{5, 5},
				},
				true,
			},
		},
		{
			"Miss TopRight",
			args{
				Rect{
					Vec2{-10, -10}, Vec2{20, 20},
				},
				Rect{
					Vec2{15, 15}, Vec2{10, 10},
				},
			},
			want{
				Rect{},
				false,
			},
		},
		{
			"Miss BotLeft",
			args{
				Rect{
					Vec2{-10, -10}, Vec2{20, 20},
				},
				Rect{
					Vec2{-20, 20}, Vec2{10, 10},
				},
			},
			want{
				Rect{},
				false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, ok := overlap(tt.args.r, tt.args.o); ok != tt.want.ok || !reflect.DeepEqual(got, tt.want.inter) {
				t.Errorf("overlap(%v) = (%s, %v), want %v", tt.args, got, ok, tt.want)
			}
		})
	}
}
