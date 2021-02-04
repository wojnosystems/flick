package parse

import (
	"github.com/wojnosystems/go-optional/v2"
	"os"
)

type fileIsOptional struct {
	pathToFile optional.String
	original   FileUnmarshaler
}

func (o *fileIsOptional) Unmarshal(c interface{}) (err error) {
	o.pathToFile.IfSet(func(path string) {
		var fileHandle *os.File
		fileHandle, err = os.Open(path)
		if err != nil {
			err = nil
			return
		}
		defer func() {
			_ = fileHandle.Close()
		}()
		err = o.original.UnmarshalFile(fileHandle, c)
	})
	return
}

func FileIsOptional(pathToFile optional.String, unmarshaler FileUnmarshaler) Unmarshaler {
	return &fileIsOptional{
		pathToFile: pathToFile,
		original:   unmarshaler,
	}
}
