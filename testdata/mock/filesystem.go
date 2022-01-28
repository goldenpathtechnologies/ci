package mock

import (
	"errors"
	"io/fs"
	"os"
	"strings"
	"time"
)

const osPathSeparator = string(os.PathSeparator)

type VirtualFileSystem interface {
	Pwd() fs.FileInfo
	Cd(dirName string) (fs.FileInfo, error)
	Ls(dirName string) ([]fs.FileInfo, error)
	ReadLink(linkName string) (string, error)
}

type FileSystem struct {
	rootNode   *FileNode
	currentDir *FileNode
}

func NewMockFileSystem(seed []*FileNode, depth, maxItemsPerDir int) *FileSystem {
	rootDir := File{
		FileName:    osPathSeparator,
		FileSize:    0,
		FileMode:    fs.ModeDir | fs.ModePerm,
		FileModTime: time.Now(),
	}

	rootNode := &FileNode{
		File:     rootDir,
		Children: seed,
		Parent:   nil,
	}

	for i := range seed {
		seed[i].Parent = rootNode
	}

	return &FileSystem{
		rootNode:   generateMockFiles(rootNode, depth, maxItemsPerDir),
		currentDir: rootNode,
	}
}

func (m *FileSystem) Pwd() fs.FileInfo {
	return m.currentDir.File
}

func (m *FileSystem) Cd(dirName string) (fs.FileInfo, error) {
	dirName = NormalizePath(dirName)

	switch dirName {
	case osPathSeparator:
		m.currentDir = m.rootNode
		return m.currentDir.File, nil
	case ".":
		return m.currentDir.File, nil
	case "..":
		if m.currentDir.Parent == nil {
			return nil, errors.New("currently at the root directory")
		} else {
			parentDir := m.currentDir.Parent
			m.currentDir = parentDir
			return m.currentDir.File, nil
		}
	}

	var dirNamePart []string
	if strings.ContainsAny(dirName, osPathSeparator) {
		dirNamePart = strings.SplitN(strings.TrimRight(dirName, osPathSeparator), osPathSeparator, 2)

		if len(dirNamePart) > 1 && dirNamePart[0] == "" {
			m.currentDir = m.rootNode
			return m.Cd(dirNamePart[1])
		}
	} else {
		dirNamePart = append(dirNamePart, dirName)
	}

	// TODO: Refactor to use findNode() here
	for _, node := range m.currentDir.Children {
		if node.File.Name() == dirNamePart[0] {
			if node.File.IsDir() {
				m.currentDir = node
				if len(dirNamePart) > 1 {
					return m.Cd(dirNamePart[1])
				} else {
					return node.File, nil
				}
			} else {
				return nil, errors.New("'" + node.File.Name() + "' is not a directory")
			}
		}
	}
	return nil, errors.New("'" + dirNamePart[0] + "' does not exist")
}

func (m *FileSystem) Ls(dirName string) ([]fs.FileInfo, error) {
	var (
		files []fs.FileInfo
		node  *FileNode
	)

	dirName = NormalizePath(dirName)

	switch dirName {
	case osPathSeparator:
		node = m.rootNode
	case "..":
		node = m.currentDir.Parent
	case ".":
		node = m.currentDir
	case "":
		node = m.currentDir
	default:
		if strings.HasPrefix(dirName, osPathSeparator) {
			node = m.rootNode
		} else {
			node = m.currentDir
		}
		nodeNames := strings.Split(strings.Trim(dirName, osPathSeparator), osPathSeparator)
		for _, name := range nodeNames {
			var err error
			if node, err = findNode(name, node.Children); err != nil {
				return nil, errors.New("'" + name + "' is not a directory")
			}
		}
	}

	for _, n := range node.Children {
		files = append(files, n.File)
	}

	return files, nil
}

func (m *FileSystem) ReadLink(linkName string) (string, error) {
	var (
		node *FileNode
		err  error
	)

	linkName = NormalizePath(linkName)

	switch linkName {
	case osPathSeparator:
		node = m.rootNode
	case ".":
		node = m.currentDir
	case "..":
		node = m.currentDir.Parent
	default:
		if node, err = findNode(linkName, m.currentDir.Children); err != nil {
			return "", err
		}
	}

	path := node.File.Name()

	if node.Parent == nil {
		return path, nil
	}

	for ok := true; ok; ok = node.Parent != nil {
		path = strings.Trim(node.Parent.File.Name(), osPathSeparator) + osPathSeparator + path
		node = node.Parent
	}

	return path, nil
}

func NormalizePath(path string) string {
	if osPathSeparator == "/" && strings.ContainsAny(path, "\\") {
		return strings.ReplaceAll(path, "\\", osPathSeparator)
	} else if osPathSeparator == "\\" && strings.ContainsAny(path, "/") {
		return strings.ReplaceAll(path, "/", osPathSeparator)
	} else {
		return path
	}

}