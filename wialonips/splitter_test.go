package wialonips

import (
	"bufio"
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
			wantCount: 9,
			wantLen:   5536,
			wantErr:   assert.NoError,
		},
		{
			name:      "0002.data",
			args:      args{path: "./testdata/0002.data"},
			wantCount: 11,
			wantLen:   1113,
			wantErr:   assert.NoError,
		},
		{
			name:      "0003.data",
			args:      args{path: "./testdata/0003.data"},
			wantCount: 28,
			wantLen:   10375,
			wantErr:   assert.NoError,
		},
		{
			name:      "0004.data",
			args:      args{path: "./testdata/0004.data"},
			wantCount: 12,
			wantLen:   9793,
			wantErr:   assert.NoError,
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
			name: "empty crlf",
			args: args{
				data:  []byte("\r\n"),
				atEOF: false,
			},
			wantAdvance: 0,
			wantToken:   nil,
			wantErr:     true,
		},
		{
			name: "empty crlf second",
			args: args{
				data:  []byte("\r\n_"),
				atEOF: false,
			},
			wantAdvance: 0,
			wantToken:   nil,
			wantErr:     true,
		},
		{
			name: "one line cut",
			args: args{
				data:  []byte("1test\r\n_"),
				atEOF: false,
			},
			wantAdvance: 7,
			wantToken:   []byte("1test\r\n"),
			wantErr:     false,
		},
		{
			name: "one line",
			args: args{
				data:  []byte("2test\r\n"),
				atEOF: false,
			},
			wantAdvance: 7,
			wantToken:   []byte("2test\r\n"),
			wantErr:     false,
		},
		{
			name: "one line eof",
			args: args{
				data:  []byte("#test\r\n"),
				atEOF: true,
			},
			wantAdvance: 7,
			wantToken:   []byte("#test\r\n"),
			wantErr:     false,
		},
		{
			name: "eof",
			args: args{
				data:  []byte("1test"),
				atEOF: true,
			},
			wantAdvance: 0,
			wantToken:   nil,
			wantErr:     true,
		},
		{
			name: "not a new line",
			args: args{
				data:  []byte("#test\n"),
				atEOF: false,
			},
			wantAdvance: 0,
			wantToken:   nil,
			wantErr:     false,
		},
		{
			name: "not a new line eof",
			args: args{
				data:  []byte("#test\n"),
				atEOF: true,
			},
			wantAdvance: 0,
			wantToken:   nil,
			wantErr:     true,
		},
		{
			name: "alone break",
			args: args{
				data:  []byte("\n"),
				atEOF: false,
			},
			wantAdvance: 0,
			wantToken:   nil,
			wantErr:     true,
		},
		{
			name: "alone break eof",
			args: args{
				data:  []byte("\n"),
				atEOF: true,
			},
			wantAdvance: 0,
			wantToken:   nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
