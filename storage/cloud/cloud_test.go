package cloud

import (
	gs "cloud.google.com/go/storage"
	"context"
	"github.com/google/uuid"
	"github.com/justdomepaul/gin-storage/pkg/config"
	"github.com/justdomepaul/gin-storage/pkg/errorhandler"
	"github.com/justdomepaul/gin-storage/storage"
	"github.com/stretchr/testify/suite"
	"google.golang.org/api/option"
	"os"
	"strings"
	"testing"
	"time"
)

type CloudSuite struct {
	suite.Suite
	ctx    context.Context
	cancel func()
	client *gs.Client
}

func (suite *CloudSuite) SetupSuite() {
	suite.ctx, suite.cancel = context.WithTimeout(context.Background(), 10*time.Second)
}

func (suite *CloudSuite) TearDownSuite() {
	suite.cancel()
}

func (suite *CloudSuite) SetupTest() {
	suite.NoError(os.Setenv("STORAGE_EMULATOR_HOST", "localhost:9023"))
	var opts []option.ClientOption
	c, err := gs.NewClient(suite.ctx, opts...)
	suite.NoError(err)
	suite.client = c
}

func (suite *CloudSuite) TestUploadMethod() {
	f, err := os.Open("./image.png")
	suite.NoError(err)
	defer f.Close()

	type want struct {
		PathPrefix string
		Match      bool
	}

	testCases := []struct {
		Label  string
		Media  config.Media
		Prefix string
		Want   want
	}{
		{
			Label: "Upload media into root path",
			Media: config.Media{
				BucketName: "staging.megaphone.appspot.com",
				PrefixPath: "/",
			},
			Prefix: "",
			Want: want{
				PathPrefix: "/",
				Match:      true,
			},
		},
		{
			Label: "Upload media into prefix path",
			Media: config.Media{
				BucketName: "staging.megaphone.appspot.com",
				PrefixPath: "/media/",
			},
			Prefix: "",
			Want: want{
				PathPrefix: "/media/",
				Match:      true,
			},
		},
		{
			Label: "Upload media into root path and sub prefix",
			Media: config.Media{
				BucketName: "staging.megaphone.appspot.com",
				PrefixPath: "/",
			},
			Prefix: "sub",
			Want: want{
				PathPrefix: "/sub/",
				Match:      true,
			},
		},
		{
			Label: "Upload media into prefix path and sub prefix",
			Media: config.Media{
				BucketName: "staging.megaphone.appspot.com",
				PrefixPath: "/media/",
			},
			Prefix: "sub",
			Want: want{
				PathPrefix: "/media/sub/",
				Match:      true,
			},
		},
	}

	for _, tc := range testCases {
		result, err := NewFile(tc.Media, suite.client).Upload(suite.ctx, tc.Prefix, f)
		suite.NoError(err)
		suite.Equal(tc.Want.Match, strings.HasPrefix(result, tc.Want.PathPrefix))
	}
}

func (suite *CloudSuite) TestGetURLMethod() {
	f, err := os.Open("./image.png")
	suite.NoError(err)
	defer f.Close()

	type want struct {
		PathPrefix string
		Match      bool
		Error      error
	}

	testCases := []struct {
		Label     string
		Media     config.Media
		Prefix    string
		ErrorPath string
		Want      want
	}{
		{
			Label: "Get media Public URL",
			Media: config.Media{
				BucketName: "staging.megaphone.appspot.com",
				PrefixPath: "/media/",
			},
			Prefix: "sub",
			Want: want{
				PathPrefix: "https://storage.googleapis.com/staging.megaphone.appspot.com/media/sub/",
				Match:      true,
			},
		},
		{
			Label: "Get media Public URL path format error",
			Media: config.Media{
				BucketName: "staging.megaphone.appspot.com",
				PrefixPath: "/media/",
			},
			Prefix:    "sub",
			ErrorPath: "/media/sub/",
			Want: want{
				Error: errorhandler.ErrFileUpdate,
			},
		},
		{
			Label: "Get media Public URL file not exist",
			Media: config.Media{
				BucketName: "staging.megaphone.appspot.com",
				PrefixPath: "/media/",
			},
			Prefix:    "sub",
			ErrorPath: "/media/sub",
			Want: want{
				Error: errorhandler.ErrFileNotExist,
			},
		},
	}

	for _, tc := range testCases {
		file := NewFile(tc.Media, suite.client)
		result := tc.ErrorPath
		if tc.Want.Error == nil {
			result, err = file.Upload(suite.ctx, tc.Prefix, f)
			suite.NoError(err)
		}
		url, err := file.GetURL(suite.ctx, result)
		if tc.Want.Error != nil {
			suite.ErrorIs(err, tc.Want.Error)
		} else {
			suite.NoError(err)
			suite.Equal(tc.Want.Match, strings.HasPrefix(url, tc.Want.PathPrefix))
		}
	}
}

func (suite *CloudSuite) TestRemoveMethod() {
	f, err := os.Open("./image.png")
	suite.NoError(err)
	defer f.Close()

	type want struct {
		Error error
	}

	testCases := []struct {
		Label     string
		Media     config.Media
		Prefix    string
		ErrorPath string
		Want      want
	}{
		{
			Label: "Remove media",
			Media: config.Media{
				BucketName: "staging.megaphone.appspot.com",
				PrefixPath: "/media/",
			},
			Prefix: "sub",
			Want:   want{},
		},
		{
			Label: "Remove media path format fail",
			Media: config.Media{
				BucketName: "staging.megaphone.appspot.com",
				PrefixPath: "/media/",
			},
			Prefix:    "sub",
			ErrorPath: "/media/sub/",
			Want: want{
				Error: errorhandler.ErrFileRemove,
			},
		},
		{
			Label: "Remove Not Exist media",
			Media: config.Media{
				BucketName: "staging.megaphone.appspot.com",
				PrefixPath: "/media/",
			},
			Prefix:    "sub",
			ErrorPath: "/media/sub",
			Want: want{
				Error: errorhandler.ErrFileRemove,
			},
		},
	}

	for _, tc := range testCases {
		file := NewFile(tc.Media, suite.client)
		result := tc.ErrorPath
		if tc.Want.Error == nil {
			result, err = file.Upload(suite.ctx, tc.Prefix, f)
			suite.NoError(err)
		}
		if tc.Want.Error != nil {
			suite.ErrorIs(file.Remove(suite.ctx, result), tc.Want.Error)
		} else {
			suite.NoError(file.Remove(suite.ctx, result))
		}
	}
}

func (suite *CloudSuite) TestListMethod() {
	type want struct {
		FolderNames []string
		FolderPaths []string
		QueryPrefix string
		URLPrefix   string
		Error       error
	}

	testCases := []struct {
		Label       string
		Media       config.Media
		QueryPrefix string
		Delimiter   string
		ErrorPath   string
		Want        want
	}{
		{
			Label: "List media folders & files",
			Media: config.Media{
				BucketName: "staging.megaphone.appspot.com",
			},
			QueryPrefix: "/media/sub/",
			Delimiter:   "/",
			Want: want{
				FolderNames: []string{"child", "master", "root"},
				FolderPaths: []string{"/media/sub/child/", "/media/sub/master/", "/media/sub/root/"},
				QueryPrefix: "/media/sub/",
				URLPrefix:   "https://storage.googleapis.com/staging.megaphone.appspot.com/media/sub/",
			},
		},
		{
			Label: "List media files",
			Media: config.Media{
				BucketName: "staging.megaphone.appspot.com",
			},
			QueryPrefix: "/media/sub/child",
			Want: want{
				QueryPrefix: "/media/sub/child/",
				URLPrefix:   "https://storage.googleapis.com/staging.megaphone.appspot.com/media/sub/child/",
			},
		},
		{
			Label: "List audio folders & files",
			Media: config.Media{
				BucketName: "staging.megaphone.appspot.com",
			},
			QueryPrefix: "/audio/sub",
			Delimiter:   "/",
			Want: want{
				FolderNames: []string{"child"},
				FolderPaths: []string{"/audio/sub/child/"},
				QueryPrefix: "/audio/sub/child/",
				URLPrefix:   "https://storage.googleapis.com/staging.megaphone.appspot.com/audio/sub/child/",
			},
		},
	}

	seedsUpload := []struct {
		PrefixPath string
		SubPath    string
	}{
		{
			PrefixPath: "/media/",
			SubPath:    "",
		},
		{
			PrefixPath: "/media/",
			SubPath:    "",
		},
		{
			PrefixPath: "/media/",
			SubPath:    "sub",
		},
		{
			PrefixPath: "/media/",
			SubPath:    "sub/child",
		},
		{
			PrefixPath: "/media/",
			SubPath:    "sub/master",
		},
		{
			PrefixPath: "/media/",
			SubPath:    "sub/root",
		},
		{
			PrefixPath: "/audio/",
			SubPath:    "",
		},
		{
			PrefixPath: "/audio/",
			SubPath:    "sub/child",
		},
	}

	for _, item := range seedsUpload {
		f, err := os.Open("./image.png")
		suite.NoError(err)
		result, err := NewFile(config.Media{
			BucketName: "staging.megaphone.appspot.com",
			PrefixPath: item.PrefixPath,
		}, suite.client).Upload(suite.ctx, item.SubPath, f)
		suite.NoError(err)
		suite.NotNil(result)
		suite.NoError(f.Close())
	}

	for _, tc := range testCases {
		file := NewFile(tc.Media, suite.client)
		var fs []storage.File
		q := storage.Query{}
		if tc.QueryPrefix != "" {
			q = storage.WithFileCloudPrefix(q, tc.QueryPrefix)
		}
		if tc.Delimiter != "" {
			q = storage.WithFileCloudDelimiter(q, tc.Delimiter)
		}
		suite.NoError(file.List(suite.ctx, q, func(file storage.File) error {
			fs = append(fs, file)
			return nil
		}))
		for _, item := range fs {
			suite.T().Logf("%+v", item)
			if folderName, folderPath, exist := item.FolderInfo(); exist {
				suite.Contains(tc.Want.FolderNames, folderName)
				suite.Contains(tc.Want.FolderPaths, folderPath)
				continue
			}
			suite.True(strings.HasPrefix(item.Path(), tc.Want.QueryPrefix))
			suite.True(strings.HasPrefix(item.GetURL(), tc.Want.URLPrefix))
			uid, err := uuid.Parse(item.Name())
			suite.NoError(err)
			suite.NotNil(uid)
			folderName, folderPath, exist := item.FolderInfo()
			if !exist {
				suite.T().Log(item.Size())
				suite.T().Log(item.CreatedTime())
				suite.T().Log(item.ModTime())
				suite.T().Log(item.ModTime())
			} else {
				suite.T().Log(folderName)
				suite.T().Log(folderPath)
			}
		}
	}
}

func TestCloudSuite(t *testing.T) {
	suite.Run(t, new(CloudSuite))
}
