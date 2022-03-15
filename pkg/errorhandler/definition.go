package errorhandler

import (
	"fmt"
	"github.com/cockroachdb/errors"
)

var (
	ErrFailCloseSession  = errors.New("fail to close connection")
	ErrDriveNotExist     = errors.New("file drive not exist")
	ErrFileNotExist      = errors.New("file not exist")
	ErrFileUpload        = errors.New("fail to upload file")
	ErrFailGenerateUUID  = fmt.Errorf("%w: fail to generate uuid", ErrFileUpload)
	ErrFileUpdate        = errors.New("fail to update file")
	ErrFileRemove        = errors.New("fail to remove file")
	ErrGetFile           = errors.New("fail to get file")
	ErrInitialFileClient = errors.New("fail to initial file client")
)
