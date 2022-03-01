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
	Fields           []FileEnumType
	CloudDelimiter   string        // google cloud field
	CloudPrefix      string        // google cloud field
	CloudVersions    bool          // google cloud field
	CloudStartOffset string        // google cloud field
	CloudEndOffset   string        // google cloud field
	CloudProjection  gs.Projection // google cloud field
}

func WithFileCloudDelimiter(condition Query, delimiter string) Query {
	condition.Fields = append(condition.Fields, FileCloudDelimiter)
	condition.CloudDelimiter = delimiter
	return condition
}

func WithFileCloudPrefix(condition Query, prefix string) Query {
	condition.Fields = append(condition.Fields, FileCloudPrefix)
	condition.CloudPrefix = prefix
	return condition
}

func WithFileCloudVersions(condition Query, versions bool) Query {
	condition.Fields = append(condition.Fields, FileCloudVersions)
	condition.CloudVersions = versions
	return condition
}

func WithFileCloudStartOffset(condition Query, startOffset string) Query {
	condition.Fields = append(condition.Fields, FileCloudStartOffset)
	condition.CloudStartOffset = startOffset
	return condition
}

func WithFileCloudEndOffset(condition Query, endOffset string) Query {
	condition.Fields = append(condition.Fields, FileCloudEndOffset)
	condition.CloudEndOffset = endOffset
	return condition
}

func WithFileCloudProjection(condition Query, projection gs.Projection) Query {
	condition.Fields = append(condition.Fields, FileCloudProjection)
	condition.CloudProjection = projection
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
