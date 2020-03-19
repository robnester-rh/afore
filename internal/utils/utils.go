package utils

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Filesys = afero.NewOsFs()
var Logger *zap.SugaredLogger
var err error

func ValidateConfigPath(path string) (bool, error) {
	valid := true
	afs := &afero.Afero{Fs: Filesys}
	_, err = afs.Fs.Stat(path)
	if err != nil {
		valid = false
		err = errors.Errorf("Error validating config path: %s", err.Error())
	}
	return valid, err
}

func GetLogEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func GetConsoleEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = nil
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func GetLogWriter() (zapcore.WriteSyncer, error) {
	afs := &afero.Afero{Fs: Filesys}
	logFileName := "afore.log"
	tmpFilePath := fmt.Sprintf("%s/afore/", os.TempDir())

	_ = afs.Fs.MkdirAll(tmpFilePath, 0755)
	logFile, _ := afs.Create(fmt.Sprintf("%s/%s", tmpFilePath, logFileName))
	return zapcore.AddSync(logFile), nil
}

func GenerateCommand(name string, run func(cmd *cobra.Command, args []string)) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: run,
	}
}
