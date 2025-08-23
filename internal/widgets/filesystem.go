package widgets

import (
	"github.com/spf13/afero"
)

// FileSystem interface for testing
type FileSystem interface {
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm int) error
}

// RealFileSystem implements FileSystem using actual OS filesystem
type RealFileSystem struct{}

func (rfs *RealFileSystem) ReadFile(filename string) ([]byte, error) {
	return afero.ReadFile(afero.NewOsFs(), filename)
}

func (rfs *RealFileSystem) WriteFile(filename string, data []byte, perm int) error {
	return afero.WriteFile(afero.NewOsFs(), filename, data, 0644)
}

// TestFileSystem implements FileSystem using in-memory filesystem for testing
type TestFileSystem struct {
	fs afero.Fs
}

func NewTestFileSystem() *TestFileSystem {
	return &TestFileSystem{
		fs: afero.NewMemMapFs(),
	}
}

func (tfs *TestFileSystem) ReadFile(filename string) ([]byte, error) {
	return afero.ReadFile(tfs.fs, filename)
}

func (tfs *TestFileSystem) WriteFile(filename string, data []byte, perm int) error {
	return afero.WriteFile(tfs.fs, filename, data, 0644)
}

func (tfs *TestFileSystem) CreateFile(filename string, content string) error {
	return afero.WriteFile(tfs.fs, filename, []byte(content), 0644)
}
