package logger

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// WriterFile writer into file
type WriterFile struct {
	// write log order by order and  atomic incr maxLinesCurLines and maxSizeCurSize
	sync.RWMutex

	logger *Logger

	Filename   string
	FileWriter *os.File

	// Rotate mode:minute|hour|day
	RotateMode string

	// Rotate minute
	MaxMinutes     int64
	minuteOpenDate int
	minuteOpenTime time.Time

	// Rotate hourly
	MaxHours       int64
	hourlyOpenDate int
	hourlyOpenTime time.Time

	// Rotate daily
	MaxDays       int64
	dailyOpenDate int
	dailyOpenTime time.Time

	// Rotate bool
	FilePerm   string
	FolderPerm string
	RotatePerm string

	// like "project.log", project is fileNameOnly and .log is suffix
	fileNameOnly string
	suffix       string
}

// NewLogFile create a LogWriter returning as os.File.
func (logger *Logger) NewLogFile(writer *WriterFile) *os.File {
	if writer.RotatePerm == "" {
		writer.RotatePerm = "0660"
	}
	if writer.FilePerm == "" {
		writer.FilePerm = "0660"
	}
	if writer.FolderPerm == "" {
		writer.FolderPerm = "0775"
	}

	if writer.RotateMode == "mimute" {
		writer.RotateMode = "minute"
	} else if writer.RotateMode == "hour" {
		writer.RotateMode = "hour"
	} else if writer.RotateMode == "day" {
		writer.RotateMode = "day"
	}

	writer.logger = logger
	writer.suffix = filepath.Ext(writer.Filename)
	writer.fileNameOnly = strings.TrimSuffix(writer.Filename, writer.suffix)
	if writer.suffix == "" {
		writer.suffix = ".log"
	}

	_ = writer.initLogFile()
	return writer.FileWriter
}

// init file logger. create log file and set to locker-inside file writer.
func (w *WriterFile) initLogFile() error {
	file, err := w.createLogFile()
	if err != nil {
		return err
	}

	if w.FileWriter != nil {
		w.FileWriter.Close()
	}
	w.logger.SetOutput(file)
	w.FileWriter = file
	return w.initFd()
}

func (w *WriterFile) needRotateMinutes(minute int) bool {
	return w.RotateMode == "minute" && minute != w.minuteOpenDate
}

func (w *WriterFile) needRotateHourly(hour int) bool {
	return w.RotateMode == "hour" && hour != w.hourlyOpenDate
}

func (w *WriterFile) needRotateDaily(day int) bool {
	return w.RotateMode == "day" && day != w.dailyOpenDate
}

func (w *WriterFile) createLogFile() (*os.File, error) {
	// Open the log file
	filePerm, err := strconv.ParseInt(w.FilePerm, 8, 64)
	if err != nil {
		return nil, err
	}

	// Open the log folder
	folderPerm, err := strconv.ParseInt(w.FolderPerm, 8, 64)
	if err != nil {
		return nil, err
	}

	filepath := path.Dir(w.Filename)
	os.MkdirAll(filepath, os.FileMode(folderPerm))

	fd, err := os.OpenFile(w.Filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(filePerm))
	if err == nil {
		// Make sure file perm is user set perm cause of `os.OpenFile` will obey umask
		os.Chmod(w.Filename, os.FileMode(filePerm))
	}
	return fd, err
}

func (w *WriterFile) initFd() error {
	fd := w.FileWriter
	_, err := fd.Stat()
	if err != nil {
		// log.Printf("get stat err: %s", err)
		return fmt.Errorf("get stat err: %s", err)
	}

	w.minuteOpenTime = time.Now()
	w.minuteOpenDate = w.minuteOpenTime.Minute()
	w.hourlyOpenTime = time.Now()
	w.hourlyOpenDate = w.hourlyOpenTime.Hour()
	w.dailyOpenTime = time.Now()
	w.dailyOpenDate = w.dailyOpenTime.Day()

	if w.RotateMode == "minute" {
		go w.minuteRotate(w.minuteOpenTime)
	} else if w.RotateMode == "hour" {
		go w.hourlyRotate(w.hourlyOpenTime)
	} else if w.RotateMode == "day" {
		go w.dailyRotate(w.dailyOpenTime)
	}

	return nil
}

func (w *WriterFile) minuteRotate(openTime time.Time) {
	y, m, d := openTime.Add(1 * time.Minute).Date()
	h, i, _ := openTime.Add(1 * time.Minute).Clock()
	nextMinute := time.Date(y, m, d, h, i, 0, 0, openTime.Location())
	tm := time.NewTimer(time.Duration(nextMinute.UnixNano() - openTime.UnixNano() + 100))
	<-tm.C
	w.Lock()
	if w.needRotateMinutes(time.Now().Minute()) {
		if err := w.doRotate(time.Now()); err != nil {
			fmt.Fprintf(os.Stderr, "WriterFile(%q): %s\n", w.Filename, err)
		}
	}
	w.Unlock()
}

func (w *WriterFile) hourlyRotate(openTime time.Time) {
	y, m, d := openTime.Add(1 * time.Hour).Date()
	h, _, _ := openTime.Add(1 * time.Hour).Clock()
	nextHour := time.Date(y, m, d, h, 0, 0, 0, openTime.Location())
	tm := time.NewTimer(time.Duration(nextHour.UnixNano() - openTime.UnixNano() + 100))
	<-tm.C
	w.Lock()
	if w.needRotateHourly(time.Now().Hour()) {
		if err := w.doRotate(time.Now()); err != nil {
			fmt.Fprintf(os.Stderr, "WriterFile(%q): %s\n", w.Filename, err)
		}
	}
	w.Unlock()
}

func (w *WriterFile) dailyRotate(openTime time.Time) {
	y, m, d := openTime.Add(24 * time.Hour).Date()
	nextDay := time.Date(y, m, d, 0, 0, 0, 0, openTime.Location())
	tm := time.NewTimer(time.Duration(nextDay.UnixNano() - openTime.UnixNano() + 100))
	<-tm.C
	w.Lock()
	if w.needRotateDaily(time.Now().Day()) {
		if err := w.doRotate(time.Now()); err != nil {
			fmt.Fprintf(os.Stderr, "WriterFile(%q): %s\n", w.Filename, err)
		}
	}
	w.Unlock()
}

// DoRotate means it need to write file in new file: new file name like xx.2013-01-01.log (daily) or xx.001.log (by line or size)
func (w *WriterFile) doRotate(logTime time.Time) error {
	fName := ""
	format := ""

	var openTime time.Time
	rotatePerm, err := strconv.ParseInt(w.RotatePerm, 8, 64)
	if err != nil {
		return err
	}

	_, err = os.Lstat(w.Filename)
	if err != nil {
		// even if the file is not exist or other ,we should RESTART the logger
		goto RESTART_LOGGER
	}

	if w.RotateMode == "minute" {
		format = "200601021504"
		openTime = w.minuteOpenTime
	} else if w.RotateMode == "hour" {
		format = "2006010215"
		openTime = w.hourlyOpenTime
	} else if w.RotateMode == "day" {
		format = "20060102"
		openTime = w.dailyOpenTime
	}

	fName = w.fileNameOnly + fmt.Sprintf(".%s%s", openTime.Format(format), w.suffix)
	_, err = os.Lstat(fName)

	// return error if the last file checked still existed
	if err == nil {
		// log.Printf("rotate: Cannot find free log number to rename %s", w.Filename)
		return fmt.Errorf("rotate: Cannot find free log number to rename %s", w.Filename)
	}

	// close FileWriter before rename
	w.FileWriter.Close()

	// Rename the file to its new found name
	// even if occurs error,we MUST guarantee to  restart new logger
	err = os.Rename(w.Filename, fName)
	if err != nil {
		goto RESTART_LOGGER
	}

	err = os.Chmod(fName, os.FileMode(rotatePerm))

RESTART_LOGGER:

	initLogErr := w.initLogFile()

	if initLogErr != nil {
		// log.Printf("rotate StartLogger: %s", initLogErr)
		return fmt.Errorf("rotate StartLogger: %s", initLogErr)
	}

	if err != nil {
		// log.Printf("rotate: %s", err)
		return fmt.Errorf("rotate: %s", err)
	}
	return nil
}
