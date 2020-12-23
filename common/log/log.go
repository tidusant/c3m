package log

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type LogLevel int

const (
	DebugLevel LogLevel = 0
	InfoLevel  LogLevel = 1
	WarnLevel  LogLevel = 2
	ErrorLevel LogLevel = 3
)

const (
	DefaultLevel = DebugLevel
)

var (
	OutputFile            *os.File
	IncludeDateInFileName = true
	AppendExisting        = false
	outputFilename        string
	outputFilenameDate    string
	messagePrefix         string
	outputLevel           LogLevel = DefaultLevel

	debugLogger, infoLogger, warnLogger, errorLogger *log.Logger
)

func DefaultOutput(level LogLevel) {
	SetOutput(os.Stderr, level)
}

func RedirectStdOut() {
	if OutputFile == nil {
		return
	}

	os.Stderr = OutputFile
	os.Stdout = OutputFile
}

func SetOutputFile(name string, level LogLevel) {
	// Ensure the logs directory exists
	fi, err := os.Stat("logs")
	if err != nil || !fi.IsDir() {
		os.Mkdir("logs", os.ModeDir|os.ModePerm)
	}

	// Close previous log file (if created)
	CloseOutputFile()

	now := time.Now().Format("20060102")
	var logFilePath string
	if IncludeDateInFileName {
		logFilePath = fmt.Sprintf(`logs%c%s-%s.log`, os.PathSeparator, strings.ToLower(name), now)
	} else {
		logFilePath = fmt.Sprintf(`logs%c%s.log`, os.PathSeparator, strings.ToLower(name))
	}

	if !AppendExisting {
		existingCount := 1
		originalLogFilePath := logFilePath
		if _, err = os.Stat(logFilePath); err == nil {
			for {
				existingCount++
				logFilePath = strings.Replace(originalLogFilePath, ".log", fmt.Sprintf("-%d.log", existingCount), -1)
				if _, err = os.Stat(logFilePath); os.IsNotExist(err) {
					break
				}
			}
		}

		OutputFile, err = os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	} else {
		OutputFile, err = os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	}

	if err != nil {
		log.Fatalf("Error creating log file: %v", err)
		panic(fmt.Sprintf("Error creating log file: %v", err))
		return
	}

	outputFilename = name
	outputLevel = level
	outputFilenameDate = now

	SetOutput(OutputFile, level)
}

func SetOutput(w io.Writer, level LogLevel) {
	log.SetOutput(w)

	debugLogger = getLogger(w, "TRACE ", level > DebugLevel)
	infoLogger = getLogger(w, "INFO ", level > InfoLevel)
	warnLogger = getLogger(w, "WARN ", level > WarnLevel)
	errorLogger = getLogger(w, "ERROR ", level > ErrorLevel)
}
func getLogger(w io.Writer, prefix string, discard bool) *log.Logger {
	if discard {
		return log.New(ioutil.Discard, "", 0)
	}

	return log.New(w, prefix, log.Ldate|log.Ltime|log.Lshortfile)
}

func CloseOutputFile() {
	if OutputFile != nil {
		OutputFile.Close()
		OutputFile = nil
	}
}

func SetPrefix(prefix string) {
	messagePrefix = prefix
}
func LogMsg(s string) {
	Println(time.Now().Format("2006-01-02 15:04:05") + " >> " + s)
}
func Println(args ...interface{}) {
	beforeLog()
	log.Println(args...)
}
func Printf(format string, args ...interface{}) {
	beforeLog()
	Println(fmt.Sprintf(format, args...))
}
func Debug(args ...interface{}) {
	logPrint(debugLogger, args...)
}
func Debugf(format string, args ...interface{}) {
	logPrintf(debugLogger, format, args...)
}
func Info(args ...interface{}) {
	logPrint(infoLogger, args...)
}
func Infof(format string, args ...interface{}) {
	logPrintf(infoLogger, format, args...)
}
func Warn(args ...interface{}) {
	logPrint(warnLogger, args...)
}
func Warnf(format string, args ...interface{}) {
	logPrintf(warnLogger, format, args...)
}
func Error(args ...interface{}) {
	logPrint(errorLogger, args...)
}
func Errorf(format string, args ...interface{}) {
	logPrintf(errorLogger, format, args...)
}

func beforeLog() {
	// Check if we need to change the output log file name
	if !IncludeDateInFileName || OutputFile == nil || outputFilenameDate == "" {
		return
	}

	now := time.Now().Format("20060102")
	if now == outputFilenameDate {
		return
	}

	SetOutputFile(outputFilename, outputLevel)
}
func logPrint(logger *log.Logger, args ...interface{}) {
	beforeLog()
	if messagePrefix != "" {
		args = append([]interface{}{messagePrefix}, args...)
	}

	if logger == nil {
		log.Print(args...)
	} else {
		logger.Print(args...)
	}
}
func logPrintf(logger *log.Logger, format string, args ...interface{}) {
	beforeLog()
	if messagePrefix != "" {
		format = messagePrefix + format
	}

	if logger == nil {
		log.Printf(format, args...)
	} else {
		logger.Printf(format, args...)
	}
}
