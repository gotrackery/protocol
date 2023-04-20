package wialonretr

import (
	"bufio"
	"os"
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
			wantCount: 1,
			wantLen:   1005,
			wantErr:   assert.NoError,
		},
		{
			name:      "0002.data",
			args:      args{path: "./testdata/0002.data"},
			wantCount: 1,
			wantLen:   234,
			wantErr:   assert.NoError,
		},
		{
			name:      "0003.data",
			args:      args{path: "./testdata/0003.data"},
			wantCount: 9080,
			wantLen:   1081778,
			wantErr:   assert.NoError,
		},
		{
			name:      "0004.data",
			args:      args{path: "./testdata/0004.data"},
			wantCount: 202,
			wantLen:   114162,
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
