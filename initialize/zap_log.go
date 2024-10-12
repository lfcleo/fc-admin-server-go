package initialize

import (
	"fc-admin-server-go/pkg/zaplogger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func ZapLogSetUp() *zap.Logger {
	//获取编码器
	encoder := zaplogger.GetEncoder()

	//info文件WriteSyncer
	infoFileWriteSyncer := zaplogger.GetInfoWriterSyncer()
	//error文件WriteSyncer
	errorFileWriteSyncer := zaplogger.GetErrorWriterSyncer()
	//panic文件WriteSyncer
	panicFileWriteSyncer := zaplogger.GetPanicWriterSyncer()
	//fatal文件WriteSyncer
	fatalFileWriteSyncer := zaplogger.GetFatalWriterSyncer()

	//同时输出到控制台 和 指定的日志文件中
	infoFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), zap.InfoLevel)
	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), zap.ErrorLevel)
	panicFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(panicFileWriteSyncer, zapcore.AddSync(os.Stdout)), zap.PanicLevel)
	fatalFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(fatalFileWriteSyncer, zapcore.AddSync(os.Stdout)), zap.FatalLevel)

	//将infoFileCore， errorFileCore ，panicFileCore 加入core切片
	var coreArr []zapcore.Core
	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, errorFileCore)
	coreArr = append(coreArr, panicFileCore)
	coreArr = append(coreArr, fatalFileCore)

	//生成Logger
	return zap.New(zapcore.NewTee(coreArr...), zap.AddCaller()) //zap.AddCaller() 显示文件名 和 行号
}
