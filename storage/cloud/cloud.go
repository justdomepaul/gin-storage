package cloud

import (
	gs "cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/justdomepaul/gin-storage/pkg/config"
	"github.com/justdomepaul/gin-storage/pkg/database/cloud"
	"github.com/justdomepaul/gin-storage/pkg/panicerrorhandler"
	zapTool "github.com/justdomepaul/gin-storage/pkg/zap"
	"github.com/justdomepaul/gin-storage/storage"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/api/iterator"
	"io"
	"net/url"
	"path"
	"regexp"
	"strings"
)

const StorageDomain = "https://storage.googleapis.com/"

func init() {
	f, fn, err := getFile()
	if err != nil {
		panic(err)
	}
	storage.Register(f, fn)
}

func getFile() (storage.IFile, func(), error) {
	cg := config.Core{}
	media := config.Media{}
	st := config.Cloud{}

	for _, item := range []interface{}{&cg, &media, &st} {
		err := envconfig.Process("", item)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %s", panicerrorhandler.ErrInitialFileClient, err.Error())
		}
	}
	logger, err := zapTool.NewLogger(cg)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", panicerrorhandler.ErrInitialFileClient, err.Error())
	}
	cloudStorage, fn, err := cloud.NewExtendStorageDatabase(logger, st)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", panicerrorhandler.ErrInitialFileClient, err.Error())
	}
	return NewFile(media, cloudStorage), fn, nil
}

var fileClauseFn = map[storage.FileEnumType]func(source storage.Query, condition *gs.Query) error{
	storage.FileCloudDelimiter:   withFileDelimiter,
	storage.FileCloudPrefix:      withFilePrefix,
	storage.FileCloudVersions:    withFileVersions,
	storage.FileCloudStartOffset: withFileStartOffset,
	storage.FileCloudEndOffset:   withFileEndOffset,
	storage.FileCloudProjection:  withFileProjection,
}

func withFileDelimiter(source storage.Query, condition *gs.Query) error {
	if err := validator.New().Var(source.Delimiter, `required`); err != nil {
		return err
	}
	condition.Delimiter = source.Delimiter
	suffix := ""
	if !strings.HasSuffix(source.Prefix, "/") {
		suffix = "/"
	}
	condition.Prefix = source.Prefix + suffix
	return nil
}

func withFilePrefix(source storage.Query, condition *gs.Query) error {
	if err := validator.New().Var(source.Prefix, `required`); err != nil {
		return err
	}
	if !strings.HasPrefix(source.Delimiter, "/") {
		condition.Prefix = source.Prefix
	}
	return nil
}

func withFileVersions(source storage.Query, condition *gs.Query) error {
	condition.Versions = source.Versions
	return nil
}

func withFileStartOffset(source storage.Query, condition *gs.Query) error {
	if err := validator.New().Var(source.StartOffset, `required`); err != nil {
		return err
	}
	condition.StartOffset = source.StartOffset
	return nil
}

func withFileEndOffset(source storage.Query, condition *gs.Query) error {
	if err := validator.New().Var(source.EndOffset, `required`); err != nil {
		return err
	}
	condition.EndOffset = source.EndOffset
	return nil
}

func withFileProjection(source storage.Query, condition *gs.Query) error {
	condition.Projection = source.Projection
	return nil
}

func toFileClauses(source storage.Query) (*gs.Query, error) {
	q := &gs.Query{}
	for _, op := range source.Fields {
		if err := fileClauseFn[op](source, q); err != nil {
			return q, err
		}
	}
	return q, nil
}

// NewFile method
func NewFile(env config.Media, session cloud.ISession) *Cloud {
	return &Cloud{
		env:     env,
		session: session,
	}
}

type Cloud struct {
	env     config.Media
	session cloud.ISession
}

func (st *Cloud) Upload(ctx context.Context, prefix string, f io.ReadCloser) (string, error) {
	mediaID, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("%w: %s", panicerrorhandler.ErrFailGenerateUUID, err.Error())
	}
	pt := mediaID.String()
	if st.env.PrefixPath != "" {
		pt = path.Join(st.env.PrefixPath, prefix, mediaID.String())
	}
	if err := verifyPath(pt); err != nil {
		return "", fmt.Errorf("%w: %s", panicerrorhandler.ErrFileUpload, err.Error())
	}
	wc := st.session.Bucket(st.env.BucketName).Object(pt).NewWriter(ctx)
	if _, err := io.Copy(wc, f); err != nil {
		return "", fmt.Errorf("%w: %s", panicerrorhandler.ErrFileUpload, err.Error())
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("%w: %s", panicerrorhandler.ErrFailCloseSession, err.Error())
	}
	return pt, nil
}

func (st *Cloud) GetURL(ctx context.Context, route string) (string, error) {
	if err := verifyPath(route); err != nil {
		return "", fmt.Errorf("%w: %s", panicerrorhandler.ErrFileUpdate, err.Error())
	}
	_, err := st.session.Bucket(st.env.BucketName).Object(route).Update(ctx, gs.ObjectAttrsToUpdate{
		PredefinedACL: "publicRead",
	})
	if errors.Is(err, gs.ErrObjectNotExist) {
		return "", panicerrorhandler.ErrFileNotExist
	}
	if err != nil {
		return "", fmt.Errorf("%w: %s", panicerrorhandler.ErrFileUpdate, err.Error())
	}
	return getPublicURL(st.env.BucketName, route)
}

func (st *Cloud) Remove(ctx context.Context, route string) error {
	if err := verifyPath(route); err != nil {
		return fmt.Errorf("%w: %s", panicerrorhandler.ErrFileRemove, err.Error())
	}
	if err := st.session.Bucket(st.env.BucketName).Object(route).Delete(ctx); err != nil {
		return fmt.Errorf("%w: %s", panicerrorhandler.ErrFileRemove, err.Error())
	}
	return nil
}

func (st *Cloud) List(ctx context.Context, query storage.Query, h storage.IterHandler) error {
	q, err := toFileClauses(query)
	if err != nil {
		return err
	}
	bucket := st.session.Bucket(st.env.BucketName)
	return iterFiles(ctx, bucket, q, h)
}

func iterFiles(ctx context.Context, handler *gs.BucketHandle, q *gs.Query, h storage.IterHandler) error {
	iter := handler.Objects(ctx, q)
	for {
		attrs, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return fmt.Errorf("%w: %s", panicerrorhandler.ErrGetFile, err.Error())
		}
		publicURL, err := getPublicURL(attrs.Bucket, attrs.Name)
		if err != nil {
			return fmt.Errorf("%w: %s", panicerrorhandler.ErrGetFile, err.Error())
		}
		node := &File{}
		if attrs.Prefix == "" {
			node.Handle = handler
			node.FilePath = attrs.Name
			node.PublicURL = publicURL
			node.MediaLink = attrs.MediaLink
			node.ContentType = attrs.ContentType
			node.FileSize = attrs.Size
			node.Created = attrs.Created
			node.Updated = attrs.Updated
		} else {
			node.Folder = &Folder{
				Name: strings.TrimSuffix(strings.TrimPrefix(attrs.Prefix, q.Prefix), "/"),
				Path: attrs.Prefix,
			}
		}
		err = h(node)
		if err != nil {
			return fmt.Errorf("%w: %s", panicerrorhandler.ErrGetFile, err.Error())
		}
	}
	return nil
}

func getPublicURL(bucketName, route string) (string, error) {
	u, err := url.Parse(StorageDomain)
	if err != nil {
		return "", fmt.Errorf("%w: %s", panicerrorhandler.ErrFileUpdate, err.Error())
	}
	u.Path = path.Join(u.Path, bucketName, route)
	return u.String(), nil
}

func verifyPath(path string) error {
	// strong condition by Cloud Storage
	if len(path) <= 0 ||
		len(path) > 1024 ||
		strings.HasPrefix(path, ".well-known/acme-challenge/") ||
		path == "." ||
		path == ".." {
		return errors.New("invalid path")
	}
	// soft
	matched, err := regexp.MatchString(`^(/*[\w\-.()$%& ]+)+$`, path)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("invalid path")
	}
	return nil
}
