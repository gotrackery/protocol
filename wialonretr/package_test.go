package wialonretr

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTest(t *testing.T) {
	// s := "9504000038363637393530333530313631383500643552cb000000070bbb0000001301036164665f7265675f74696d6500643552cc0bbb0000001301036164665f7265675f74696d6500643552cc0bbb0000000d00016d73675f747970650041000bbb00000010000170726f746f00464c4558332e30000bbb0000001500056d73675f6e756d6265720000000000000bdf890bbb0000001100036576656e745f636f646500000017050bbb0000000d000373746174757300000000000bbb0000001100036d6f64756c65735f737400000000a90bbb0000001200036d6f64756c65735f73743200000000800bbb0000000a000367736d000000000e0bbb0000001a00056c6173745f76616c69645f74696d650000000000643552cb0bbb0000001500036e61765f726376725f737461746500000000010bbb00000010000376616c69645f6e617600000000010bbb0000000b00037361747300000000130bbb000000270102706f73696e666f001e48f104d5f94240ef3000673cd04b409a99999999195d4000210158130bbb0000001200046d696c656167650000000060b0acf2400bbb000000180004696e7465725f6d696c6561676500000000a01f54b43f0bbb0000001200047077725f65787400508d976e12833c400bbb0000001200047077725f696e74004e621058393410400bbb0000000f000461646331008b6ce7fba9713c400bbb0000000f00046164633200d7a3703d0ad711400bbb00000012010361766c5f6f75747075747300000000000bbb000000170004656e67696e655f686f757273000d951da6eea0b3400bbb0000000e0003736174735f676c00000000090bbb0000000f0003736174735f677073000000000a0bbb0000000f0003736174735f67616c00000000000bbb000000100003736174735f636f6d7000000000000bbb000000100003736174735f6265696400000000000bbb0000000f0003736174735f646f7200000000000bbb000000110003736174735f69726e737300000000000bbb000000100003736174735f717a737300000000000bbb0000000f000468646f7000000000000000e03f0bbb0000000f000470646f7000000000000000f03f0bbb0000000e000363656c6c5f696400000049be0bbb0000000a00036c61630000006a000bbb0000000a00036d636300000000fa0bbb0000000a00036d6e6300000000630bbb0000000f000372785f6c6576656c00000000b40bbb0000000f000363656c6c5f696431000000cec90bbb0000000b00036c6163310000006a000bbb0000000b00036d63633100000000fa0bbb0000000b00036d6e633100000000630bbb00000010000372785f6c6576656c3100000000a50bbb0000000f000363656c6c5f69643200000031d40bbb0000000b00036c6163320000006a080bbb0000000b00036d63633200000000fa0bbb0000000b00036d6e633200000000630bbb00000010000372785f6c6576656c3200000000a30bbb0000001300056c62735f74696d650000000000643552b00bbb000000140003636f6e6e656374696f6e5f7374000000000a0bbb0000000f0003616363656c5f737400000000140bbb00000011010361766c5f696e7075747300000000017302000038363637393530333439383836373300643552bd000000070bbb0000001301036164665f7265675f74696d6500643552d00bbb0000001301036164665f7265675f74696d6500643552cf0bbb0000000d00016d73675f747970650041000bbb00000010000170726f746f00464c4558332e30000bbb0000001500056d73675f6e756d62657200000000000015acd80bbb0000001100036576656e745f636f6465000000170b0bbb0000000d000373746174757300000000000bbb0000001100036d6f64756c65735f737400000000a90bbb0000001200036d6f64756c65735f73743200000000800bbb0000000a000367736d000000001f0bbb0000001a00056c6173745f76616c69645f74696d650000000000643552bc0bbb0000001500036e61765f726376725f737461746500000000010bbb00000010000376616c69645f6e617600000000010bbb0000000b00037361747300000000140bbb000000270102706f73696e666f001dee10ece5c8424055185b0872ca4b400000000000000000001a0071140bbb0000001200046d696c6561676500000000c0149106410bbb0000001200047077725f65787400b81e85eb51783c400bbb0000001200047077725f696e7400c1caa145b67310400bbb0000000f0004616463310004560e2db25d3c400bbb0000000f0004616463320060e5d022dbf910400bbb00000012010361766c5f6f75747075747300000000000bbb000000170004656e67696e655f686f757273004286ca0ea3c6ba400bbb000000140003636f6e6e656374696f6e5f7374000000000a0bbb0000000f0003616363656c5f737400000000140bbb00000011010361766c5f696e707574730000000001"
	// s := "74000000333533393736303133343435343835004B0BFB70000000030BBB000000270102706F73696E666F00A027AFDF5D9848403AC7253383DD4B400000000000805A40003601460B0BBB0000001200047077725F657874002B8716D9CE973B400BBB00000011010361766C5F696E707574730000000001"
	// s := "7500000038363235333130343537363130373000643552be000000070bbb000000270102706f73696e666f00cc7a319413f14b40d7c05609161b4b403333333333636440005300b2150bbb0000001200047077725f65787400dd24068195433c400bbb0000001200047077725f696e7400448b6ce7fba91040"
	// s := "7500000038363235333130343537363130373000643552d0000000070bbb000000270102706f73696e666f008c101e6d1cf14b40029a081b9e1a4b409a99999999796540005300ad140bbb0000001200047077725f6578740017d9cef753433c400bbb0000001200047077725f696e7400448b6ce7fba91040"
	// s := "7500000038363637393530333831373736313200643552d0000000070bbb000000270102706f73696e666f0086e63a8db4084c401d2098a3c7834b409a99999999a96a4000000103130bbb0000001200047077725f65787400bc74931804563c400bbb0000001200047077725f696e7400c74b378941601040"
	// s := "6c00000033303137383500643552bf000000070bbb000000270102706f73696e666f00840d4faf94fd4b401b12f758fad44a400000000000000000000000aa120bbb0000001200047077725f65787400a01a2fdd24263a400bbb0000001200047077725f696e74008b6ce7fba9711140"
	// s := "e802000038363637393530333832323334353700643552af000000030bbb000000270002706f73696e666f0026c3f17c064444405ddf878384264b400000000000000000004d0118100bbb0000000a01034e756d00002aaae20bbb0000000c01034576656e7400000017040bbb0000000f0103546573744d6f646500000000000bbb000000180103416c61726d4e6f74696669636174696f6e00000000000bbb0000000c0103416c61726d00000000000bbb0000000d010341726d65643100000000000bbb00000011010345766163756174696f6e00000000000bbb000000160103506f776572536176696e674d6f646500000000000bbb0000001e0103416363656c65726f6d6574657243616c6962726174656400000000000bbb0000000a010347534d00000000100bbb0000000c010356616c696400000000010bbb0000001001034c61737456616c696400643552af0bbb0000000e01034d696c6561676500000483cb0bbb0000000e010356736f757263650000006c1d0bbb0000000e01036261747465727900000010090bbb0000000b010341696e3100000000430bbb0000000b010344494e3100000000000bbb0000000b010344494e3200000000000bbb0000000b010344494e3300000000000bbb0000000b010344494e3400000000000bbb0000000b010344494e3500000000000bbb0000000b010344494e3600000000000bbb0000000b010344494e3700000000000bbb0000000b010344494e3800000000000bbb0000000d01036675656c5f31000000082a0bbb0000000d010374656d705f3200000000040bbb0000000f010448444f5000666666666666e63f0bbb0000000f010450444f5000333333333333f33f0bbb000000110103525334383554656d703100000000170bbb00000010010368617273685f61636300000000000bbb00000015010368617273685f627265616b696e6700000000000bbb00000012010368617273685f636f726e6500000000000bbb0000000f010347736d5374617465000000000a"
	s := "e6000000383636313932303332363230393937006436565b000000030bbb000000270102706f73696e666f002ebc6c4e38de42405969520ababb4b400000000000a0654000000100030bbb00000011010361766c5f696e7075747300000000000bbb0000001200047077725f696e74003d0ad7a3703d0c400bbb000000150004656774735f6d745f69640000000000540e48410bbb0000001200047077725f65787400d9cef753e3a5ab3f0bbb0000000f00046c6c7331000000000040ffef400bbb0000000f000467736d31000000000000c055400bbb0000000f00046c6c7332000000000040ffef40"
	bytes, err := hex.DecodeString(s)
	require.NoError(t, err)
	fmt.Println(string(bytes))
	fmt.Println(bytes)
	fmt.Println(int(binary.LittleEndian.Uint32(bytes[:4])) + 4)
	fmt.Println(len(bytes))
}

// TestScanRealPackage is based on recorded real TCP payload with tcpdump and tcpflow.
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