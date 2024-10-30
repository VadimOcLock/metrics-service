package migrations

import "testing"

func Test_dbNameByDSN(t *testing.T) {
	type args struct {
		dsn string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Valid DSN",
			args:    args{dsn: "postgres://postgres:postgres@localhost:5432/postgres_db?sslmode=disable"},
			want:    "postgres_db",
			wantErr: false,
		},
		{
			name:    "Invalid scheme",
			args:    args{dsn: "mysql://user:password@localhost:3306/db"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "No database name",
			args:    args{dsn: "postgres://postgres:postgres@localhost:5432/"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Empty DSN",
			args:    args{dsn: ""},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Invalid URL",
			args:    args{dsn: "invalid_dsn"},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dbNameByDSN(tt.args.dsn)
			if (err != nil) != tt.wantErr {
				t.Errorf("dbNameByDSN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("dbNameByDSN() got = %v, want %v", got, tt.want)
			}
		})
	}
}
