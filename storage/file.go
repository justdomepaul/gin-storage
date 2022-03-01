package storage

import (
	gs "cloud.google.com/go/storage"
	"context"
	"io"
	"time"
)

type FileEnumType int

const (
	FileCloudDelimiter FileEnumType = iota
	FileCloudPrefix
	FileCloudVersions
	FileCloudStartOffset
	FileCloudEndOffset
	FileCloudProjection
)

type Query struct {
	Fields      []FileEnumType
	Delimiter   string        // google cloud field
	Prefix      string        // google cloud field
	Versions    bool          // google cloud field
	StartOffset string        // google cloud field
	EndOffset   string        // google cloud field
	Projection  gs.Projection // google cloud field
}

func WithFileDelimiter(condition Query, delimiter string) Query {
	condition.Fields = append(condition.Fields, FileCloudDelimiter)
	condition.Delimiter = delimiter
	return condition
}

func WithFilePrefix(condition Query, prefix string) Query {
	condition.Fields = append(condition.Fields, FileCloudPrefix)
	condition.Prefix = prefix
	return condition
}

func WithFileVersions(condition Query, versions bool) Query {
	condition.Fields = append(condition.Fields, FileCloudVersions)
	condition.Versions = versions
	return condition
}

func WithFileStartOffset(condition Query, startOffset string) Query {
	condition.Fields = append(condition.Fields, FileCloudStartOffset)
	condition.StartOffset = startOffset
	return condition
}

func WithFileEndOffset(condition Query, endOffset string) Query {
	condition.Fields = append(condition.Fields, FileCloudEndOffset)
	condition.EndOffset = endOffset
	return condition
}

func WithFileProjection(condition Query, projection gs.Projection) Query {
	condition.Fields = append(condition.Fields, FileCloudProjection)
	condition.Projection = projection
	return condition
}

type File interface {
	FolderInfo() (name string, path string, exist bool)
	Path() string
	Name() string
	Size() (int64, error)
	CreatedTime() (time.Time, error)
	ModTime() (time.Time, error)
	NewWriter(ctx context.Context) (writer io.WriteCloser, closeFn func() error)
	NewReader(ctx context.Context) (reader io.ReadCloser, closeFn func() error, err error)
	Remove(ctx context.Context) error
	GetURL() string
}

type IterHandler func(file File) error

type IFile interface {
	Upload(ctx context.Context, prefix string, f io.ReadCloser) (string, error)
	GetURL(ctx context.Context, path string) (string, error)
	Remove(ctx context.Context, path string) error
	List(ctx context.Context, q Query, h IterHandler) error
}
