package log

import (
	"os"

	"github.com/mattn/go-isatty"
	"github.com/rdeusser/x/zappretty"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() (*zap.Logger, zap.AtomicLevel) {
	atom := zap.NewAtomicLevel()
	cfg := zap.NewProductionEncoderConfig()

	zappretty.Register(cfg)

	cliEncoder := zappretty.NewCLIEncoder(cfg)
	jsonEncoder := zapcore.NewJSONEncoder(cfg)

	leveler := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= atom.Level()
	})

	core := zapcore.NewCore(jsonEncoder, os.Stdout, leveler)

	if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		core = zapcore.NewCore(cliEncoder, os.Stdout, leveler)
	}

	logger := zap.New(core)
	defer logger.Sync()

	return logger, atom
}
