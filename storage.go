package gin_storage

import (
	"encoding/json"
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	errorhandlerTool "github.com/justdomepaul/gin-storage/pkg/errorhandler"
	"github.com/justdomepaul/gin-storage/storage"
	"github.com/justdomepaul/toolbox/errorhandler"
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
		prefixRouter.GET("", handler.List)
		prefixRouter.POST("", handler.Upload)
		prefixRouter.POST("/multiple", handler.Batch)
		prefixRouter.PUT("", handler.Publicize)
		prefixRouter.PUT("/multiple", handler.MultiplePublicize)
		prefixRouter.DELETE("/:id", handler.Remove)
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
		panic(errorhandler.NewErrVariable(err))
	}
	f, err := file.Open()
	if err != nil {
		panic(errorhandler.NewErrExecute(err))
	}
	path, err := fh.storage.Upload(c, c.Request.FormValue("prefix"), f)
	if err != nil {
		panic(errorhandler.NewErrDBExecute(err))
	}
	c.JSON(http.StatusOK, gin.H{
		"path": path,
	})
}

type BatchFile struct {
	Filename string `json:"filename,omitempty" validate:"required"`
	Path     string `json:"path,omitempty" validate:"required"`
}

type PublicizeURL struct {
	Filename string `json:"filename,omitempty" validate:"required"`
	URL      string `json:"url,omitempty" validate:"required"`
}

func (fh FileHandler) Batch(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		panic(errorhandler.NewErrVariable(err))
	}
	var responsePaths []BatchFile
	for _, file := range form.File["file[]"] {
		f, err := file.Open()
		if err != nil {
			panic(errorhandler.NewErrExecute(err))
		}
		path, err := fh.storage.Upload(c, c.Request.FormValue("prefix"), f)
		if err != nil {
			panic(errorhandler.NewErrDBExecute(err))
		}
		responsePaths = append(responsePaths, BatchFile{
			Filename: file.Filename,
			Path:     path,
		})
	}
	c.JSON(http.StatusOK, responsePaths)
}

func (fh FileHandler) Publicize(c *gin.Context) {
	req := struct {
		Path string `json:"path,omitempty" validate:"required"`
	}{}
	defer c.Request.Body.Close()
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		panic(errorhandler.NewErrJSONUnmarshal(err))
	}
	if err := validator.New().Struct(&req); err != nil {
		panic(errorhandler.NewErrVariable(err))
	}
	url, err := fh.storage.GetURL(c, req.Path)
	if err != nil {
		panic(errorhandler.NewErrDBExecute(err))
	}
	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

func (fh FileHandler) MultiplePublicize(c *gin.Context) {
	req := struct {
		Paths []BatchFile `json:"paths,omitempty" validate:"required,min=1,dive"`
	}{}
	defer c.Request.Body.Close()
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		panic(errorhandler.NewErrJSONUnmarshal(err))
	}
	if err := validator.New().Struct(&req); err != nil {
		panic(errorhandler.NewErrVariable(err))
	}
	var responseURLs []PublicizeURL
	for _, path := range req.Paths {
		url, err := fh.storage.GetURL(c, path.Path)
		if err != nil {
			panic(errorhandler.NewErrDBExecute(err))
		}
		responseURLs = append(responseURLs, PublicizeURL{
			Filename: path.Filename,
			URL:      url,
		})
	}
	c.JSON(http.StatusOK, responseURLs)
}

func (fh FileHandler) Remove(c *gin.Context) {
	req := struct {
		Path string `validate:"required"`
	}{
		Path: c.Param("id"),
	}
	defer c.Request.Body.Close()
	if err := validator.New().Struct(&req); err != nil {
		panic(errorhandler.NewErrVariable(err))
	}
	if err := fh.storage.Remove(c, req.Path); err != nil {
		panic(errorhandler.NewErrDBExecute(err))
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
	}); errors.Is(err, errorhandlerTool.ErrGetFile) {
		panic(errorhandler.NewErrDBRowNotFound(err))
	} else if err != nil {
		panic(errorhandler.NewErrDBExecute(err))
	}
	c.JSON(http.StatusOK, files)
}
