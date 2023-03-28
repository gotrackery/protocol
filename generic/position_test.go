package generic

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/gotrackery/protocol"
	"gopkg.in/guregu/null.v4"
)

func TestAttributes_AppendNullInt(t *testing.T) {
	type args struct {
		name string
		i    null.Int
	}
	tests := []struct {
		name string
		a    protocol.Attributes
		args args
		want protocol.Attributes
	}{
		{
			name: "add to nil",
			a:    nil,
			args: args{
				name: "test",
				i: null.Int{NullInt64: sql.NullInt64{
					Int64: 4,
					Valid: true,
				}},
			},
			want: protocol.Attributes{"test": int64(4)},
		},
		{
			name: "add to not nil",
			a:    protocol.Attributes{"early": "test"},
			args: args{
				name: "test",
				i: null.Int{NullInt64: sql.NullInt64{
					Int64: 4,
					Valid: true,
				}},
			},
			want: protocol.Attributes{"early": "test", "test": int64(4)},
		},
		{
			name: "invalid add to not nil",
			a:    protocol.Attributes{"early": "test"},
			args: args{
				name: "test",
				i: null.Int{NullInt64: sql.NullInt64{
					Int64: 4,
					Valid: false,
				}},
			},
			want: protocol.Attributes{"early": "test"},
		},
		{
			name: "invalid add to not nil",
			a:    nil,
			args: args{
				name: "test",
				i: null.Int{NullInt64: sql.NullInt64{
					Int64: 4,
					Valid: false,
				}},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.AppendNullInt(tt.args.name, tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppendNullInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributes_AppendNullFloat(t *testing.T) {
	type args struct {
		name string
		f    null.Float
	}
	tests := []struct {
		name string
		a    protocol.Attributes
		args args
		want protocol.Attributes
	}{
		{
			name: "add to nil",
			a:    nil,
			args: args{
				name: "test",
				f: null.Float{NullFloat64: sql.NullFloat64{
					Float64: 4.4,
					Valid:   true,
				}},
			},
			want: protocol.Attributes{"test": 4.4},
		},
		{
			name: "add to not nil",
			a:    protocol.Attributes{"early": "test"},
			args: args{
				name: "test",
				f: null.Float{NullFloat64: sql.NullFloat64{
					Float64: 4.4,
					Valid:   true,
				}},
			},
			want: protocol.Attributes{"early": "test", "test": 4.4},
		},
		{
			name: "invalid add to not nil",
			a:    protocol.Attributes{"early": "test"},
			args: args{
				name: "test",
				f: null.Float{NullFloat64: sql.NullFloat64{
					Float64: 4.4,
					Valid:   false,
				}},
			},
			want: protocol.Attributes{"early": "test"},
		},
		{
			name: "invalid add to not nil",
			a:    nil,
			args: args{
				name: "test",
				f: null.Float{NullFloat64: sql.NullFloat64{
					Float64: 4.4,
					Valid:   false,
				}},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.AppendNullFloat(tt.args.name, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppendNullFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributes_AppendNullFloatSlice(t *testing.T) {
	type args struct {
		name string
		s    []null.Float
	}
	tests := []struct {
		name string
		a    protocol.Attributes
		args args
		want protocol.Attributes
	}{
		{
			name: "add to nil",
			a:    nil,
			args: args{
				name: "test",
				s: []null.Float{
					{NullFloat64: sql.NullFloat64{
						Float64: 4.4,
						Valid:   true,
					}},
				},
			},
			want: protocol.Attributes{"test_0": 4.4},
		},
		{
			name: "add to not nil",
			a:    protocol.Attributes{"early": "test"},
			args: args{
				name: "test",
				s: []null.Float{
					{NullFloat64: sql.NullFloat64{
						Float64: 4.4,
						Valid:   true,
					}},
					{NullFloat64: sql.NullFloat64{
						Float64: 4.4,
						Valid:   true,
					}},
				},
			},
			want: protocol.Attributes{"early": "test", "test_0": 4.4, "test_1": 4.4},
		},
		{
			name: "add to not nil",
			a:    protocol.Attributes{"early": "test"},
			args: args{
				name: "test",
				s: []null.Float{
					{NullFloat64: sql.NullFloat64{
						Float64: 4.4,
						Valid:   false,
					}},
					{NullFloat64: sql.NullFloat64{
						Float64: 4.4,
						Valid:   true,
					}},
				},
			},
			want: protocol.Attributes{"early": "test", "test_1": 4.4},
		},
		{
			name: "invalid add to not nil",
			a:    nil,
			args: args{
				name: "test",
				s: []null.Float{
					{NullFloat64: sql.NullFloat64{
						Float64: 4.4,
						Valid:   false,
					}},
					{NullFloat64: sql.NullFloat64{
						Float64: 4.4,
						Valid:   false,
					}},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.AppendNullFloatSlice(tt.args.name, tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppendNullFloatSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributes_AppendNullString(t *testing.T) {
	type args struct {
		name string
		s    null.String
	}
	tests := []struct {
		name string
		a    protocol.Attributes
		args args
		want protocol.Attributes
	}{
		{
			name: "add to nil",
			a:    nil,
			args: args{
				name: "test",
				s: null.String{NullString: sql.NullString{
					String: "4",
					Valid:  true,
				}},
			},
			want: protocol.Attributes{"test": "4"},
		},
		{
			name: "add to not nil",
			a:    protocol.Attributes{"early": "test"},
			args: args{
				name: "test",
				s: null.String{NullString: sql.NullString{
					String: "4",
					Valid:  true,
				}},
			},
			want: protocol.Attributes{"early": "test", "test": "4"},
		},
		{
			name: "invalid add to not nil",
			a:    protocol.Attributes{"early": "test"},
			args: args{
				name: "test",
				s: null.String{NullString: sql.NullString{
					String: "4",
					Valid:  false,
				}},
			},
			want: protocol.Attributes{"early": "test"},
		},
		{
			name: "invalid add to not nil",
			a:    nil,
			args: args{
				name: "test",
				s: null.String{NullString: sql.NullString{
					String: "4",
					Valid:  false,
				}},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.AppendNullString(tt.args.name, tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppendNullString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttributes_GetFloat64(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		a       protocol.Attributes
		args    args
		wantF   float64
		wantErr bool
	}{
		{
			name:    "empty",
			a:       nil,
			args:    args{key: "test"},
			wantF:   0,
			wantErr: true,
		},
		{
			name:    "not a float",
			a:       protocol.Attributes{"test": "test"},
			args:    args{key: "test"},
			wantF:   0,
			wantErr: true,
		},
		{
			name:    "float",
			a:       protocol.Attributes{"test": 4.4},
			args:    args{key: "test"},
			wantF:   4.4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotF, err := tt.a.GetFloat64(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotF != tt.wantF {
				t.Errorf("GetFloat64() gotF = %v, want %v", gotF, tt.wantF)
			}
		})
	}
}

func TestAttributes_GetInt64(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		a       protocol.Attributes
		args    args
		wantI   int64
		wantErr bool
	}{
		{
			name:    "empty",
			a:       nil,
			args:    args{key: "test"},
			wantI:   0,
			wantErr: true,
		},
		{
			name:    "not a int",
			a:       protocol.Attributes{"test": "test"},
			args:    args{key: "test"},
			wantI:   0,
			wantErr: true,
		},
		{
			name:    "int",
			a:       protocol.Attributes{"test": int64(4)},
			args:    args{key: "test"},
			wantI:   4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotI, err := tt.a.GetInt64(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotI != tt.wantI {
				t.Errorf("GetInt64() gotI = %v, want %v", gotI, tt.wantI)
			}
		})
	}
}

func TestAttributes_GetString(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		a       protocol.Attributes
		args    args
		wantF   string
		wantErr bool
	}{
		{
			name:    "empty",
			a:       nil,
			args:    args{key: "test"},
			wantF:   "",
			wantErr: true,
		},
		{
			name:    "not a string",
			a:       protocol.Attributes{"test": 4.4},
			args:    args{key: "test"},
			wantF:   "",
			wantErr: true,
		},
		{
			name:    "string",
			a:       protocol.Attributes{"test": "4"},
			args:    args{key: "test"},
			wantF:   "4",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotF, err := tt.a.GetString(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotF != tt.wantF {
				t.Errorf("GetString() gotF = %v, want %v", gotF, tt.wantF)
			}
		})
	}
}

func TestAttributes_GetType(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		a       protocol.Attributes
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "empty",
			a:       nil,
			args:    args{key: "test"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "no key",
			a:       protocol.Attributes{"tst": 4},
			args:    args{key: "test"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "int",
			a:       protocol.Attributes{"test": int64(4)},
			args:    args{key: "test"},
			want:    "int64",
			wantErr: false,
		},
		{
			name:    "float",
			a:       protocol.Attributes{"test": 4.4},
			args:    args{key: "test"},
			want:    "float64",
			wantErr: false,
		},
		{
			name:    "string",
			a:       protocol.Attributes{"test": "4.4"},
			args:    args{key: "test"},
			want:    "string",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.GetType(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetType() got = %v, want %v", got, tt.want)
			}
		})
	}
}
