package files

import (
	"io"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

//local is an implementation of the Storage interface which works
// with local disk on the current machine
type Local struct {
	maxFileSize int64 // max number of buyes for files
	basePath    string
}

// NewLocal creates a new Local filesytem with the given base path
// basePath is the base directory to save files to
// maxSize is the max number of bytes that a file can be
func NewLocal(basePath string, maxSize int64) (*Local, error) {
	p, err := filepath.Abs(basePath)

	if err != nil {
		return nil, err
	}

	return &Local{basePath: p, maxFileSize: maxSize}, nil
}

func (l *Local) Save(path string, contents io.Reader) error {
	// get full path
	fp := l.fullPath(path)

	// get the directory and make sure it exists
	d := filepath.Dir(fp)
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return xerrors.Errorf("Unable to create directory: %w", err)
	}

	// if the file exists delete it
	_, err = os.Stat(fp)
	if err == nil {
		err = os.Remove(fp)
		if err != nil {
			return xerrors.Errorf("Unable to delete to file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		//if this is anything other than a not exists error
		return xerrors.Errorf("Unable to get file info: %w", err)
	}

	//create a new file at the path
	f, err := os.Create(fp)
	if err != nil {
		return xerrors.Errorf("Unable to create file: %w", err)
	}

	//wirte the contents to the new file
	//make sure we are not writing greater than max bytes
	_, err = io.Copy(f, contents)
	if err != nil {
		return xerrors.Errorf("Unable to write to file: %w", err)
	}
	return nil
}

func (l *Local) fullPath(path string) string {

	// append the given path to the base path
	return filepath.Join(l.basePath, path)
}

func (l *Local) MaxFileSize() int64 {
	return l.maxFileSize
}
