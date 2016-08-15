package l

import (
	"path/filepath"
	"runtime"

	"github.com/Sirupsen/logrus"
)

func Locate(fields logrus.Fields) logrus.Fields {
	_, path, line, ok := runtime.Caller(1)
	if ok {
		_, file := filepath.Split(path)
		fields["file"] = file
		fields["line"] = line
	}
	return fields
}
