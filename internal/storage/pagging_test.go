package storage

import (
	"reflect"
	"testing"
)

func TestPagging(t *testing.T) {
	type args struct {
		page     int
		pageSize int
	}
	tests := map[string]struct {
		args args
		want Page
	}{
		"1": {
			args: args{
				page:     0,
				pageSize: 0,
			},
			want: Page{
				Limit:  DefaultLimit,
				Offset: DefaultOffset,
			},
		},
		"2": {
			args: args{
				page:     100,
				pageSize: 0,
			},
			want: Page{
				Limit:  DefaultLimit,
				Offset: (100 - 1) * DefaultLimit,
			},
		},
		"3": {
			args: args{
				page:     1,
				pageSize: 3000,
			},
			want: Page{
				Limit:  MaxLimit,
				Offset: 0,
			},
		},
		"4": {
			args: args{
				page:     1,
				pageSize: 20,
			},
			want: Page{
				Limit:  20,
				Offset: 0,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := Pagging(&tt.args.page, &tt.args.pageSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pagging() = %v, want %v", got, tt.want)
			}
		})
	}
}
