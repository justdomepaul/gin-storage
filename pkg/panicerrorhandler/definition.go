package panicerrorhandler

import (
	"fmt"
	"github.com/cockroachdb/errors"
)

const (
	ErrProcessAuthenticate    = "errAuthenticate"
	ErrDbAlreadyExists        = "errDBAlreadyExists"
	ErrDbConnection           = "errDBConnection"
	ErrDbDisconnection        = "errDBDisConnection"
	ErrDbExecute              = "errDBExecute"
	ErrDbRowNotFound          = "errDBRowNotFound"
	ErrDbUpdateNoEffect       = "errDBUpdateNoEffect"
	ErrProcessExecute         = "errExecute"
	ErrGrpcConnection         = "errGRPCConnection"
	ErrGrpcExecute            = "errGRPCExecute"
	ErrProcessInvalidArgument = "errProcessInvalidArgument"
	ErrJsonMarshal            = "errJSONMarshal"
	ErrJsonUnmarshal          = "errJSONUnmarshal"
	ErrJwtExecute             = "errJWTExecute"
	ErrDataNotFound           = "errDataNotFound"
	ErrProcessPermissionDeny  = "errPermissionDeny"
	ErrProcessServerExecute   = "errServerExecute"
	ErrProcessVariable        = "errVariable"
)

var (
	ErrNoRows            = errors.New("no rows in result set")
	ErrUpdateNoEffect    = errors.New("no rows effected")
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
