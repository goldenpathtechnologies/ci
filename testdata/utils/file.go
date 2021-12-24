package utils

import (
	"errors"
	"github.com/goldenpathtechnologies/ci/internal/pkg/utils"
	"github.com/google/uuid"
	"io/fs"
	"math/rand"
	"runtime"
	"strings"
	"time"
)

type MockFile struct {
	FileName string
	FileSize int64
	FileMode fs.FileMode
	FileModTime time.Time
}

func (t MockFile) Name() string {
	return t.FileName
}

func (t MockFile) Size() int64 {
	return t.FileSize
}

func (t MockFile) Mode() fs.FileMode {
	return t.FileMode
}

func (t MockFile) ModTime() time.Time {
	return t.FileModTime
}

func (t MockFile) IsDir() bool {
	return t.FileMode & fs.ModeDir != 0
}

func (t MockFile) Sys() interface{} {
	return t
}

func GenerateMockFile() MockFile {
	return MockFile{
		FileName:    uuid.NewString(),
		FileSize:    rand.Int63(),
		FileMode:    0 | fs.ModePerm,
		FileModTime: time.Now(),
	}
}

func GenerateMockDirectory() MockFile {
	return MockFile{
		FileName:    uuid.NewString(),
		FileSize:    0,
		FileMode:    fs.ModeDir | fs.ModePerm,
		FileModTime: time.Now(),
	}
}

type VirtualFileSystem interface {
	Pwd() fs.FileInfo
	Cd(dirName string) (fs.FileInfo, error)
	Ls() []fs.FileInfo
	ReadLink(linkName string) (string, error)
}

type MockFileNode struct {
	File MockFile
	Children []*MockFileNode
	Parent *MockFileNode
}

type MockFileSystem struct {
	rootNode *MockFileNode
	currentDir *MockFileNode
}

func (m *MockFileSystem) Pwd() fs.FileInfo {
	return m.currentDir.File
}

func (m *MockFileSystem) Cd(dirName string) (fs.FileInfo, error) {
	dirName = normalizePath(dirName)

	switch dirName {
	case utils.OsPathSeparator:
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
	if strings.ContainsAny(dirName, utils.OsPathSeparator) {
		dirNamePart = strings.SplitN(strings.TrimRight(dirName, utils.OsPathSeparator), utils.OsPathSeparator, 2)

		if len(dirNamePart) > 1 && dirNamePart[0] == "" {
			m.currentDir = m.rootNode
			return m.Cd(dirNamePart[1])
		}
	} else {
		dirNamePart = append(dirNamePart, dirName)
	}

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

func normalizePath(path string) string {
	if utils.OsPathSeparator == "/" && strings.ContainsAny(path, "\\") {
		return strings.ReplaceAll(path, "\\", utils.OsPathSeparator)
	} else if utils.OsPathSeparator == "\\" && strings.ContainsAny(path, "/") {
		return strings.ReplaceAll(path, "/", utils.OsPathSeparator)
	} else {
		return path
	}
}

func (m *MockFileSystem) Ls() []fs.FileInfo {
	var files []fs.FileInfo

	for _, node := range m.currentDir.Children {
		files = append(files, node.File)
	}

	return files
}

func (m *MockFileSystem) ReadLink(linkName string) (string, error) {
	var (
		node *MockFileNode
		err error
	)

	linkName = normalizePath(linkName)

	if node, err = findNode(linkName, m.currentDir.Children); err != nil {
		return "", err
	}

	switch linkName {
	case utils.OsPathSeparator:
		node = m.rootNode
	case ".":
		node = m.currentDir
	case "..":
		node = m.currentDir.Parent
	}

	path := node.File.Name()

	for ok := true; ok; ok = node.Parent != nil {
		path = strings.Trim(node.Parent.File.Name(), utils.OsPathSeparator) + utils.OsPathSeparator + path
		node = node.Parent
	}

	return path, nil
}

func findNode(fileName string, nodes []*MockFileNode) (*MockFileNode, error) {
	for _, node := range nodes {
		if node.File.Name() == fileName {
			return node, nil
		}
	}
	return nil, errors.New("node not found")
}

func NewMockFileSystem(seed []*MockFileNode, depth, maxItemsPerDir int) *MockFileSystem {
	rootDir := MockFile{
		FileName:    "",
		FileSize:    0,
		FileMode:    fs.ModeDir | fs.ModePerm,
		FileModTime: time.Now(),
	}

	if runtime.GOOS == "windows" {
		rootDir.FileName = "C:" + utils.OsPathSeparator
	} else {
		rootDir.FileName = utils.OsPathSeparator
	}

	rootNode := &MockFileNode{
		File:     rootDir,
		Children: seed,
		Parent:   nil,
	}

	for i := range seed {
		seed[i].Parent = rootNode
	}

	return &MockFileSystem{
		rootNode:   generateMockFiles(rootNode, depth, maxItemsPerDir),
		currentDir: rootNode,
	}
}

func generateMockFiles(node *MockFileNode, depth, maxItemsPerDir int) *MockFileNode {
	if depth <= 0 {
		return node
	}

	numItemsPerDir := rand.Intn(maxItemsPerDir)

	if node.Children == nil {
		node.Children = []*MockFileNode{}
	}

	for i := len(node.Children) - 1; i < numItemsPerDir; i++ {
		var newNode *MockFileNode
		if rand.Intn(2) == 0 {
			newNode = &MockFileNode{
				File:     GenerateMockFile(),
				Children: nil,
				Parent:   node,
			}
		} else {
			newNode = &MockFileNode{
				File:     GenerateMockDirectory(),
				Children: nil,
				Parent:   node,
			}
		}
		node.Children = append(node.Children, newNode)
	}

	for i := range node.Children {
		if node.Children[i].File.IsDir() {
			node.Children[i] = generateMockFiles(node.Children[i], depth-1, maxItemsPerDir)
		}
	}

	return node
}
