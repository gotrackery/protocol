package generic

import (
	"reflect"
	"testing"

	"github.com/peterstace/simplefeatures/geom"
)

func TestAxisWGS84_Float64(t *testing.T) {
	type fields struct {
		Coordinate float64
		Cardinal   CardinalAxis
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "Zero North",
			fields: fields{
				Coordinate: 0.0,
				Cardinal:   North,
			},
			want: 0.0,
		},
		{
			name: "Zero South",
			fields: fields{
				Coordinate: 0.0,
				Cardinal:   South,
			},
			want: 0.0,
		},
		{
			name: "Zero West",
			fields: fields{
				Coordinate: 0.0,
				Cardinal:   West,
			},
			want: 0.0,
		},
		{
			name: "Zero East",
			fields: fields{
				Coordinate: 0.0,
				Cardinal:   East,
			},
			want: 0.0,
		},
		{
			name: "North",
			fields: fields{
				Coordinate: 5425.98960,
				Cardinal:   North,
			},
			want: 54.43316,
		},
		{
			name: "South",
			fields: fields{
				Coordinate: 5425.98960,
				Cardinal:   South,
			},
			want: -54.43316,
		},
		{
			name: "East",
			fields: fields{
				Coordinate: 5425.98960,
				Cardinal:   East,
			},
			want: 54.43316,
		},
		{
			name: "West",
			fields: fields{
				Coordinate: 5425.98960,
				Cardinal:   West,
			},
			want: -54.43316,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AxisWGS84{
				Coordinate: tt.fields.Coordinate,
				Cardinal:   tt.fields.Cardinal,
			}
			if got := a.Float64(); got != tt.want {
				t.Errorf("Float64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCardinalAxis_Sign(t *testing.T) {
	tests := []struct {
		name string
		c    CardinalAxis
		want float64
	}{
		{
			name: "North",
			c:    North,
			want: 1.0,
		},
		{
			name: "East",
			c:    East,
			want: 1.0,
		},
		{
			name: "South",
			c:    South,
			want: -1.0,
		},
		{
			name: "West",
			c:    West,
			want: -1.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Sign(); got != tt.want {
				t.Errorf("Sign() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseAxisWGS84(t *testing.T) {
	type args struct {
		coord string
		c     CardinalAxis
	}
	tests := []struct {
		name    string
		args    args
		want    AxisWGS84
		wantErr bool
	}{
		{
			name: "Zero North",
			args: args{
				coord: "0.0",
				c:     North,
			},
			want: AxisWGS84{
				Coordinate: 0.0,
				Cardinal:   North,
			},
			wantErr: false,
		},
		{
			name: "North",
			args: args{
				coord: "5425.98960",
				c:     North,
			},
			want: AxisWGS84{
				Coordinate: 5425.98960,
				Cardinal:   North,
			},
			wantErr: false,
		},
		{
			name: "South",
			args: args{
				coord: "5425.98960",
				c:     South,
			},
			want: AxisWGS84{
				Coordinate: 5425.98960,
				Cardinal:   South,
			},
			wantErr: false,
		},
		{
			name: "East",
			args: args{
				coord: "5425.98960",
				c:     East,
			},
			want: AxisWGS84{
				Coordinate: 5425.98960,
				Cardinal:   East,
			},
			wantErr: false,
		},
		{
			name: "West",
			args: args{
				coord: "5425.98960",
				c:     West,
			},
			want: AxisWGS84{
				Coordinate: 5425.98960,
				Cardinal:   West,
			},
			wantErr: false,
		},
		{
			name: "North NA",
			args: args{
				coord: "NA",
				c:     North,
			},
			want: AxisWGS84{
				Coordinate: 0.0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAxisWGS84(tt.args.coord, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAxisWGS84() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseAxisWGS84() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseCardinalAxis(t *testing.T) {
	type args struct {
		raw string
	}
	tests := []struct {
		name    string
		args    args
		wantCp  CardinalAxis
		wantErr bool
	}{
		{
			name:    "North",
			args:    args{raw: "N"},
			wantCp:  North,
			wantErr: false,
		},
		{
			name:    "South",
			args:    args{raw: "S"},
			wantCp:  South,
			wantErr: false,
		},
		{
			name:    "East",
			args:    args{raw: "E"},
			wantCp:  East,
			wantErr: false,
		},
		{
			name:    "West",
			args:    args{raw: "W"},
			wantCp:  West,
			wantErr: false,
		},
		{
			name:    "abracadabra",
			args:    args{raw: "North"},
			wantCp:  "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCp, err := ParseCardinalAxis(tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCardinalAxis() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCp != tt.wantCp {
				t.Errorf("ParseCardinalAxis() gotCp = %v, want %v", gotCp, tt.wantCp)
			}
		})
	}
}

func TestParsePointWGS84(t *testing.T) {
	type args struct {
		lon     string
		cardLon CardinalAxis
		lat     string
		cardLat CardinalAxis
	}
	tests := []struct {
		name    string
		args    args
		want    PointWGS84
		wantErr bool
	}{
		{
			name: "valid NE",
			args: args{
				lon:     "05425.98960",
				cardLon: "N",
				lat:     "05425.98960",
				cardLat: "E",
			},
			want: PointWGS84{
				Lon: AxisWGS84{
					Coordinate: 5425.98960,
					Cardinal:   North,
				},
				Lat: AxisWGS84{
					Coordinate: 5425.98960,
					Cardinal:   East,
				},
				Valid: true,
			},
			wantErr: false,
		},
		{
			name: "valid SW",
			args: args{
				lon:     "05425.98960",
				cardLon: "S",
				lat:     "05425.98960",
				cardLat: "W",
			},
			want: PointWGS84{
				Lon: AxisWGS84{
					Coordinate: 5425.98960,
					Cardinal:   South,
				},
				Lat: AxisWGS84{
					Coordinate: 5425.98960,
					Cardinal:   West,
				},
				Valid: true,
			},
			wantErr: false,
		},
		{
			name: "valid zero NE",
			args: args{
				lon:     "00000.0",
				cardLon: "N",
				lat:     "00000.0",
				cardLat: "E",
			},
			want: PointWGS84{
				Lon: AxisWGS84{
					Coordinate: 0.0,
					Cardinal:   North,
				},
				Lat: AxisWGS84{
					Coordinate: 0.0,
					Cardinal:   East,
				},
				Valid: true,
			},
			wantErr: false,
		},
		{
			name: "valid SW",
			args: args{
				lon:     "00000.0",
				cardLon: "S",
				lat:     "00000.0",
				cardLat: "W",
			},
			want: PointWGS84{
				Lon: AxisWGS84{
					Coordinate: 0.0,
					Cardinal:   South,
				},
				Lat: AxisWGS84{
					Coordinate: 0.0,
					Cardinal:   West,
				},
				Valid: true,
			},
			wantErr: false,
		},
		{
			name: "invalid NA NE",
			args: args{
				lon:     "NA",
				cardLon: "N",
				lat:     "00000.0",
				cardLat: "E",
			},
			want: PointWGS84{
				Lon: AxisWGS84{
					Coordinate: 0.0,
					Cardinal:   "",
				},
				Lat: AxisWGS84{
					Coordinate: 0.0,
					Cardinal:   "",
				},
				Valid: false,
			},
			wantErr: true,
		},
		{
			name: "invalid NA SW",
			args: args{
				lon:     "00000.0",
				cardLon: "S",
				lat:     "NA",
				cardLat: "W",
			},
			want: PointWGS84{
				Lon: AxisWGS84{
					Coordinate: 0.0,
					Cardinal:   "",
				},
				Lat: AxisWGS84{
					Coordinate: 0.0,
					Cardinal:   "",
				},
				Valid: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePointWGS84(tt.args.lon, tt.args.cardLon, tt.args.lat, tt.args.cardLat)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePointWGS84() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePointWGS84() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPointWGS84_LocationXYZ(t *testing.T) {
	type args struct {
		z float64
	}
	type fields struct {
		Lon   AxisWGS84
		Lat   AxisWGS84
		Valid bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Location
	}{
		{
			name: "valid NE",
			fields: fields{
				Lon: AxisWGS84{
					Coordinate: 05425.98960,
					Cardinal:   North,
				},
				Lat: AxisWGS84{
					Coordinate: 05425.98960,
					Cardinal:   East,
				},
				Valid: true,
			},
			args: args{z: 100.0},
			want: Location{
				Coordinates: geom.Coordinates{
					XY: geom.XY{
						X: 54.43316,
						Y: 54.43316,
					},
					Z:    100.0,
					Type: geom.DimXYZ,
				},
				Valid: true,
			},
		},
		{
			name: "valid SW",
			fields: fields{
				Lon: AxisWGS84{
					Coordinate: 05425.98960,
					Cardinal:   South,
				},
				Lat: AxisWGS84{
					Coordinate: 05425.98960,
					Cardinal:   West,
				},
				Valid: true,
			},
			args: args{z: 100.0},
			want: Location{
				Coordinates: geom.Coordinates{
					XY: geom.XY{
						X: -54.43316,
						Y: -54.43316,
					},
					Z:    100.0,
					Type: geom.DimXYZ,
				},
				Valid: true,
			},
		},
		{
			name: "valid zero NE",
			fields: fields{
				Lon: AxisWGS84{
					Coordinate: 0.0,
					Cardinal:   North,
				},
				Lat: AxisWGS84{
					Coordinate: 0.0,
					Cardinal:   East,
				},
				Valid: true,
			},
			args: args{z: 100.0},
			want: Location{
				Coordinates: geom.Coordinates{
					XY: geom.XY{
						X: 0.0,
						Y: 0.0,
					},
					Z:    100.0,
					Type: geom.DimXYZ,
				},
				Valid: true,
			},
		},
		{
			name: "valid zero SW",
			fields: fields{
				Lon: AxisWGS84{
					Coordinate: 0.0,
					Cardinal:   South,
				},
				Lat: AxisWGS84{
					Coordinate: 0.0,
					Cardinal:   West,
				},
				Valid: true,
			},
			args: args{z: 100.0},
			want: Location{
				Coordinates: geom.Coordinates{
					XY: geom.XY{
						X: 0.0,
						Y: 0.0,
					},
					Z:    100.0,
					Type: geom.DimXYZ,
				},
				Valid: true,
			},
		},
		{
			name: "invalid zero NE",
			fields: fields{
				Lon: AxisWGS84{
					Coordinate: 0.0,
					Cardinal:   North,
				},
				Lat: AxisWGS84{
					Coordinate: 0.0,
					Cardinal:   East,
				},
				Valid: false,
			},
			args: args{z: 100.0},
			want: Location{
				Coordinates: geom.Coordinates{
					XY: geom.XY{
						X: 0.0,
						Y: 0.0,
					},
					Z:    100.0,
					Type: geom.DimXYZ,
				},
				Valid: false,
			},
		},
		{
			name: "invalid zero SW",
			fields: fields{
				Lon: AxisWGS84{
					Coordinate: 0.0,
					Cardinal:   South,
				},
				Lat: AxisWGS84{
					Coordinate: 0.0,
					Cardinal:   West,
				},
				Valid: false,
			},
			args: args{z: 100.0},
			want: Location{
				Coordinates: geom.Coordinates{
					XY: geom.XY{
						X: 0.0,
						Y: 0.0,
					},
					Z:    100.0,
					Type: geom.DimXYZ,
				},
				Valid: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := PointWGS84{
				Lon:   tt.fields.Lon,
				Lat:   tt.fields.Lat,
				Valid: tt.fields.Valid,
			}
			if got := w.LocationXYZ(tt.args.z); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LocationXYZ() = %v, want %v", got, tt.want)
			}
		})
	}
}
