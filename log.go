package log

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Logger struct {
	mu     sync.Mutex
	out    io.Writer
	format string
	utc    bool
	buf    []byte
}

func New(out io.Writer, format string) *Logger {
	return &Logger{out: out, format: format}
}

func (l *Logger) UseUTC() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.utc = true
}

func (l *Logger) Output(s string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	if l.utc {
		now = now.UTC()
	}
	l.buf = l.buf[:0]
	l.buf = append(l.buf, now.Format(l.format)...)
	l.buf = append(l.buf, ' ')
	l.buf = append(l.buf, s...)
	if l.buf[len(l.buf)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.out.Write(l.buf)
	return err
}

func (l *Logger) Fatal(v ...interface{}) {
	l.Output(fmt.Sprint(v...))
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Output(fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l *Logger) Fatalln(v ...interface{}) {
	l.Output(fmt.Sprintln(v...))
	os.Exit(1)
}

func (l *Logger) Print(v ...interface{}) {
	l.Output(fmt.Sprint(v...))
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.Output(fmt.Sprintf(format, v...))
}

func (l *Logger) Println(v ...interface{}) {
	l.Output(fmt.Sprintln(v...))
}

var defaultLogger = New(os.Stderr, time.RFC3339)

func UseUTC() {
	defaultLogger.UseUTC()
}

func Output(s string) error {
	return defaultLogger.Output(s)
}

func Fatal(v ...interface{}) {
	defaultLogger.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	defaultLogger.Fatalf(format, v...)
}

func Fatalln(v ...interface{}) {
	defaultLogger.Fatalln(v...)
}

func Print(v ...interface{}) {
	defaultLogger.Print(v...)
}

func Printf(format string, v ...interface{}) {
	defaultLogger.Printf(format, v...)
}

func Println(v ...interface{}) {
	defaultLogger.Println(v...)
}
