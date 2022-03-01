package cloud

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
	"path/filepath"
	"time"
)

type Folder struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
}

type File struct {
	Handle      *storage.BucketHandle `json:"-"`
	FilePath    string                `json:"path,omitempty"`
	PublicURL   string                `json:"public_url,omitempty"`
	MediaLink   string                `json:"media_link,omitempty"`
	ContentType string                `json:"content_type,omitempty"`
	FileSize    int64                 `json:"size,omitempty"`
	Created     time.Time             `json:"created,omitempty"`
	Updated     time.Time             `json:"updated,omitempty"`
	Folder      *Folder               `json:"folders,omitempty"`
}

func (f *File) FolderInfo() (name string, path string, exist bool) {
	if f.Folder == nil {
		return "", "", false
	}
	return f.Folder.Name, f.Folder.Path, true
}

func (f *File) Path() string { return f.FilePath }

func (f *File) Name() string { return filepath.Base(f.FilePath) }

func (f *File) Size() (int64, error) { return f.FileSize, nil }

func (f *File) CreatedTime() (time.Time, error) { return f.Created, nil }

func (f *File) ModTime() (time.Time, error) { return f.Updated, nil }

func (f *File) NewWriter(ctx context.Context) (writer io.WriteCloser, closeFn func() error) {
	wc := f.Handle.Object(f.FilePath).NewWriter(ctx)
	return wc, func() error {
		return wc.Close()
	}
}

func (f *File) NewReader(ctx context.Context) (reader io.ReadCloser, closeFn func() error, err error) {
	rc, err := f.Handle.Object(f.FilePath).NewReader(ctx)
	if err != nil {
		return nil, nil, err
	}
	return rc, func() error {
		return rc.Close()
	}, nil
}

// Remove a file but returns ErrNotFound if not found.
func (f *File) Remove(ctx context.Context) error {
	return f.Handle.Object(f.FilePath).Delete(ctx)
}

// GetURL fetches the file URL for downloading but returns ErrNotFound if no found.
func (f *File) GetURL() string { return f.PublicURL }
