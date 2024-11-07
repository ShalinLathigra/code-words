package math

import (
	"reflect"
	"testing"
)

func Test_Vec2(t *testing.T) {
	type args struct {
		a Vec2
		b Vec2
	}
	type want struct {
		add      Vec2
		subtract Vec2
		dot      int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"Zero'd",
			args{
				Vec2{},
				Vec2{},
			},
			want{
				add:      Vec2{},
				subtract: Vec2{},
				dot:      0,
			},
		},
		{
			"+ to -",
			args{
				Vec2{100, 100},
				Vec2{-200, -200},
			},
			want{
				add:      Vec2{-100, -100},
				subtract: Vec2{300, 300},
				dot:      -40000,
			},
		},
		{
			"- to +",
			args{
				Vec2{-199, -12},
				Vec2{20, 200},
			},
			want{
				add:      Vec2{-179, 188},
				subtract: Vec2{-219, -212},
				dot:      -6380,
			},
		},
		{
			"- to -",
			args{
				Vec2{-100, -100},
				Vec2{-200, -200},
			},
			want{
				add:      Vec2{-300, -300},
				subtract: Vec2{100, 100},
				dot:      40000,
			},
		},
		{
			"alternating",
			args{
				Vec2{-100, 100},
				Vec2{200, -200},
			},
			want{
				add:      Vec2{100, -100},
				subtract: Vec2{-300, 300},
				dot:      -40000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want.add) {
				t.Errorf("Add() = %v, want %v", got, tt.want.add)
			}
			if got := Subtract(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want.subtract) {
				t.Errorf("Subtract() = %v, want %v", got, tt.want.subtract)
			}
			if got := Dot(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want.dot) {
				t.Errorf("Dot() = %v, want %v", got, tt.want.dot)
			}
		})
	}
}

func Test_Shrink(t *testing.T) {

	type args struct {
		r Rect
		s int
	}
	type want struct {
		r  Rect
		ok bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"s > r.size / 2",
			args{
				Rect{
					Vec2{-100, -100}, Vec2{201, 201},
				},
				150,
			},
			want{
				Rect{},
				false,
			},
		},
		{
			"s == r.size/2",
			args{
				Rect{
					Vec2{-100, -100}, Vec2{200, 200},
				},
				100,
			},
			want{
				Rect{},
				false,
			},
		},
		{
			"s < r.size/2",
			args{
				Rect{
					Vec2{-100, -100}, Vec2{201, 201},
				},
				100,
			},
			want{
				Rect{
					Vec2{}, Vec2{1, 1},
				},
				true,
			},
		},
		{
			"s == 0",
			args{
				Rect{
					Vec2{-100, -100}, Vec2{201, 201},
				},
				0,
			},
			want{
				Rect{
					Vec2{-100, -100}, Vec2{201, 201},
				},
				true,
			},
		},
		{
			"s < 0",
			args{
				Rect{
					Vec2{-100, -100}, Vec2{201, 201},
				},
				-100,
			},
			want{
				Rect{
					Vec2{-200, -200}, Vec2{401, 401},
				},
				true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, ok := Shrink(tt.args.r, tt.args.s); ok != tt.want.ok || !reflect.DeepEqual(got, tt.want.r) {
				t.Errorf("Contains(%v) = %v, want %v", tt.args, got, tt.want)
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
			if got := Contains(tt.args.r, tt.args.p); got != tt.want {
				t.Errorf("Contains(%v) = %v, want %v", tt.args, got, tt.want)
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
			if got, ok := Overlap(tt.args.r, tt.args.o); ok != tt.want.ok || !reflect.DeepEqual(got, tt.want.inter) {
				t.Errorf("Overlap(%v) = (%s, %v), want %v", tt.args, got, ok, tt.want)
			}
		})
	}
}
