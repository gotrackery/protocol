package wialonretr

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScanRealPackage(t *testing.T) {
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
			args:      args{path: "../test/wialonretr/0001.data"},
			wantCount: 1,
			wantLen:   1005,
			wantErr:   assert.NoError,
		},
		{
			name:      "0002.data",
			args:      args{path: "../test/wialonretr/0002.data"},
			wantCount: 1,
			wantLen:   234,
			wantErr:   assert.NoError,
		},
		{
			name:      "0003.data",
			args:      args{path: "../test/wialonretr/0003.data"},
			wantCount: 9080,
			wantLen:   1081778,
			wantErr:   assert.NoError,
		},
		{
			name:      "0004.data",
			args:      args{path: "../test/wialonretr/0004.data"},
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
			scanner := bufio.NewScanner(f)
			scanner.Split(ScanPackage)
			for scanner.Scan() {
				ln += len(scanner.Bytes())
				cnt++
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

func Test_scanBlock(t *testing.T) {
	type args struct {
		data  string
		atEOF bool
	}
	tests := []struct {
		name        string
		args        args
		wantAdvance int
		wantToken   []byte
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "empty",
			args: args{
				data:  "",
				atEOF: false,
			},
			wantAdvance: 0,
			wantToken:   nil,
			wantErr:     assert.NoError,
		},
		{
			name: "leading 1",
			args: args{
				data:  "310bbb320bbb34",
				atEOF: false,
			},
			wantAdvance: 3,
			wantToken:   []byte{0x31},
			wantErr:     assert.NoError,
		},
		{
			name: "leading 1 eof",
			args: args{
				data:  "310bbb320bbb34",
				atEOF: true,
			},
			wantAdvance: 3,
			wantToken:   []byte{0x31},
			wantErr:     assert.NoError,
		},
		{
			name: "leading empty",
			args: args{
				data:  "0bbb320bbb34",
				atEOF: false,
			},
			wantAdvance: 2,
			wantToken:   nil,
			wantErr:     assert.NoError,
		},
		{
			name: "leading empty eof",
			args: args{
				data:  "0bbb320bbb34",
				atEOF: false,
			},
			wantAdvance: 2,
			wantToken:   nil,
			wantErr:     assert.NoError,
		},
		{
			name: "not all data",
			args: args{
				data:  "323334",
				atEOF: false,
			},
			wantAdvance: 0,
			wantToken:   nil,
			wantErr:     assert.NoError,
		},
		{
			name: "rest data",
			args: args{
				data:  "323334",
				atEOF: true,
			},
			wantAdvance: 3,
			wantToken:   []byte{0x32, 0x33, 0x34},
			wantErr:     assert.NoError,
		},
		{
			name: "alone block sep",
			args: args{
				data:  "0bbb",
				atEOF: false,
			},
			wantAdvance: 2,
			wantToken:   nil,
			wantErr:     assert.NoError,
		},
		{
			name: "alone block sep eof",
			args: args{
				data:  "0bbb",
				atEOF: true,
			},
			wantAdvance: 2,
			wantToken:   nil,
			wantErr:     assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := hex.DecodeString(tt.args.data)
			require.NoError(t, err)
			gotAdvance, gotToken, err := scanBlock(b, tt.args.atEOF)
			if !tt.wantErr(t, err, fmt.Sprintf("scanBlock(%v, %v)", tt.args.data, tt.args.atEOF)) {
				return
			}
			assert.Equalf(t, tt.wantAdvance, gotAdvance, "scanBlock(%v, %v)", tt.args.data, tt.args.atEOF)
			assert.Equalf(t, tt.wantToken, gotToken, "scanBlock(%v, %v)", tt.args.data, tt.args.atEOF)
		})
	}
}

func Test_scanBlockRealData(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name:      "first",
			args:      args{data: "0bbb0000001301036164665f7265675f74696d6500643552cc0bbb0000001301036164665f7265675f74696d6500643552cc0bbb0000000d00016d73675f747970650041000bbb00000010000170726f746f00464c4558332e30000bbb0000001500056d73675f6e756d6265720000000000000bdf890bbb0000001100036576656e745f636f646500000017050bbb0000000d000373746174757300000000000bbb0000001100036d6f64756c65735f737400000000a90bbb0000001200036d6f64756c65735f73743200000000800bbb0000000a000367736d000000000e0bbb0000001a00056c6173745f76616c69645f74696d650000000000643552cb0bbb0000001500036e61765f726376725f737461746500000000010bbb00000010000376616c69645f6e617600000000010bbb0000000b00037361747300000000130bbb000000270102706f73696e666f001e48f104d5f94240ef3000673cd04b409a99999999195d4000210158130bbb0000001200046d696c656167650000000060b0acf2400bbb000000180004696e7465725f6d696c6561676500000000a01f54b43f0bbb0000001200047077725f65787400508d976e12833c400bbb0000001200047077725f696e74004e621058393410400bbb0000000f000461646331008b6ce7fba9713c400bbb0000000f00046164633200d7a3703d0ad711400bbb00000012010361766c5f6f75747075747300000000000bbb000000170004656e67696e655f686f757273000d951da6eea0b3400bbb0000000e0003736174735f676c00000000090bbb0000000f0003736174735f677073000000000a0bbb0000000f0003736174735f67616c00000000000bbb000000100003736174735f636f6d7000000000000bbb000000100003736174735f6265696400000000000bbb0000000f0003736174735f646f7200000000000bbb000000110003736174735f69726e737300000000000bbb000000100003736174735f717a737300000000000bbb0000000f000468646f7000000000000000e03f0bbb0000000f000470646f7000000000000000f03f0bbb0000000e000363656c6c5f696400000049be0bbb0000000a00036c61630000006a000bbb0000000a00036d636300000000fa0bbb0000000a00036d6e6300000000630bbb0000000f000372785f6c6576656c00000000b40bbb0000000f000363656c6c5f696431000000cec90bbb0000000b00036c6163310000006a000bbb0000000b00036d63633100000000fa0bbb0000000b00036d6e633100000000630bbb00000010000372785f6c6576656c3100000000a50bbb0000000f000363656c6c5f69643200000031d40bbb0000000b00036c6163320000006a080bbb0000000b00036d63633200000000fa0bbb0000000b00036d6e633200000000630bbb00000010000372785f6c6576656c3200000000a30bbb0000001300056c62735f74696d650000000000643552b00bbb000000140003636f6e6e656374696f6e5f7374000000000a0bbb0000000f0003616363656c5f737400000000140bbb00000011010361766c5f696e707574730000000001"},
			wantCount: 52,
			wantErr:   assert.NoError,
		},
		{
			name:      "second",
			args:      args{data: "0bbb000000270102706f73696e666f00840d4faf94fd4b401b12f758fad44a400000000000000000000000aa120bbb0000001200047077725f65787400a01a2fdd24263a400bbb0000001200047077725f696e74008b6ce7fba9711140"},
			wantCount: 3,
			wantErr:   assert.NoError,
		},
		{
			name:      "third",
			args:      args{data: "0bbb000000270102706f73696e666f0086e63a8db4084c401d2098a3c7834b409a99999999a96a4000000103130bbb0000001200047077725f65787400bc74931804563c400bbb0000001200047077725f696e7400c74b378941601040"},
			wantCount: 3,
			wantErr:   assert.NoError,
		},
		{
			name:      "fourth",
			args:      args{data: "0bbb000000270102706f73696e666f002ebc6c4e38de42405969520ababb4b400000000000a0654000000100030bbb00000011010361766c5f696e7075747300000000000bbb0000001200047077725f696e74003d0ad7a3703d0c400bbb000000150004656774735f6d745f69640000000000540e48410bbb0000001200047077725f65787400d9cef753e3a5ab3f0bbb0000000f00046c6c7331000000000040ffef400bbb0000000f000467736d31000000000000c055400bbb0000000f00046c6c7332000000000040ffef40"},
			wantCount: 8,
			wantErr:   assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := hex.DecodeString(tt.args.data)
			require.NoError(t, err)

			var (
				cnt,
				ln int
			)
			scanner := bufio.NewScanner(bytes.NewReader(b))
			scanner.Split(scanBlock)
			for scanner.Scan() {
				ln += len(scanner.Bytes())
				ln += 2
				cnt++
				// fmt.Println(hex.EncodeToString(scanner.Bytes()))
				// fmt.Println(scanner.Text())
			}
			if !tt.wantErr(t, scanner.Err()) {
				t.Errorf("ScanPackage() got unexpected error result = %v", scanner.Err())
			}
			if tt.wantCount != cnt {
				t.Errorf("ScanPackage() got blocks = %v, want %v", cnt, tt.wantCount)
			}
			if len(tt.args.data) != ln*2 {
				t.Errorf("ScanPackage() got length = %v, want %v", ln*2, len(tt.args.data))
			}
		})
	}
}

func TestPacket_Decode(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name                  string
		args                  args
		want                  Packet
		wantHasLocation       bool
		wantHasDigitalInputs  bool
		wantHasDigitalOutputs bool
		wantHasAlerts         bool
		wantHasDriverID       bool
		wantErr               bool
	}{
		{
			name: "spec example",
			args: args{data: "74000000333533393736303133343435343835005D515DBB000000030BBB000000270102706F73696E666F00A027AFDF5D9848403AC7253383DD4B400000000000805A40003601460B0BBB0000001200047077725F657874002B8716D9CE973B400BBB00000011010361766C5F696E707574730000000001"},
			want: Packet{
				DeviceID:     "353976013445485",
				RegisteredAt: time.Unix(1565613499, 0),
				DataBlocks: map[string]DataBlock{
					PosInfoName: {
						securityParam: hiddenParam,
						name:          PosInfoName,
						Value: PositionInfo{
							Lon:    49.1903648,
							Lat:    55.7305664,
							Alt:    106.0,
							Speed:  54,
							Course: 326,
							Sats:   11,
						},
					},
					"pwr_ext": {
						securityParam: shownParam,
						name:          "pwr_ext",
						Value:         27.593,
					},
					AvlInputsName: {
						securityParam: hiddenParam,
						name:          AvlInputsName,
						Value:         int32(1),
					},
				},
				bitFlags: 3,
				err:      nil,
			},
			wantHasLocation:       true,
			wantHasDigitalInputs:  true,
			wantHasDigitalOutputs: false,
			wantHasAlerts:         false,
			wantHasDriverID:       false,
			wantErr:               false,
		},
		{
			name: "another example",
			args: args{data: "6c00000033303137383500643552bf000000070bbb000000270102706f73696e666f00840d4faf94fd4b401b12f758fad44a400000000000000000000000aa120bbb0000001200047077725f65787400a01a2fdd24263a400bbb0000001200047077725f696e74008b6ce7fba9711140"},
			want: Packet{
				DeviceID:     "301785",
				RegisteredAt: time.Unix(1681216191, 0),
				DataBlocks: map[string]DataBlock{
					PosInfoName: {
						securityParam: hiddenParam,
						name:          PosInfoName,
						Value: PositionInfo{
							Lon:    55.9811,
							Lat:    53.66389,
							Alt:    0.0,
							Speed:  0,
							Course: 170,
							Sats:   18,
						},
					},
					"pwr_ext": {
						securityParam: shownParam,
						name:          "pwr_ext",
						Value:         26.149,
					},
					"pwr_int": {
						securityParam: shownParam,
						name:          "pwr_int",
						Value:         4.361,
					},
				},
				bitFlags: 7,
				err:      nil,
			},
			wantHasLocation:       true,
			wantHasDigitalInputs:  true,
			wantHasDigitalOutputs: true,
			wantHasAlerts:         false,
			wantHasDriverID:       false,
			wantErr:               false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := hex.DecodeString(tt.args.data)
			require.NoError(t, err)

			var got Packet
			err = got.Decode(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// fmt.Println(got)
			// fmt.Println(tt.want)

			assert.Equalf(t, tt.wantHasLocation, got.HasLocation(),
				"Decode() got.HasLocation() = %v, want %v", got.HasLocation(), tt.wantHasLocation)
			assert.Equalf(t, tt.wantHasDigitalInputs, got.HasDigitalInputs(),
				"Decode() got.HasDigitalInputs() = %v, want %v", got.HasDigitalInputs(), tt.wantHasDigitalInputs)
			assert.Equalf(t, tt.wantHasDigitalOutputs, got.HasDigitalOutputs(),
				"Decode() got.HasDigitalOutputs() = %v, want %v", got.HasDigitalOutputs(), tt.wantHasDigitalOutputs)
			assert.Equalf(t, tt.wantHasAlerts, got.HasAlerts(),
				"Decode() got.HasAlerts() = %v, want %v", got.HasAlerts(), tt.wantHasAlerts)
			assert.Equalf(t, tt.wantHasDriverID, got.HasDriverID(),
				"Decode() got.HasDriverID() = %v, want %v", got.HasDriverID(), tt.wantHasDriverID)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
