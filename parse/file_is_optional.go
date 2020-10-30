package parse

import (
	"github.com/wojnosystems/go-optional"
	"os"
)

type fileIsOptional struct {
	pathToFile optional.String
	original   FileUnmarshaler
}

func (o *fileIsOptional) Unmarshal(c interface{}) (err error) {
	if !o.pathToFile.IsSet() {
		return
	}
	fileHandle, err := os.Open(o.pathToFile.Value())
	if err != nil {
		err = nil
		return
	}
	defer func() {
		_ = fileHandle.Close()
	}()
	err = o.original.UnmarshalFile(fileHandle, c)
	return
}

func FileIsOptional(pathToFile optional.String, unmarshaler FileUnmarshaler) Unmarshaler {
	return &fileIsOptional{
		pathToFile: pathToFile,
		original:   unmarshaler,
	}
}
