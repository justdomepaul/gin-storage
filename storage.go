package gin_storage

import (
	"encoding/json"
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/justdomepaul/gin-storage/pkg/panicerrorhandler"
	"github.com/justdomepaul/gin-storage/storage"
	"net/http"
)

const (
	// DefaultPrefix url prefix of pprof
	DefaultPrefix = "/storage"
)

func getPrefix(prefixOptions ...string) string {
	prefix := DefaultPrefix
	if len(prefixOptions) > 0 {
		prefix = prefixOptions[0]
	}
	return prefix
}

func Register(r *gin.Engine, prefixOptions ...string) (closeFn func()) {
	return RouteRegister(&(r.RouterGroup), prefixOptions...)
}

func RouteRegister(rg *gin.RouterGroup, prefixOptions ...string) (closeFn func()) {
	fileStorage, fn := storage.Load()
	handler := NewFileHandler(fileStorage)

	prefixRouter := rg.Group(getPrefix(prefixOptions...))
	{
		prefixRouter.POST("", handler.Upload)
		prefixRouter.PUT("", handler.Publicize)
		prefixRouter.DELETE("", handler.Remove)
		prefixRouter.GET("", handler.List)
	}

	return fn
}

func NewFileHandler(storage storage.IFile) *FileHandler {
	return &FileHandler{
		storage: storage,
	}
}

type FileHandler struct {
	storage storage.IFile
}

func (fh FileHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		panic(panicerrorhandler.NewErrExecute(err))
	}
	f, err := file.Open()
	if err != nil {
		panic(panicerrorhandler.NewErrExecute(err))
	}
	path, err := fh.storage.Upload(c, c.Request.FormValue("prefix"), f)
	if err != nil {
		panic(panicerrorhandler.NewErrDBExecute(err))
	}
	c.JSON(http.StatusOK, gin.H{
		"path": path,
	})
}

func (fh FileHandler) Publicize(c *gin.Context) {
	req := struct {
		Path string `json:"path,omitempty" validate:"required"`
	}{}
	defer c.Request.Body.Close()
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		panic(panicerrorhandler.NewErrJSONUnmarshal(err))
	}
	if err := validator.New().Struct(&req); err != nil {
		panic(panicerrorhandler.NewErrVariable(err))
	}
	url, err := fh.storage.GetURL(c, req.Path)
	if err != nil {
		panic(panicerrorhandler.NewErrDBExecute(err))
	}
	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

func (fh FileHandler) Remove(c *gin.Context) {
	req := struct {
		Path string `json:"path,omitempty" validate:"required"`
	}{}
	defer c.Request.Body.Close()
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		panic(panicerrorhandler.NewErrJSONUnmarshal(err))
	}
	if err := validator.New().Struct(&req); err != nil {
		panic(panicerrorhandler.NewErrVariable(err))
	}
	if err := fh.storage.Remove(c, req.Path); err != nil {
		panic(panicerrorhandler.NewErrDBExecute(err))
	}
	c.String(http.StatusOK, "ok")
}

func (fh FileHandler) List(c *gin.Context) {
	q := storage.Query{}
	if c.Query("delimiter") != "" {
		q = storage.WithFileCloudDelimiter(q, c.Query("delimiter"))
	}
	if c.Query("prefix") != "" {
		q = storage.WithFileCloudPrefix(q, c.Query("prefix"))
	}

	var files []storage.File
	if err := fh.storage.List(c, q, func(file storage.File) error {
		files = append(files, file)
		return nil
	}); errors.Is(err, panicerrorhandler.ErrGetFile) {
		panic(panicerrorhandler.NewErrDBRowNotFound(err))
	} else if err != nil {
		panic(panicerrorhandler.NewErrDBExecute(err))
	}
	c.JSON(http.StatusOK, files)
}
