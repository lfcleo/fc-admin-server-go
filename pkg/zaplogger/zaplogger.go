package zaplogger

import (
	"fc-admin-server-go/pkg/config"
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"
	"time"
)

//https://blog.csdn.net/weixin_48536164/article/details/126588267

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func levelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var level string
	switch l {
	case zapcore.DebugLevel:
		level = "[DEBUG]"
	case zapcore.InfoLevel:
		level = "[INFO]"
	case zapcore.WarnLevel:
		level = "[WARN]"
	case zapcore.ErrorLevel:
		level = "[ERROR]"
	case zapcore.DPanicLevel:
		level = "[DPANIC]"
	case zapcore.PanicLevel:
		level = "[PANIC]"
	case zapcore.FatalLevel:
		level = "[FATAL]"
	default:
		level = fmt.Sprintf("[LEVEL(%d)]", l)
	}
	enc.AppendString(level)
}

func shortCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("[%s]", caller.TrimmedPath()))
}

func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    levelEncoder, //zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,  //指定时间格式
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   shortCallerEncoder, //zapcore.ShortCallerEncoder,
	}
}

// GetEncoder core 三个参数之  Encoder 获取编码器
func GetEncoder() zapcore.Encoder {
	//自定义编码配置,下方NewJSONEncoder输出如下的日志格式
	//{"L":"[INFO]","T":"2022-09-16 14:24:59.552","C":"[prototest/main.go:113]","M":"name = xiaoli, age = 18"}
	return zapcore.NewJSONEncoder(NewEncoderConfig())

	//下方NewConsoleEncoder输出如下的日志格式
	//2022-09-16 14:26:02.933 [INFO]  [prototest/main.go:113] name = xiaoli, age = 18
	//return zapcore.NewConsoleEncoder(NewEncoderConfig())
}

// GetInfoWriterSyncer core 三个参数之  日志输出路径
func GetInfoWriterSyncer() zapcore.WriteSyncer {
	//引入第三方库 Lumberjack 加入日志切割功能，日志文件每 10MB 会切割并且在当前目录下最多保存 5 个日志文件
	infoLumberIO := &lumberjack.Logger{
		Filename:   config.Data.Zap.InfoFilename,
		MaxSize:    config.Data.Zap.MaxSize,    //文件最大内存
		MaxBackups: config.Data.Zap.MaxBackups, //最大备份数
		MaxAge:     config.Data.Zap.MaxAge,     //保存最大天数
		Compress:   false,                      //Compress确定是否应该使用gzip压缩已旋转的日志文件。默认值是不执行压缩。
	}
	return zapcore.AddSync(infoLumberIO)
}

func GetErrorWriterSyncer() zapcore.WriteSyncer {
	//引入第三方库 Lumberjack 加入日志切割功能，日志文件每 10MB 会切割并且在当前目录下最多保存 5 个日志文件
	lumberWriteSyncer := &lumberjack.Logger{
		Filename:   config.Data.Zap.ErrorFilename,
		MaxSize:    config.Data.Zap.MaxSize,    //文件最大内存
		MaxBackups: config.Data.Zap.MaxBackups, //最大备份数
		MaxAge:     config.Data.Zap.MaxAge,     //保存最大天数
		Compress:   false,                      //Compress确定是否应该使用gzip压缩已旋转的日志文件。默认值是不执行压缩。
	}
	return zapcore.AddSync(lumberWriteSyncer)
}

func GetPanicWriterSyncer() zapcore.WriteSyncer {
	//引入第三方库 Lumberjack 加入日志切割功能，日志文件每 10MB 会切割并且在当前目录下最多保存 5 个日志文件
	lumberWriteSyncer := &lumberjack.Logger{
		Filename:   config.Data.Zap.PanicFilename,
		MaxSize:    config.Data.Zap.MaxSize,    //文件最大内存
		MaxBackups: config.Data.Zap.MaxBackups, //最大备份数
		MaxAge:     config.Data.Zap.MaxAge,     //保存最大天数
		Compress:   false,                      //Compress确定是否应该使用gzip压缩已旋转的日志文件。默认值是不执行压缩。
	}
	return zapcore.AddSync(lumberWriteSyncer)
}

func GetFatalWriterSyncer() zapcore.WriteSyncer {
	//引入第三方库 Lumberjack 加入日志切割功能，日志文件每 10MB 会切割并且在当前目录下最多保存 5 个日志文件
	lumberWriteSyncer := &lumberjack.Logger{
		Filename:   config.Data.Zap.FatalFilename,
		MaxSize:    config.Data.Zap.MaxSize,    //文件最大内存
		MaxBackups: config.Data.Zap.MaxBackups, //最大备份数
		MaxAge:     config.Data.Zap.MaxAge,     //保存最大天数
		Compress:   false,                      //Compress确定是否应该使用gzip压缩已旋转的日志文件。默认值是不执行压缩。
	}
	return zapcore.AddSync(lumberWriteSyncer)
}
