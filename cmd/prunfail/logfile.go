// chris 090515

package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"io/ioutil"
)

const maxLogSize = 16384

// maxLogSize is small enough that we can afford to just read the entire
// thing into memory given the benefit in simplicity.

type logFile struct {
	path     string
	data     string
	failures int
}

var errBadLog = errors.New("corrupted log file")

func newLogFile(path string) (*logFile, error) {
	flag := os.O_RDONLY | os.O_CREATE
	file, err := os.OpenFile(path, flag, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		return nil, err
	}
	lf := &logFile{path: path, data: string(data)}
	if err := lf.parse(); err != nil {
		return nil, err
	}
	return lf, nil
}

func (lf *logFile) parse() error {
	const footerprefix = "fail "

	if len(lf.data) == 0 {
		lf.failures = 0
		return nil
	}

	i := strings.LastIndex(lf.data, footerprefix)
	if i == -1 {
		lf.failures = 0
		return errBadLog
	}

	x := lf.data[i+len(footerprefix):]
	if !strings.HasSuffix(x, "\n") {
		return errBadLog
	}
	x = x[:len(x)-1] // Strip off trailing newline.

	failures64, err := strconv.ParseInt(x, 0, 0)
	if err != nil {
		return errors.New("corrupted log file: " + err.Error())
	}
	lf.failures = int(failures64)
	return nil
}

func (lf *logFile) write(data []byte) error {
	flag := os.O_WRONLY | os.O_TRUNC | os.O_CREATE
	file, err := os.OpenFile(lf.path, flag, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	newdata := lf.data + string(data)
	if len(newdata) > maxLogSize {
		newdata = newdata[len(newdata)-maxLogSize:]
	}
	newdata += fmt.Sprintf("%s fail %d\n", time.Now(), lf.failures)
	_, err2 := io.WriteString(file, newdata)
	return err2
}
