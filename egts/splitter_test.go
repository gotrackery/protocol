package egts

import (
	"bufio"
	"encoding/hex"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitter_Splitter(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantLen   int
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name:      "0001.data",
			args:      args{path: "./testdata/0001.data"},
			wantCount: 41,
			wantLen:   5872,
			wantErr:   assert.NoError,
		},
		{
			name:      "0002.data",
			args:      args{path: "./testdata/0002.data"},
			wantCount: 59,
			wantLen:   6346,
			wantErr:   assert.NoError,
		},
		{
			name:      "0003.data",
			args:      args{path: "./testdata/0003.data"},
			wantCount: 77,
			wantLen:   8673,
			wantErr:   assert.NoError,
		},
		{
			name:      "0004.data",
			args:      args{path: "./testdata/0004.data"},
			wantCount: 395,
			wantLen:   31967,
			wantErr:   assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.args.path)
			require.NoError(t, err)
			defer func() {
				err = f.Close()
				require.NoError(t, err)
			}()

			var (
				cnt,
				ln int
			)

			sp := NewSplitter()
			scanner := bufio.NewScanner(f)
			scanner.Split(sp.Splitter())
			for scanner.Scan() {
				ln += len(scanner.Bytes())
				cnt++
				assert.NoError(t, sp.Error())
				assert.Equal(t, sp.badData, []byte(nil))
				// fmt.Println(hex.EncodeToString(scanner.Bytes()))
			}
			if !tt.wantErr(t, scanner.Err()) {
				t.Errorf("ScanPackage() got unexpected error result = %v", scanner.Err())
			}
			if tt.wantCount != cnt {
				t.Errorf("ScanPackage() got packages = %v, want %v", cnt, tt.wantCount)
			}
			if tt.wantLen != ln {
				t.Errorf("ScanPackage() got length = %v, want %v", ln, tt.wantLen)
			}
		})
	}
}

func TestScanPackage(t *testing.T) {
	bytes, _ := hex.DecodeString("0100010b003a000588013a2f00068801cdd3450202021018006e29ca1880b9b4a1a7f31d3391e600e4f4de588010000000110400081400001804000100000012030000000051e1")
	type args struct {
		data  []byte
		atEOF bool
	}
	tests := []struct {
		name        string
		args        args
		wantAdvance int
		wantToken   []byte
		wantErr     bool
	}{
		{
			name: "empty",
			args: args{
				data:  []byte(""),
				atEOF: false,
			},
			wantAdvance: 0,
			wantToken:   nil,
			wantErr:     false,
		},
		{
			name: "empty eof",
			args: args{
				data:  []byte(""),
				atEOF: true,
			},
			wantAdvance: 0,
			wantToken:   nil,
			wantErr:     false,
		},
		{
			name: "empty eof",
			args: args{
				data:  []byte(""),
				atEOF: true,
			},
			wantAdvance: 0,
			wantToken:   nil,
			wantErr:     false,
		},
		{
			name: "valid data",
			args: args{
				data:  bytes,
				atEOF: false,
			},
			wantAdvance: 71,
			wantToken:   bytes,
			wantErr:     false,
		},
		{
			name: "invalid data at start",
			args: args{
				data:  bytes[1:],
				atEOF: false,
			},
			wantAdvance: 0,
			wantToken:   nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// fmt.Println(hex.EncodeToString(tt.args.data))
			f := NewSplitter().Splitter()
			gotAdvance, gotToken, err := f(tt.args.data, tt.args.atEOF)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanPackage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAdvance != tt.wantAdvance {
				t.Errorf("ScanPackage() gotAdvance = %v, want %v", gotAdvance, tt.wantAdvance)
			}
			if !reflect.DeepEqual(gotToken, tt.wantToken) {
				t.Errorf("ScanPackage() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}
