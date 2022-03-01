package gin_storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/justdomepaul/gin-storage/pkg/panicerrorhandler"
	"github.com/justdomepaul/gin-storage/storage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

type testIFile struct {
	mock.Mock
	storage.IFile
}

func (t *testIFile) Upload(ctx context.Context, prefix string, f io.ReadCloser) (string, error) {
	args := t.Called(ctx, prefix, f)
	return args.Get(0).(string), args.Error(1)
}

func (t *testIFile) GetURL(ctx context.Context, path string) (string, error) {
	args := t.Called(ctx, path)
	return args.Get(0).(string), args.Error(1)
}

func (t *testIFile) Remove(ctx context.Context, path string) error {
	args := t.Called(ctx, path)
	return args.Error(0)
}

func (t *testIFile) List(ctx context.Context, q storage.Query, h storage.IterHandler) error {
	args := t.Called(ctx, q, h)
	return args.Error(0)
}

func NewMockGinServer() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(
		gin.Logger(),
		panicerrorhandler.GinPanicErrorHandler("system", "Mock Test Gin Server"))
	return router
}

// Get method
func Get(uri string, headers map[string]string, router *gin.Engine) ([]byte, error) {
	req := httptest.NewRequest(http.MethodGet, uri, nil)
	return getBody(req, headers, router)
}

// PostFile method
func PostFile(uri string, forms map[string]io.Reader, headers map[string]string, router *gin.Engine) ([]byte, error) {
	var (
		b   bytes.Buffer
		err error
	)
	w := multipart.NewWriter(&b)
	for key, r := range forms {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an file
		if _, ok := r.(*os.File); ok {
			h := make(textproto.MIMEHeader)
			h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s.txt"`, key, time.Now().Format("20060102150405")))
			h.Set("Content-Type", "image/png")
			if fw, err = w.CreatePart(h); err != nil {
				break
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				break
			}
		}
		if _, err := io.Copy(fw, r); err != nil {
			break
		}
	}
	w.Close()

	req := httptest.NewRequest(http.MethodPost, uri, &b)
	req.Header.Add("Content-Type", w.FormDataContentType())
	return getBody(req, headers, router)
}

// PutJSON method
func PutJSON(uri string, param map[string]interface{}, headers map[string]string, router *gin.Engine) ([]byte, error) {
	jsonByte, _ := json.Marshal(param)
	req := httptest.NewRequest(http.MethodPut, uri, bytes.NewReader(jsonByte))
	return getBody(req, headers, router)
}

// DeleteJSON method
func DeleteJSON(uri string, param map[string]interface{}, headers map[string]string, router *gin.Engine) ([]byte, error) {
	jsonByte, _ := json.Marshal(param)
	req := httptest.NewRequest(http.MethodDelete, uri, bytes.NewReader(jsonByte))
	return getBody(req, headers, router)
}

func getBody(req *http.Request, headers map[string]string, router *gin.Engine) ([]byte, error) {
	for key, header := range headers {
		req.Header.Set(key, header)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	log.Println(result.StatusCode)
	if result.StatusCode >= 400 {
		return nil, errors.New(fmt.Sprintf("request error by code: %d", result.StatusCode))
	}

	body, _ := io.ReadAll(result.Body)
	return body, nil
}

type StorageSuite struct {
	suite.Suite
	ctx    context.Context
	cancel func()
}

func (suite *StorageSuite) SetupSuite() {
	suite.ctx, suite.cancel = context.WithTimeout(context.Background(), 10*time.Second)
}

func (suite *StorageSuite) TearDownSuite() {
	suite.cancel()
}

func (suite *StorageSuite) TestRegister() {
	testIFile := &testIFile{}
	storage.Register(testIFile, func() {})
	defer storage.Unload()
	route := gin.New()
	Register(route)
	suite.T().Log(route)
	suite.Equal(DefaultPrefix, route.Routes()[0].Path)
	suite.Equal(http.MethodPost, route.Routes()[0].Method)
	suite.Equal(DefaultPrefix, route.Routes()[1].Path)
	suite.Equal(http.MethodPut, route.Routes()[1].Method)
	suite.Equal(DefaultPrefix, route.Routes()[2].Path)
	suite.Equal(http.MethodDelete, route.Routes()[2].Method)
	suite.Equal(DefaultPrefix, route.Routes()[3].Path)
	suite.Equal(http.MethodGet, route.Routes()[3].Method)
	suite.T().Log(route.Routes()[0].Path)
	suite.T().Log(storage.FILE)
}

func (suite *StorageSuite) TestUpload() {
	type want struct {
		Path        string
		UploadError error
	}

	testCases := []struct {
		Label  string
		Prefix string
		Want   want
	}{
		{
			Label:  "Upload media",
			Prefix: "test",
			Want: want{
				Path: "test/testPath",
			},
		},
		{
			Label:  "Upload media upload fail",
			Prefix: "test",
			Want: want{
				Path:        "test/testPath",
				UploadError: panicerrorhandler.ErrFileUpload,
			},
		},
	}
	for _, tc := range testCases {
		func() {
			f, errOpen := os.Open("./storage/cloud/image.png")
			suite.NoError(errOpen)
			defer f.Close()

			testIFile := &testIFile{}
			testIFile.On("Upload", mock.Anything, tc.Prefix, mock.Anything).Return(tc.Want.Path, tc.Want.UploadError)
			storage.Register(testIFile, func() {})
			defer storage.Unload()
			route := NewMockGinServer()
			Register(route)

			forms := map[string]io.Reader{
				"file":   f,
				"prefix": strings.NewReader(tc.Prefix),
			}

			resp, err := PostFile("/storage", forms, map[string]string{}, route)
			if tc.Want.UploadError == nil {
				var result map[string]interface{}
				suite.NoError(json.Unmarshal(resp, &result))
				suite.Equal(tc.Want.Path, result["path"])
			} else {
				suite.Error(err)
			}
		}()
	}
}

func (suite *StorageSuite) TestPublicize() {
	type want struct {
		URL            string
		PublicizeError error
	}

	testCases := []struct {
		Label string
		Path  string
		Want  want
	}{
		{
			Label: "Publicize media",
			Path:  "test/testPath",
			Want: want{
				URL: "https://storage.googleapis.com/test/testPath",
			},
		},
		{
			Label: "Publicize media error",
			Path:  "test/testPath",
			Want: want{
				PublicizeError: panicerrorhandler.ErrFileUpdate,
			},
		},
	}
	for _, tc := range testCases {
		func() {
			testIFile := &testIFile{}
			testIFile.On("GetURL", mock.Anything, tc.Path).Return(tc.Want.URL, tc.Want.PublicizeError)
			storage.Register(testIFile, func() {})
			defer storage.Unload()
			route := NewMockGinServer()
			Register(route)

			resp, err := PutJSON("/storage", map[string]interface{}{
				"path": tc.Path,
			}, map[string]string{}, route)
			if tc.Want.PublicizeError == nil {
				var result map[string]interface{}
				suite.NoError(json.Unmarshal(resp, &result))
				suite.Equal(tc.Want.URL, result["url"])
			} else {
				suite.Error(err)
			}
		}()
	}
}

func (suite *StorageSuite) TestRemove() {
	type want struct {
		RemoveError error
	}

	testCases := []struct {
		Label string
		Path  string
		Want  want
	}{
		{
			Label: "Remove media",
			Path:  "test/testPath",
			Want:  want{},
		},
		{
			Label: "Remove media error",
			Path:  "test/testPath",
			Want: want{
				RemoveError: panicerrorhandler.ErrFileRemove,
			},
		},
	}
	for _, tc := range testCases {
		func() {
			testIFile := &testIFile{}
			testIFile.On("Remove", mock.Anything, tc.Path).Return(tc.Want.RemoveError)
			storage.Register(testIFile, func() {})
			defer storage.Unload()
			route := NewMockGinServer()
			Register(route)

			resp, err := DeleteJSON("/storage", map[string]interface{}{
				"path": tc.Path,
			}, map[string]string{}, route)
			if tc.Want.RemoveError == nil {
				suite.Equal("ok", string(resp))
			} else {
				suite.Error(err)
			}
		}()
	}
}

func (suite *StorageSuite) TestList() {
	type want struct {
		ListError error
	}

	testCases := []struct {
		Label     string
		Path      string
		Delimiter string
		Prefix    string
		Want      want
	}{
		{
			Label:     "List media",
			Delimiter: "/",
			Prefix:    "/test",
			Want:      want{},
		},
		{
			Label: "List media error",
			Want: want{
				ListError: panicerrorhandler.ErrGetFile,
			},
		},
	}
	for _, tc := range testCases {
		func() {
			q := storage.Query{}
			if tc.Delimiter != "" {
				q = storage.WithFileCloudDelimiter(q, tc.Delimiter)
			}
			if tc.Prefix != "" {
				q = storage.WithFileCloudPrefix(q, tc.Prefix)
			}
			testIFile := &testIFile{}
			testIFile.On("List", mock.Anything, q, mock.Anything).Return(tc.Want.ListError)
			storage.Register(testIFile, func() {})
			defer storage.Unload()
			route := NewMockGinServer()
			Register(route)

			u := &url.Values{}
			if tc.Delimiter != "" {
				u.Set("delimiter", tc.Delimiter)
			}
			if tc.Prefix != "" {
				u.Set("prefix", tc.Prefix)
			}
			resp, err := Get("/storage?"+u.Encode(), map[string]string{}, route)
			if tc.Want.ListError == nil {
				suite.T().Log(string(resp))
			} else {
				suite.Error(err)
			}
		}()
	}
}

func TestStorageSuite(t *testing.T) {
	suite.Run(t, new(StorageSuite))
}
