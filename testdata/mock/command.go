package mock

import (
	"io/fs"
)

type DirectoryCommands struct {
	ReadDirectoryFunc   func(dirname string) ([]fs.FileInfo, error)
	GetAbsolutePathFunc func(path string) (string, error)
	ScanDirectoryFunc   func(path string, callback func(dirName string)) error
}

func (m *DirectoryCommands) ReadDirectory(dirname string) ([]fs.FileInfo, error) {
	return m.ReadDirectoryFunc(dirname)
}

func (m *DirectoryCommands) GetAbsolutePath(path string) (string, error) {
	return m.GetAbsolutePathFunc(path)
}

func (m *DirectoryCommands) ScanDirectory(path string, callback func(dirName string)) error {
	return m.ScanDirectoryFunc(path, callback)
}

func NewDirectoryCommandsForVirtualFileSystem(fileSystem VirtualFileSystem) *DirectoryCommands {
	return &DirectoryCommands{
		ReadDirectoryFunc: func(dirname string) ([]fs.FileInfo, error) {
			return fileSystem.Ls(dirname)
		},
		GetAbsolutePathFunc: func(path string) (string, error) {
			return fileSystem.ReadLink(path)
		},
		ScanDirectoryFunc: func(path string, callback func(dirName string)) error {
			files, err := fileSystem.Ls(path)
			if err != nil {
				return err
			}
			for _, file := range files {
				callback(file.Name())
			}
			return nil
		},
	}
}