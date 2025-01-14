package testlog

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/ten-protocol/go-ten/lib/gethfork/debug"

	gethlog "github.com/ethereum/go-ethereum/log"
)

// must be set during setup
var (
	logFile = ""
	testlog gethlog.Logger
)

func Logger() gethlog.Logger {
	return testlog
}

func LogFile() string {
	return logFile
}

type Cfg struct {
	LogDir      string // directory for the log file
	TestType    string // type of test (comes before timestamp in filename so sorted file list will block these together)
	TestSubtype string // test subtype (comes after timestamp in filename so sorted file list will show latest of different subtypes together)
	LogLevel    slog.Level
}

// Setup will direct logs to a timestamped log file with a standard naming pattern, useful for simulations etc.
func Setup(cfg *Cfg) *os.File {
	// Create a folder for the logs if none exists.
	err := os.MkdirAll(cfg.LogDir, 0o700)
	if err != nil {
		panic(err)
	}
	timeFormatted := time.Now().Format("2006-01-02_15-04-05")
	f, err := os.CreateTemp(cfg.LogDir, fmt.Sprintf("%s-%s-%s-*.txt", cfg.TestType, timeFormatted, cfg.TestSubtype))
	if err != nil {
		panic(err)
	}
	logFile = f.Name()

	err = debug.Setup("terminal", logFile, false, 10000000, 0, 0, false, false, cfg.LogLevel, "")
	if err != nil {
		panic(err)
	}

	testlog = gethlog.New()
	return f
}

// SetupSysOut will direct the test logs to stdout
func SetupSysOut() {
	err := debug.Setup("terminal", "", false, 10000000, 0, 0, false, false, slog.LevelDebug, "")
	if err != nil {
		panic(err)
	}
	testlog = gethlog.New()
}
