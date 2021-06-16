package importlib

import (
	"github.com/shoriwe/gplasma/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileSystem interface {
	ChangeDirectoryRelative(string) *errors.Error
	ChangeDirectoryFullPath(string) *errors.Error
	ChangeDirectoryToFileLocation(string) *errors.Error
	ResetPath() // This should reset the path to root of the file system
	OpenRelative(string) (io.ReadSeekCloser, error)
	ExistsRelative(string) bool
	ListDirectory() ([]string, error)
	AbsolutePwd() string
	RelativePwd() string
}

type RealFileSystem struct {
	current []string
	root    string
}

func (r *RealFileSystem) ChangeDirectoryRelative(path string) *errors.Error {
	path = filepath.Clean(path)
	r.current = append(r.current, strings.Split(path, string(filepath.Separator))...)
	return nil
}

func (r *RealFileSystem) ChangeDirectoryFullPath(path string) *errors.Error {
	path = filepath.Clean(path)
	r.current = strings.Split(path, string(filepath.Separator))
	return nil
}

func (r *RealFileSystem) ChangeDirectoryToFileLocation(path string) *errors.Error {
	path = filepath.Clean(path)
	directory, _ := filepath.Split(path)
	return r.ChangeDirectoryRelative(directory)
}
func (r *RealFileSystem) ResetPath() {
	r.current = nil
}

func (r *RealFileSystem) OpenRelative(path string) (io.ReadSeekCloser, error) {
	path = filepath.Clean(path)
	filePath := filepath.Join(r.root, filepath.Join(r.current...), path)
	return os.Open(filePath)
}

func (r *RealFileSystem) ExistsRelative(path string) bool {
	path = filepath.Clean(path)
	filePath := filepath.Join(r.root, filepath.Join(r.current...), path)
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	// What if there is an error?
	panic("Unsupported behavior")
}

func (r *RealFileSystem) ListDirectory() ([]string, error) {
	filePath := filepath.Join(r.root, filepath.Join(r.current...))
	files, readError := os.ReadDir(filePath)
	if readError != nil {
		return nil, readError
	}
	var result []string
	for _, entry := range files {
		result = append(result, entry.Name())
	}
	return result, nil
}

func (r *RealFileSystem) AbsolutePwd() string {
	calculatedPath := filepath.Join(r.root, filepath.Join(r.current...))
	absolute, err := filepath.Abs(calculatedPath)
	if err != nil {
		return calculatedPath
	}
	return absolute
}

func (r *RealFileSystem) RelativePwd() string {
	return filepath.Join(filepath.Join(r.current...))
}

func NewRealFileSystem(path string) *RealFileSystem {
	path = filepath.Clean(path)
	return &RealFileSystem{
		root: path,
	}
}
