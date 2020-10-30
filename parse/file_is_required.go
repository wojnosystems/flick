package parse

import (
	"errors"
	"github.com/wojnosystems/go-optional"
	"os"
)

var ErrNoFile = errors.New("file not provided, but was required")

type fileIsRequired struct {
	pathToFile optional.String
	original   FileUnmarshaler
}

func (o *fileIsRequired) Unmarshal(c interface{}) (err error) {
	if !o.pathToFile.IsSet() {
		return ErrNoFile
	}
	fileHandle, err := os.Open(o.pathToFile.Value())
	if err != nil {
		return
	}
	defer func() {
		_ = fileHandle.Close()
	}()
	err = o.original.UnmarshalFile(fileHandle, c)
	return
}

func FileIsRequired(pathToFile optional.String, unmarshaller FileUnmarshaler) Unmarshaler {
	return &fileIsRequired{
		pathToFile: pathToFile,
		original:   unmarshaller,
	}
}
