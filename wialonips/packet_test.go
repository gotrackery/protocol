package wialonips

import (
	"encoding/hex"
	"reflect"
	"testing"
	"time"

	"github.com/gotrackery/protocol/common"
	"gopkg.in/guregu/null.v4"
)

func TestDecode(t *testing.T) {
	type args struct {
		data string
		v    Version
		imei string
	}
	tests := []struct {
		name    string
		args    args
		want    Packet
		wantErr bool
	}{
		{
			name: "Invalid Packet - 2 bytes only",
			args: args{
				data: "2358",
				v:    UnknownVersion,
			},
			want:    Packet{},
			wantErr: true,
		},
		{
			name: "Invalid Packet - unknown package type",
			args: args{
				data: "2358233836323436323033313430303536363b4e410d0a",
				v:    UnknownVersion,
			},
			want:    Packet{},
			wantErr: true,
		},
		{
			name: "Invalid Packet - bad structure",
			args: args{
				data: "234c23383632343632303331343030353636",
				v:    UnknownVersion,
			},
			want: Packet{
				Type:    LoginPacket,
				Version: UnknownVersion,
				IMEI:    "",
				Message: nil,
			},
			wantErr: true,
		},
		{
			name: "Login Packet v1.1",
			args: args{
				data: "234c233836323436323033313430303536363b4e410d0a",
				v:    UnknownVersion,
			},
			want: Packet{
				Type:    LoginPacket,
				Version: V1_1,
				IMEI:    "862462031400566",
				Message: &LoginMessage{
					Password: "NA",
					message:  message{imei: "862462031400566", ver: V1_1},
				},
			},
			wantErr: false,
		},
		{
			name: "Login Packet v1.1",
			args: args{
				data: "234c23383632343632303331343030353636",
				v:    UnknownVersion,
			},
			want: Packet{
				Type:    LoginPacket,
				Version: UnknownVersion,
				IMEI:    "",
				Message: nil,
			},
			wantErr: true,
		},
		{
			name: "Login Packet v2.0",
			args: args{
				data: "234c23322e303b3233363031303137333b4e413b3841370d0a",
				v:    UnknownVersion,
			},
			want: Packet{
				Type:    LoginPacket,
				Version: V2_0,
				IMEI:    "236010173",
				Message: &LoginMessage{
					Password: "NA",
					message:  message{imei: "236010173", ver: V2_0},
				},
			},
			wantErr: false,
		},
		{
			name: "ShortData Packet v1.1 all zero",
			args: args{
				data: "235344233234303132333b3133313232323b303030302e30303030303b533b30303030302e30303030303b573b303b303b303b300d0a",
				v:    V1_1,
				imei: "236010173",
			},
			want: Packet{
				Type:    ShortenedDataPacket,
				Version: V1_1,
				IMEI:    "236010173",
				Message: &ShortenedDataMessage{
					RegisteredAt: time.Date(2023, time.January, 24, 13, 12, 22, 0, time.UTC),
					Point: common.PointWGS84{
						Lon:   common.AxisWGS84{Cardinal: common.West},
						Lat:   common.AxisWGS84{Cardinal: common.South},
						Valid: true,
					},
					Speed:    null.NewFloat(0.0, true),
					Course:   null.NewInt(0, true),
					Altitude: null.NewFloat(0.0, true),
					Sat:      null.NewInt(0, true),
					message:  message{imei: "236010173", ver: V1_1},
				},
			},
			wantErr: false,
		},
		{
			name: "ShortData Packet v2.0 all na",
			args: args{
				data: "235344233234303132333b3131303035303b4e413b4e413b4e413b4e413b4e413b4e413b4e413b4e413b4334420d0a",
				v:    V2_0,
			},
			want: Packet{
				Type:    ShortenedDataPacket,
				Version: V2_0,
				Message: &ShortenedDataMessage{
					RegisteredAt: time.Date(2023, time.January, 24, 11, 00, 50, 0, time.UTC),
					Point: common.PointWGS84{
						Lon:   common.AxisWGS84{Cardinal: common.East},
						Lat:   common.AxisWGS84{Cardinal: common.North},
						Valid: false,
					},
					Speed:    null.NewFloat(0.0, false),
					Course:   null.NewInt(0, false),
					Altitude: null.NewFloat(0.0, false),
					Sat:      null.NewInt(0, false),
					message:  message{ver: V2_0},
				},
			},
			wantErr: false,
		},
		{
			name: "Data Packet v1.1 1",
			args: args{
				data: "2344233234303132333b3036303335303b353432352e39383936303b4e3b30343030332e32383134303b453b37392c383b3333333b3133323b31313b4e413b4e413b4e413b3b4e413b",
				v:    V1_1,
			},
			want: Packet{
				Type:    DataPacket,
				Version: V1_1,
				Message: &DataMessage{
					ShortenedDataMessage: ShortenedDataMessage{
						RegisteredAt: time.Date(2023, time.January, 24, 06, 03, 50, 0, time.UTC),
						Point: common.PointWGS84{
							Lon:   common.AxisWGS84{Coordinate: 04003.28140, Cardinal: common.East},
							Lat:   common.AxisWGS84{Coordinate: 5425.98960, Cardinal: common.North},
							Valid: true,
						},
						Speed:    null.NewFloat(79.8, true),
						Course:   null.NewInt(333, true),
						Altitude: null.NewFloat(132.0, true),
						Sat:      null.NewInt(11, true),
						message:  message{ver: V1_1},
					},
					HDOP:       null.NewFloat(0.0, false),
					Inputs:     null.NewInt(0, false),
					Outputs:    null.NewInt(0, false),
					ADC:        nil,
					IButton:    null.NewString("", false),
					Attributes: nil,
				},
			},
			wantErr: false,
		},
		{
			name: "Data Packet v1.1 2",
			args: args{
				data: "2344233234303132333b3036303335303b353432352e39383936303b4e3b30343030332e32383134303b453b37392c383b3333333b3133323b31313b4e413b4e413b4e413b3b4e413b4e6f4461746133323936313a313a313535323037393533330d0a",
				v:    V1_1,
			},
			want: Packet{
				Type:    DataPacket,
				Version: V1_1,
				Message: &DataMessage{
					ShortenedDataMessage: ShortenedDataMessage{
						RegisteredAt: time.Date(2023, time.January, 24, 06, 03, 50, 0, time.UTC),
						Point: common.PointWGS84{
							Lon:   common.AxisWGS84{Coordinate: 04003.28140, Cardinal: common.East},
							Lat:   common.AxisWGS84{Coordinate: 5425.98960, Cardinal: common.North},
							Valid: true,
						},
						Speed:    null.NewFloat(79.8, true),
						Course:   null.NewInt(333, true),
						Altitude: null.NewFloat(132.0, true),
						Sat:      null.NewInt(11, true),
						message:  message{ver: V1_1},
					},
					HDOP:    null.NewFloat(0.0, false),
					Inputs:  null.NewInt(0, false),
					Outputs: null.NewInt(0, false),
					ADC:     nil,
					IButton: null.NewString("", false),
					Attributes: common.Attributes{
						"NoData32961": int64(1552079533),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Data Packet v1.1 3",
			args: args{
				data: "2344233234303132333b3134313434333b353535372e353735373b4e3b30393531322e303135383b453b303b39323b3337303b303b4e413b3235373b4e413b3b4e413b44697374616e63653a323a3139383337372e3230333132352c506f7765723a323a32382e3138343030302c426174746572793a323a342e3033323030302c4d6f746f3a323a3232332e3037323733392c47534d3a313a312c503235303a323a3938333830382e3030303030302c503135313a323a323234332e3030303030302c503135323a323a32303030302e3030303030302c503135333a323a333331312e3030303030302c503138303a323a342e3030303030302c503138313a323a31362e3030303030302c503138323a323a31362e3030303030302c503138333a323a31362e3030303030300d0a",
				v:    V1_1,
			},
			want: Packet{
				Type:    DataPacket,
				Version: V1_1,
				Message: &DataMessage{
					ShortenedDataMessage: ShortenedDataMessage{
						RegisteredAt: time.Date(2023, time.January, 24, 14, 14, 43, 0, time.UTC),
						Point: common.PointWGS84{
							Lon:   common.AxisWGS84{Coordinate: 09512.0158, Cardinal: common.East},
							Lat:   common.AxisWGS84{Coordinate: 5557.5757, Cardinal: common.North},
							Valid: true,
						},
						Speed:    null.NewFloat(0.0, true),
						Course:   null.NewInt(92, true),
						Altitude: null.NewFloat(370.0, true),
						Sat:      null.NewInt(0, true),
						message:  message{ver: V1_1},
					},
					HDOP:    null.NewFloat(0.0, false),
					Inputs:  null.NewInt(257, true),
					Outputs: null.NewInt(0, false),
					ADC:     nil,
					IButton: null.NewString("", false),
					Attributes: common.Attributes{
						"Distance": 198377.203125,
						"Power":    28.184,
						"Battery":  4.032,
						"Moto":     223.072739,
						"GSM":      int64(1),
						"P250":     983808.0,
						"P151":     2243.0,
						"P152":     20000.0,
						"P153":     3311.0,
						"P180":     4.0,
						"P181":     16.0,
						"P182":     16.0,
						"P183":     16.0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Data Packet v1.1 4",
			args: args{
				data: "2344233234303132333b3134303933333b353535302e383531343b4e3b30333734302e373134333b453b303b303b3134362e3030303030303b31363b4e413b4e413b4e413b4e413b4e413b7072696f723a313a302c6576656e745f696f5f69643a313a302c746f74616c5f696f3a313a322c696f5f313a313a302c7077725f6578743a323a32362e3133353030302c696f5f36363a313a3236313335",
				v:    V1_1,
			},
			want: Packet{
				Type:    DataPacket,
				Version: V1_1,
				Message: &DataMessage{
					ShortenedDataMessage: ShortenedDataMessage{
						RegisteredAt: time.Date(2023, time.January, 24, 14, 9, 33, 0, time.UTC),
						Point: common.PointWGS84{
							Lon:   common.AxisWGS84{Coordinate: 03740.7143, Cardinal: common.East},
							Lat:   common.AxisWGS84{Coordinate: 5550.8514, Cardinal: common.North},
							Valid: true,
						},
						Speed:    null.NewFloat(0.0, true),
						Course:   null.NewInt(0, true),
						Altitude: null.NewFloat(146.0, true),
						Sat:      null.NewInt(16, true),
						message:  message{ver: V1_1},
					},
					HDOP:    null.NewFloat(0.0, false),
					Inputs:  null.NewInt(0, false),
					Outputs: null.NewInt(0, false),
					ADC:     []null.Float{null.NewFloat(0.0, false)},
					IButton: null.NewString("", false),
					Attributes: common.Attributes{
						"prior":       int64(0),
						"event_io_id": int64(0),
						"total_io":    int64(2),
						"io_1":        int64(0),
						"pwr_ext":     26.135000,
						"io_66":       int64(26135),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Data Packet v1.1 5",
			args: args{
				data: "2344233234303132333b3134313434333b353533332e383631323b4e3b30333734352e303636353b453b303b3139323b302e3030303030303b3235353b4e413b313b303b31312e3030303030302c31302e3030303030303b4e413b76616c69645f636f6f7264733a313a312c70726f746f3a333a756e6b6e6f776e2c6e617669676174696f6e5f73797374656d3a333a4750532b476c6f6e617373",
				v:    V1_1,
			},
			want: Packet{
				Type:    DataPacket,
				Version: V1_1,
				Message: &DataMessage{
					ShortenedDataMessage: ShortenedDataMessage{
						RegisteredAt: time.Date(2023, time.January, 24, 14, 14, 43, 0, time.UTC),
						Point: common.PointWGS84{
							Lon:   common.AxisWGS84{Coordinate: 03745.0665, Cardinal: common.East},
							Lat:   common.AxisWGS84{Coordinate: 5533.8612, Cardinal: common.North},
							Valid: true,
						},
						Speed:    null.NewFloat(0.0, true),
						Course:   null.NewInt(192, true),
						Altitude: null.NewFloat(0.0, true),
						Sat:      null.NewInt(255, true),
						message:  message{ver: V1_1},
					},
					HDOP:    null.NewFloat(0.0, false),
					Inputs:  null.NewInt(1, true),
					Outputs: null.NewInt(0, true),
					ADC:     []null.Float{null.NewFloat(11.000000, true), null.NewFloat(10.000000, true)},
					IButton: null.NewString("", false),
					Attributes: common.Attributes{
						"valid_coords":      int64(1),
						"proto":             "unknown",
						"navigation_system": "GPS+Glonass",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Data Packet v1.1 6",
			args: args{
				data: "2344233234303132333b3134313430393b353531352e393539353b4e3b30373635352e353637313b453b303b37353b3131352e3030303030303b31353b302e3630303030303b4e413b3730343634333038373b4e412c372e3833333030302c4e412c4e412c4e412c4e412c4e412c4e412c4e412c4e412c4e412c373833332e3030303030303b4e413b67736d5f7374617475733a313a322c6163635f747269676765723a313a302c6465765f7374617475733a313a3130373532",
				v:    V1_1,
			},
			want: Packet{
				Type:    DataPacket,
				Version: V1_1,
				Message: &DataMessage{
					ShortenedDataMessage: ShortenedDataMessage{
						RegisteredAt: time.Date(2023, time.January, 24, 14, 14, 9, 0, time.UTC),
						Point: common.PointWGS84{
							Lon:   common.AxisWGS84{Coordinate: 07655.5671, Cardinal: common.East},
							Lat:   common.AxisWGS84{Coordinate: 5515.9595, Cardinal: common.North},
							Valid: true,
						},
						Speed:    null.NewFloat(0.0, true),
						Course:   null.NewInt(75, true),
						Altitude: null.NewFloat(115.0, true),
						Sat:      null.NewInt(15, true),
						message:  message{ver: V1_1},
					},
					HDOP:    null.NewFloat(0.6, true),
					Inputs:  null.NewInt(0, false),
					Outputs: null.NewInt(704643087, true),
					ADC: []null.Float{
						null.NewFloat(0.0, false),
						null.NewFloat(7.8330, true),
						null.NewFloat(0.0, false),
						null.NewFloat(0.0, false),
						null.NewFloat(0.0, false),
						null.NewFloat(0.0, false),
						null.NewFloat(0.0, false),
						null.NewFloat(0.0, false),
						null.NewFloat(0.0, false),
						null.NewFloat(0.0, false),
						null.NewFloat(0.0, false),
						null.NewFloat(7833.0, true),
					},
					IButton: null.NewString("", false),
					Attributes: common.Attributes{
						"gsm_status":  int64(2),
						"acc_trigger": int64(0),
						"dev_status":  int64(10752),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Data Packet v2.0 1",
			args: args{
				data: "2344233234303132333b3133353332313b353734382e333332393b4e3b30343331322e393033383b453b38393b34323b313235303b31363b4e413b4e413b4e413b3b4e413b69676e3a313a312c6675656c313a313a323138332c6675656c323a313a323432313b344239420d0a",
				v:    V2_0,
			},
			want: Packet{
				Type:    DataPacket,
				Version: V2_0,
				Message: &DataMessage{
					ShortenedDataMessage: ShortenedDataMessage{
						RegisteredAt: time.Date(2023, time.January, 24, 13, 53, 21, 0, time.UTC),
						Point: common.PointWGS84{
							Lon:   common.AxisWGS84{Coordinate: 04312.9038, Cardinal: common.East},
							Lat:   common.AxisWGS84{Coordinate: 5748.3329, Cardinal: common.North},
							Valid: true,
						},
						Speed:    null.NewFloat(89.0, true),
						Course:   null.NewInt(42, true),
						Altitude: null.NewFloat(1250.0, true),
						Sat:      null.NewInt(16, true),
						message:  message{ver: V2_0},
					},
					HDOP:    null.NewFloat(0.0, false),
					Inputs:  null.NewInt(0, false),
					Outputs: null.NewInt(0, false),
					ADC:     nil,
					IButton: null.NewString("", false),
					Attributes: common.Attributes{
						"ign":   int64(1),
						"fuel1": int64(2183),
						"fuel2": int64(2421),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Data Packet v2.0 2",
			args: args{
				data: "2344233234303132333b3134303231343b353735302e313234393b4e3b30343135322e323736303b453b37393b3131313b313634373b31393b4e413b4e413b4e413b3b4e413b69676e3a313a312c6163635f783a313a302c6163635f793a313a302c6163635f7a3a313a2d313b384342310d0a",
				v:    V2_0,
			},
			want: Packet{
				Type:    DataPacket,
				Version: V2_0,
				Message: &DataMessage{
					ShortenedDataMessage: ShortenedDataMessage{
						RegisteredAt: time.Date(2023, time.January, 24, 14, 2, 14, 0, time.UTC),
						Point: common.PointWGS84{
							Lon:   common.AxisWGS84{Coordinate: 04152.2760, Cardinal: common.East},
							Lat:   common.AxisWGS84{Coordinate: 5750.1249, Cardinal: common.North},
							Valid: true,
						},
						Speed:    null.NewFloat(79.0, true),
						Course:   null.NewInt(111, true),
						Altitude: null.NewFloat(1647.0, true),
						Sat:      null.NewInt(19, true),
						message:  message{ver: V2_0},
					},
					HDOP:    null.NewFloat(0.0, false),
					Inputs:  null.NewInt(0, false),
					Outputs: null.NewInt(0, false),
					ADC:     nil,
					IButton: null.NewString("", false),
					Attributes: common.Attributes{
						"ign":   int64(1),
						"acc_x": int64(0),
						"acc_y": int64(0),
						"acc_z": int64(-1),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Data Packet v1.1 na",
			args: args{
				data: "2344233230303132333b3036323831313b4e413b4e413b4e413b4e413b4e413b4e413b4e413b4e413b4e413b303b303b4e413b4e413b67736d5f73743a313a332c6e61765f73743a313a302c6d773a313a30",
				v:    V1_1,
			},
			want: Packet{
				Type:    DataPacket,
				Version: V1_1,
				Message: &DataMessage{
					ShortenedDataMessage: ShortenedDataMessage{
						RegisteredAt: time.Date(2023, time.January, 20, 6, 28, 11, 0, time.UTC),
						Point: common.PointWGS84{
							Lon:   common.AxisWGS84{Cardinal: common.East},
							Lat:   common.AxisWGS84{Cardinal: common.North},
							Valid: false,
						},
						Speed:    null.NewFloat(0.0, false),
						Course:   null.NewInt(0, false),
						Altitude: null.NewFloat(0.0, false),
						Sat:      null.NewInt(0, false),
						message:  message{ver: V1_1},
					},
					HDOP:    null.NewFloat(0.0, false),
					Inputs:  null.NewInt(0, true),
					Outputs: null.NewInt(0, true),
					ADC:     []null.Float{null.NewFloat(0.0, false)},
					IButton: null.NewString("", false),
					Attributes: common.Attributes{
						"gsm_st": int64(3),
						"nav_st": int64(0),
						"mw":     int64(0),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "BlackBox Packet v1.1",
			args: args{
				data: "2342233234303132333b3134313833383b353534372e373835303b4e3b30333733342e373734303b453b303b3138343b3138383b31323b313b3b303b3b4e413b6f646f3a323a302c63656c6c5f69643a313a35303137362c6c61633a313a3339332c6d63633a313a3235302c6d6e633a313a312c67736d5f6c766c3a313a36302c74613a313a302c66726d5f76657273696f6e3a333a51434f517c3234303132333b3134313835383b303030302e303030303b4e3b30303030302e303030303b453b303b303b303b303b303b3b303b3b4e413b6f646f3a323a302c63656c6c5f69643a313a35303137362c6c61633a313a3339332c6d63633a313a3235302c6d6e633a313a312c67736d5f6c766c3a313a36312c74613a313a302c66726d5f76657273696f6e3a333a51434f51",
				v:    V1_1,
			},
			want: Packet{
				Type:    BlackBoxPacket,
				Version: V1_1,
				Message: &BlackBoxMessage{
					ShortenedMessages: nil,
					DataMessages: []DataMessage{
						{ShortenedDataMessage: ShortenedDataMessage{
							RegisteredAt: time.Date(2023, time.January, 24, 14, 18, 38, 0, time.UTC),
							Point: common.PointWGS84{
								Lon:   common.AxisWGS84{Coordinate: 03734.7740, Cardinal: common.East},
								Lat:   common.AxisWGS84{Coordinate: 5547.7850, Cardinal: common.North},
								Valid: true,
							},
							Speed:    null.NewFloat(0.0, true),
							Course:   null.NewInt(184, true),
							Altitude: null.NewFloat(188.0, true),
							Sat:      null.NewInt(12, true),
							message:  message{},
						},
							HDOP:    null.NewFloat(1.0, true),
							Inputs:  null.NewInt(0, false),
							Outputs: null.NewInt(0, true),
							ADC:     nil,
							IButton: null.NewString("", false),
							Attributes: common.Attributes{
								"odo":         0.0,
								"cell_id":     int64(50176),
								"lac":         int64(393),
								"mcc":         int64(250),
								"mnc":         int64(1),
								"gsm_lvl":     int64(60),
								"ta":          int64(0),
								"frm_version": "QCOQ",
							},
						},
						{ShortenedDataMessage: ShortenedDataMessage{
							RegisteredAt: time.Date(2023, time.January, 24, 14, 18, 58, 0, time.UTC),
							Point: common.PointWGS84{
								Lon:   common.AxisWGS84{Coordinate: 0.0, Cardinal: common.East},
								Lat:   common.AxisWGS84{Coordinate: 0.0, Cardinal: common.North},
								Valid: true,
							},
							Speed:    null.NewFloat(0.0, true),
							Course:   null.NewInt(0, true),
							Altitude: null.NewFloat(0.0, true),
							Sat:      null.NewInt(0, true),
							message:  message{},
						},
							HDOP:    null.NewFloat(0.0, true),
							Inputs:  null.NewInt(0, false),
							Outputs: null.NewInt(0, true),
							ADC:     nil,
							IButton: null.NewString("", false),
							Attributes: common.Attributes{
								"odo":         0.0,
								"cell_id":     int64(50176),
								"lac":         int64(393),
								"mcc":         int64(250),
								"mnc":         int64(1),
								"gsm_lvl":     int64(61),
								"ta":          int64(0),
								"frm_version": "QCOQ",
							},
						},
					},
					message: message{ver: V1_1},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := hex.DecodeString(tt.args.data)
			if err != nil {
				panic(err)
			}

			got := NewPacket(tt.args.v, tt.args.imei)
			err = got.Decode(data)
			// fmt.Println(tt.want.Message)
			// fmt.Println(got.Message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_String(t *testing.T) {
	tests := []struct {
		name  string
		s     Version
		want  string
		valid bool
	}{
		{
			name:  "unknown",
			s:     UnknownVersion,
			want:  "Version(0)",
			valid: false,
		},
		{
			name:  "V1_1",
			s:     V1_1,
			want:  "1.1",
			valid: true,
		},
		{
			name:  "V2_0",
			s:     V2_0,
			want:  "2.0",
			valid: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
			if got := tt.s.IsValid(); got != tt.valid {
				t.Errorf("IsValid() = %v, valid %v", got, tt.valid)
			}
		})
	}
}
