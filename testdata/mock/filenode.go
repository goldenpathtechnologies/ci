package mock

import (
	"errors"
	"io/fs"
	"math/rand"
	"strconv"
	"time"
)

type FileNode struct {
	File     File
	Children []*FileNode
	Parent   *FileNode
}

func findNode(fileName string, nodes []*FileNode) (*FileNode, error) {
	// TODO: Create type called MockFileNodeCollection of type []*MockFileNode
	//  and add this func as a member. I may also want to replace references of
	//  the slice with those of the new type.
	for _, node := range nodes {
		if node.File.Name() == fileName {
			return node, nil
		}
	}
	return nil, errors.New("node not found")
}

func generateMockFiles(node *FileNode, depth, maxItemsPerDir int) *FileNode {
	if depth <= 0 || maxItemsPerDir <= 0 {
		return node
	}

	numItemsPerDir := rand.Intn(maxItemsPerDir)

	if node.Children == nil {
		node.Children = []*FileNode{}
	} else {
		for i := 0; i < len(node.Children); i++ {
			if node.Children[i].Parent == nil {
				node.Children[i].Parent = node
			}
		}
	}

	for i := len(node.Children) - 1; i < numItemsPerDir; i++ {
		var newNode *FileNode
		if rand.Intn(2) == 0 {
			newNode = &FileNode{
				File:     GenerateMockFile(),
				Children: nil,
				Parent:   node,
			}
		} else {
			newNode = &FileNode{
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

func GenerateSeedDirectories(fileNamePrefix string, count int) []*FileNode {
	var directories []*FileNode

	for i := 0; i < count; i++ {
		directories = append(directories, &FileNode{
			File: File{
				FileName:    fileNamePrefix + strconv.Itoa(i),
				FileSize:    0,
				FileMode:    fs.ModeDir | fs.ModePerm,
				FileModTime: time.Now(),
			},
			Children: nil,
			Parent:   nil,
		})
	}

	return directories
}

func GetSampleExampleSeedDirectories() []*FileNode {
	return []*FileNode{
		{
			File: File{
				FileName:    "sample0",
				FileSize:    0,
				FileMode:    fs.ModeDir | fs.ModePerm,
				FileModTime: time.Now(),
			},
			Children: nil,
			Parent:   nil,
		},
		{
			File: File{
				FileName:    "sample1",
				FileSize:    0,
				FileMode:    fs.ModeDir | fs.ModePerm,
				FileModTime: time.Now(),
			},
			Children: nil,
			Parent:   nil,
		},
		{
			File: File{
				FileName:    "example0",
				FileSize:    0,
				FileMode:    fs.ModeDir | fs.ModePerm,
				FileModTime: time.Now(),
			},
			Children: nil,
			Parent:   nil,
		},
		{
			File: File{
				FileName:    "example1",
				FileSize:    0,
				FileMode:    fs.ModeDir | fs.ModePerm,
				FileModTime: time.Now(),
			},
			Children: nil,
			Parent:   nil,
		},
		{
			File: File{
				FileName:    "example2",
				FileSize:    0,
				FileMode:    fs.ModeDir | fs.ModePerm,
				FileModTime: time.Now(),
			},
			Children: nil,
			Parent:   nil,
		},
	}
}

func GetHierarchicalSeedDirectories() []*FileNode {
	return []*FileNode{
		{
			File: File{
				FileName:    "testA",
				FileSize:    0,
				FileMode:    fs.ModeDir | fs.ModePerm,
				FileModTime: time.Now(),
			},
			Children: []*FileNode{
				{
					File: File{
						FileName:    "testB",
						FileSize:    0,
						FileMode:    fs.ModeDir | fs.ModePerm,
						FileModTime: time.Now(),
					},
					Children: []*FileNode{
						{
							File: File{
								FileName:    "testC",
								FileSize:    0,
								FileMode:    fs.ModeDir | fs.ModePerm,
								FileModTime: time.Now(),
							},
							Children: nil,
							Parent:   nil,
						},
					},
					Parent: nil,
				},
			},
			Parent: nil,
		},
	}
}
