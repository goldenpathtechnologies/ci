package utils

import (
	"github.com/google/uuid"
	"io/fs"
	"math/rand"
	"time"
)

type TestFile struct {
	FileName string
	FileSize int64
	FileMode fs.FileMode
	FileModTime time.Time
}

func (t TestFile) Name() string {
	return t.FileName
}

func (t TestFile) Size() int64 {
	return t.FileSize
}

func (t TestFile) Mode() fs.FileMode {
	return t.FileMode
}

func (t TestFile) ModTime() time.Time {
	return t.FileModTime
}

func (t TestFile) IsDir() bool {
	return t.FileMode & fs.ModeDir != 0
}

func (t TestFile) Sys() interface{} {
	return t
}

func GenerateTestFile() TestFile {
	return TestFile{
		FileName:    uuid.NewString(),
		FileSize:    rand.Int63(),
		FileMode:    0 | fs.ModePerm,
		FileModTime: time.Now(),
	}
}

func GenerateTestDirectory() TestFile {
	return TestFile{
		FileName:    uuid.NewString(),
		FileSize:    0,
		FileMode:    fs.ModeDir | fs.ModePerm,
		FileModTime: time.Now(),
	}
}
