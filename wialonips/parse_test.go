package wialonips

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseValue(t *testing.T) {
	type args struct {
		b    byte
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantVal interface{}
		wantErr bool
	}{
		{
			name: "int",
			args: args{
				b:    intType,
				data: []byte("4444"),
			},
			wantVal: int64(4444),
			wantErr: false,
		},
		{
			name: "float",
			args: args{
				b:    floatType,
				data: []byte("44.44"),
			},
			wantVal: 44.44,
			wantErr: false,
		},
		{
			name: "string",
			args: args{
				b:    stringType,
				data: []byte("44.44"),
			},
			wantVal: "44.44",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, err := parseValue(tt.args.b, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseValue() error = %ver, wantErr %ver", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVal, tt.wantVal) {
				t.Errorf("parseValue() gotVal = %ver, want %ver", gotVal, tt.wantVal)
			}
		})
	}
}

func Test_parseCRC(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantCrc uint16
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "valid 1",
			args:    args{data: []byte("8A7")},
			wantCrc: 2215,
			wantErr: assert.NoError,
		},
		{
			name:    "valid 2",
			args:    args{data: []byte("08A7")},
			wantCrc: 2215,
			wantErr: assert.NoError,
		},
		{
			name:    "invalid 1",
			args:    args{data: []byte("A7")},
			wantCrc: 0,
			wantErr: assert.Error,
		},
		{
			name:    "invalid 2",
			args:    args{data: []byte("08A71")},
			wantCrc: 0,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCrc, err := parseCRC(tt.args.data)
			if !tt.wantErr(t, err, fmt.Sprintf("parseCRC(%v)", tt.args.data)) {
				return
			}
			assert.Equalf(t, tt.wantCrc, gotCrc, "parseCRC(%v)", tt.args.data)
		})
	}
}
