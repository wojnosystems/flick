package parse

import (
	"errors"
	"github.com/wojnosystems/go-optional/v2"
	"os"
)

var ErrNoFile = errors.New("file not provided, but was required")

type fileIsRequired struct {
	pathToFile optional.String
	original   FileUnmarshaler
}

func (o *fileIsRequired) Unmarshal(c interface{}) (err error) {
	o.pathToFile.IfSetElse(func(path string) {
		var fileHandle *os.File
		fileHandle, err = os.Open(path)
		if err != nil {
			return
		}
		defer func() {
			_ = fileHandle.Close()
		}()
		err = o.original.UnmarshalFile(fileHandle, c)
	}, func() {
		err = ErrNoFile
	})
	return
}

func FileIsRequired(pathToFile optional.String, unmarshaller FileUnmarshaler) Unmarshaler {
	return &fileIsRequired{
		pathToFile: pathToFile,
		original:   unmarshaller,
	}
}
