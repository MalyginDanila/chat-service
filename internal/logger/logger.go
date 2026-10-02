package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level — глобальный изменяемый уровень логирования.
// zap.AtomicLevel реализует http.Handler (GET/PUT), что используется в debug-сервере.
var Level = zap.NewAtomicLevel()

type Options struct {
	Level          string
	ProductionMode bool
}

func MustInit(opts Options) {
	if err := Init(opts); err != nil {
		panic(err)
	}
}

func Init(opts Options) error {
	lvl, err := zapcore.ParseLevel(opts.Level)
	if err != nil {
		return fmt.Errorf("parse level %q: %v", opts.Level, err)
	}
	Level.SetLevel(lvl)

	encCfg := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "level",
		NameKey:        "component",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	var enc zapcore.Encoder
	if opts.ProductionMode {
		encCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		enc = zapcore.NewJSONEncoder(encCfg)
	} else {
		encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		enc = zapcore.NewConsoleEncoder(encCfg)
	}

	cores := []zapcore.Core{
		zapcore.NewCore(enc, zapcore.Lock(os.Stdout), Level),
	}
	// В модуле 2 сюда добавится core для Sentry.

	l := zap.New(zapcore.NewTee(cores...), zap.AddStacktrace(zapcore.ErrorLevel))
	zap.ReplaceGlobals(l)
	return nil
}

func Sync() {
	if err := zap.L().Sync(); err != nil {
		// Sync на stdout в Windows/Linux иногда возвращает "invalid argument" — это нормально.
		_ = err
	}
}
