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
			if got := add(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want.add) {
				t.Errorf("add() = %v, want %v", got, tt.want.add)
			}
			if got := subtract(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want.subtract) {
				t.Errorf("subtract() = %v, want %v", got, tt.want.subtract)
			}
			if got := dot(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want.dot) {
				t.Errorf("dot() = %v, want %v", got, tt.want.dot)
			}
		})
	}
}
