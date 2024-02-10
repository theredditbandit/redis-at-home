package main

import (
	"reflect"
	"sync"
	"testing"
)

func Test_respHandler(t *testing.T) {
	type args struct {
		data    []byte
		dstream datastream
		wg      *sync.WaitGroup
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := respHandler(tt.args.data, tt.args.dstream, tt.args.wg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("respHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
