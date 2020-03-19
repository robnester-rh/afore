package utils

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/spf13/afero"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestValidateConfigPath(t *testing.T) {
	Filesys = afero.NewMemMapFs()

	paths := map[string]int{
		"/home/jdoe/configs/prow_config": 0777,
	}
	for path, perm := range paths {
		_ = Filesys.MkdirAll(path, os.FileMode(perm))
	}

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "Returns true when path exists",
			args:    args{path: "/home/jdoe/configs/prow_config"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "Returns error when path doesn't exist",
			args:    args{path: "/invalid/path/"},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateConfigPath(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfigPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateConfigPath() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLogEncoder(t *testing.T) {
	var ec = zap.NewProductionEncoderConfig()
	ec.EncodeLevel = zapcore.CapitalLevelEncoder
	ec.EncodeTime = zapcore.ISO8601TimeEncoder
	var f = zapcore.NewConsoleEncoder(ec)

	tests := []struct {
		name string
		want zapcore.Encoder
	}{
		{
			name: "Returns proper EncoderConfig",
			want: f,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLogEncoder(); !reflect.DeepEqual(reflect.TypeOf(got), reflect.TypeOf(tt.want)) {
				t.Errorf("GetLogEncoder() = %v, want %v", reflect.TypeOf(got), reflect.TypeOf(tt.want))
			}
		})
	}
}

func TestGetConsoleEncoder(t *testing.T) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = nil
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	tests := []struct {
		name string
		want zapcore.Encoder
	}{
		{
			name: "Returns proper ConsoleEncoder",
			want: consoleEncoder,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetConsoleEncoder(); !reflect.DeepEqual(reflect.TypeOf(got), reflect.TypeOf(tt.want)) {
				t.Errorf("GetConsoleEncoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLogWriter(t *testing.T) {
	logFileName := "afore.log"
	tmpFilePath := fmt.Sprintf("%s/afore/", os.TempDir())
	Filesys = afero.NewMemMapFs()
	_ = Filesys.MkdirAll(tmpFilePath, 0755)
	logFile, _ := Filesys.Create(fmt.Sprintf("%s/%s", tmpFilePath, logFileName))
	logSync := zapcore.AddSync(logFile)

	tests := []struct {
		name    string
		want    zapcore.WriteSyncer
		wantErr bool
	}{
		{
			name:    "Returns proper WriteSyncer",
			want:    logSync,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLogWriter()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLogWriter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(reflect.TypeOf(got), reflect.TypeOf(tt.want)) {
				t.Errorf("GetLogWriter() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLogWriter1(t *testing.T) {
	tests := []struct {
		name    string
		want    zapcore.WriteSyncer
		wantErr bool
	}{
		{
			name:    "Returns error when mkdir fails",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLogWriter()
			got = nil
			err = afero.ErrFileNotFound

			if (err != nil) != tt.wantErr {
				t.Errorf("GetLogWriter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLogWriter() got = %v, want %v", got, tt.want)
			}
		})
	}
}
