package mock

import (
	"github.com/google/uuid"
	"io/fs"
	"math/rand"
	"time"
)

type File struct {
	FileName    string
	FileSize    int64
	FileMode    fs.FileMode
	FileModTime time.Time
}

func (t File) Name() string {
	return t.FileName
}

func (t File) Size() int64 {
	return t.FileSize
}

func (t File) Mode() fs.FileMode {
	return t.FileMode
}

func (t File) ModTime() time.Time {
	return t.FileModTime
}

func (t File) IsDir() bool {
	return t.FileMode & fs.ModeDir != 0
}

func (t File) Sys() interface{} {
	return t
}

func GenerateMockFile() File {
	return File{
		FileName:    uuid.NewString(),
		FileSize:    rand.Int63(),
		FileMode:    0 | fs.ModePerm,
		FileModTime: time.Now(),
	}
}

func GenerateMockDirectory() File {
	return File{
		FileName:    uuid.NewString(),
		FileSize:    0,
		FileMode:    fs.ModeDir | fs.ModePerm,
		FileModTime: time.Now(),
	}
}
