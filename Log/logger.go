package log

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)
var Logx = NewLog("../log","Aquarius",4,10,1024*1024*100,LOG_SHIFT_BY_DAY)

const (
	LOG_SHIFT_BY_SIZE = iota
	LOG_SHIFT_BY_DAY
	LOG_SHIFT_BY_HOUR
	LOG_SHIFT_BY_MINUTE
)

const (
	LOG_LEVEL_ERROR = iota
	LOG_LEVEL_WARN
	LOG_LEVEL_INFO
	LOG_LEVEL_DEBUG
	LOG_LEVEL_STDOUT
)

type Log struct {
	log_path       string
	log_prefix     string
	log_level      int
	log_num        int
	log_size       int64
	log_shift_type int
	file           *os.File
	log_filename   string
	mutex          *sync.Mutex
}

func NewLog(log_path string, log_prefix string, log_level int, log_num int, log_size int64,
	log_shift_type int) *Log {
	l := new(Log)
	l.log_path = log_path
	l.log_prefix = log_prefix
	l.log_level = log_level
	l.log_num = log_num
	l.log_size = log_size
	l.log_shift_type = log_shift_type
	l.mutex = new(sync.Mutex)
	return l
}

func (l *Log) Close() {
	if l.file != nil {
		l.file.Close()
	}
}
//写入数据到file
func (l *Log) RawLog(content string) (err error) {
	if l.log_level != LOG_LEVEL_STDOUT {
		var file *os.File
		//在此会shift
		if file, err = l.getLogFile(); err != nil {
			return err
		}

		l.mutex.Lock()
		defer l.mutex.Unlock()
		if _, err = fmt.Fprintln(file, content); err != nil {
			return err
		}
	} else {
		fmt.Println(content)
	}
	return nil
}


func (l *Log) LogDebug(format string, v ...interface{}) {
	if l.log_level >= LOG_LEVEL_DEBUG {
		l.logFormat(LOG_LEVEL_DEBUG, format, v...)
	}
}

func (l *Log) LogInfo(format string, v ...interface{}) {
	if l.log_level >= LOG_LEVEL_INFO {
		l.logFormat(LOG_LEVEL_INFO, format, v...)
	}
}

func (l *Log) LogWarn(format string, v ...interface{}) {
	if l.log_level >= LOG_LEVEL_WARN {
		l.logFormat(LOG_LEVEL_WARN, format, v...)
	}
}

func (l *Log) LogError(format string, v ...interface{}) {
	if l.log_level >= LOG_LEVEL_ERROR {
		l.logFormat(LOG_LEVEL_ERROR, format, v...)
	}
}


func (l *Log) Log(format string, v ...interface{}) {
	l.logFormat(LOG_LEVEL_DEBUG, format, v...)
}

func (l *Log) logFormat(log_level int, format string, v ...interface{}) {
	log_content := l.formatLogContent(log_level, format, v...)
	l.RawLog(log_content)
}
//构造待写入数据
func (l *Log) formatLogContent(log_level int, format string, v ...interface{}) string {
	var log_content string
	log_content = time.Now().Format("[2006-01-02] 15:04:05.99999")
	switch log_level {
	case LOG_LEVEL_ERROR:
		log_content += " [ERROR] "
	case LOG_LEVEL_WARN:
		log_content += " [WARN] "
	case LOG_LEVEL_DEBUG:
		log_content += " [DEBUG] "
	case LOG_LEVEL_INFO:
		log_content += " [INFO] "
	default:
		log_content += " [LOG] "
	}
	log_content += fmt.Sprintf(format, v...)
	//runtime.Caller 返回当前 goroutine 的栈上的函数调用信息
	_, full_filename, line, _ := runtime.Caller(3)
	full_filename_arr := bytes.Split([]byte(full_filename), []byte(`/`))
	short_filename := string(full_filename_arr[len(full_filename_arr)-1])
	log_content += fmt.Sprintf(" %s:%d", short_filename, line)
	return log_content
}
//设置logfile
func (l *Log) reopen(filename string) (err error) {
	if l.file != nil {
		l.file.Close()
	}
	l.file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	return nil
}

func (l *Log) shiftLogFile(filename string) (err error) {
	if filename != l.log_filename {
		if err = l.reopen(filename); err != nil {
			return err
		}
		l.log_filename = filename
	} else if l.file != nil {
		var stat os.FileInfo
		if stat, err = l.file.Stat(); err != nil {
			return err
		}

		if l.log_shift_type == LOG_SHIFT_BY_SIZE {
			if stat.Size() >= l.log_size {
				l.mutex.Lock()
				defer l.mutex.Unlock()
				for i := l.log_num - 2; i >= 0; i-- {
					var next_filename string = l.getLogFileName(i)
					if _, err := os.Stat(filename); err == nil {
						os.Rename(next_filename, l.getLogFileName(i+1))
					}
				}
				if err = l.reopen(filename); err != nil {
					return err
				}
			}
		} else {
			var shift_time_str string //上个日志时间
			var out_date_time string //删除日志时间
			do_shift := false
			now := time.Now()
			if l.log_shift_type == LOG_SHIFT_BY_DAY {
				if now.Format("20060102") != stat.ModTime().Format("20060102") {
					shift_time_str = now.Add(-24 * time.Hour).Format("20060102")
					out_date_time = now.Add(-24 * time.Duration(l.log_num) * time.Hour).Format("20060102")
					do_shift = true
				}
			} else if l.log_shift_type == LOG_SHIFT_BY_HOUR {
				if now.Format("20060102-15") != stat.ModTime().Format("20060102-15") {
					shift_time_str = now.Add(-1 * time.Hour).Format("20060102-15")
					out_date_time = now.Add(-1 * time.Duration(l.log_num) * time.Hour).Format("20060102-15")
					do_shift = true
				}
			} else if l.log_shift_type == LOG_SHIFT_BY_MINUTE {
				if now.Format("20060102-1504") != stat.ModTime().Format("20060102-1504") {
					shift_time_str = now.Add(-1 * time.Minute).Format("20060102-1504")
					out_date_time = now.Add(-1 * time.Duration(l.log_num) * time.Minute).Format("20060102-1504")
					do_shift = true
				}
			}

			if do_shift {
				shift_filename := filename + "." + shift_time_str
				del_filename := filename + "." + out_date_time
				l.mutex.Lock()
				defer l.mutex.Unlock()
				go os.Remove(del_filename)
				os.Rename(filename, shift_filename)
				if err = l.reopen(filename); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (l *Log) getLogFile() (file *os.File, err error) {
	log_filename := l.getLogFileName(0)
	if err = l.shiftLogFile(log_filename); err != nil {
		return nil, err
	}

	return l.file, nil
}

func (l *Log) getLogFileName(idx int) (name string) {
	if idx == 0 {
		return fmt.Sprintf("%s/%s%s", l.log_path, l.log_prefix, ".log")
	} else {
		return fmt.Sprintf("%s/%s%d%s", l.log_path, l.log_prefix, idx, ".log")
	}
}

