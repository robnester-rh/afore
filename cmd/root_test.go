package cmd

import (
	"os"
	"reflect"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/robnester-rh/afore/internal/utils"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func Test_initLogging(t *testing.T) {
	utils.Filesys = afero.NewMemMapFs()
	initLogging()
	logExists, _ := utils.Filesys.Stat("/tmp/afore/afore.log")
	tests := []struct {
		name string
		want os.FileInfo
	}{
		{
			name: "Creates log file",
			want: logExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := logExists
			if !reflect.DeepEqual(logExists, tt.want) {
				t.Errorf("initLogging() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_initLogging2(t *testing.T) {
	Verbose = true
	utils.Filesys = afero.NewMemMapFs()
	initLogging()
	logExists, _ := utils.Filesys.Stat("/tmp/afore/afore.log")
	tests := []struct {
		name string
		want os.FileInfo
	}{
		{
			name: "Creates log file",
			want: logExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := logExists
			if !reflect.DeepEqual(logExists, tt.want) {
				t.Errorf("initLogging() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_initConfig(t *testing.T) {
	afs := setupInitConfig()
	_, _ = afs.Fs.Create("/home/jdoe/.afore.yaml")
	d1 := []byte("---\ntestKey: 1\n")
	err := afs.WriteFile("/home/jdoe/.afore.yaml", d1, 0755)
	if err != nil {
		t.Errorf(err.Error())
	}
	initConfig()
	got := viper.GetString("testKey")
	want := "1"
	if !reflect.DeepEqual(got, want) {
		t.Errorf("initConfig() failure. Got %s, want %s", got, want)
	}
}

func Test_initConfig2(t *testing.T) {
	afs := setupInitConfig()
	_, _ = afs.Fs.Create("/home/jdoe/.afore.yaml")
	_ = afs.MkdirAll("/path/to/", 0755)
	_, _ = afs.Fs.Create("/path/to/config.yaml")
	d1 := []byte("---\ntestKey: 1\n")
	d2 := []byte("---\ntestKey: 2\n")
	err := afs.WriteFile("/home/jdoe/.afore.yaml", d1, 0755)
	if err != nil {
		t.Errorf(err.Error())
	}
	err = afs.WriteFile("/path/to/config.yaml", d2, 0755)
	if err != nil {
		t.Errorf(err.Error())
	}
	cfgFile = "/path/to/config.yaml"
	initConfig()
	got := viper.GetString("testKey")
	want := "2"
	if !reflect.DeepEqual(got, want) {
		t.Errorf("initConfig() failure. Got %s, want %s", got, want)
	}
}

func setupInitConfig() *afero.Afero {
	afs := &afero.Afero{Fs: afero.NewMemMapFs()}
	utils.Filesys = afs
	initLogging()
	_ = os.Setenv("HOME", "/home/jdoe/")
	home, _ := homedir.Dir()
	viper.SetFs(afs.Fs)
	utils.Logger.Debugw("home directory", "home", home)
	_ = afs.Fs.MkdirAll(home, 0755)
	return afs
}
