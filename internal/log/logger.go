package log

import (
	"os"
	"path"
	"sync"

	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//DefaultLog ..
type DefaultLog struct {
	*zap.SugaredLogger
}

var log *DefaultLog
var name string
var logOnce sync.Once

//AppName 模块名
var AppName string = "default"

func setName() {
	name = AppName
}

//AppLog AppLog
func AppLog() *DefaultLog {
	setName()
	logOnce.Do(getLog)
	return log
}

func getLog() {
	errPri := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zapcore.ErrorLevel
	})

	infoPri := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zapcore.InfoLevel
	})

	debugPri := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zapcore.DebugLevel
	})

	infoWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path.Join(viper.GetString("log-path"), name+".log"), // 日志文件路径
		MaxSize:    100,                                                 // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 5,                                                   // 日志文件最多保存多少个备份
		MaxAge:     7,                                                   // 文件最多保存多少天
		Compress:   true,                                                // 是否压缩
		LocalTime:  true,
	})

	errWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path.Join(viper.GetString("log-path"), name+"-error.log"), // 日志文件路径
		MaxSize:    100,                                                       // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 5,                                                         // 日志文件最多保存多少个备份
		MaxAge:     7,                                                         // 文件最多保存多少天
		Compress:   true,                                                      // 是否压缩
		LocalTime:  true,
	})

	var cfg zapcore.EncoderConfig
	if viper.GetBool("release") {
		cfg = zap.NewProductionEncoderConfig()
	} else {
		cfg = zap.NewDevelopmentEncoderConfig()
	}

	enc := zapcore.NewConsoleEncoder(cfg)
	debugWriter := zapcore.Lock(os.Stdout)
	core := zapcore.NewTee(
		zapcore.NewCore(enc, debugWriter, debugPri),
		zapcore.NewCore(enc, infoWriter, infoPri),
		zapcore.NewCore(enc, errWriter, errPri),
	)

	z1 := zap.New(core, zap.AddCaller())
	log = &DefaultLog{
		SugaredLogger: z1.Sugar(),
	}
}
