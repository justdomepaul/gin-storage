package storage

import "github.com/justdomepaul/gin-storage/pkg/errorhandler"

var (
	FILE    IFile
	CLOSEFN func()
)

func Register(input IFile, closeFn func()) {
	FILE = input
	CLOSEFN = closeFn
}

func Load() (IFile, func()) {
	if FILE == nil || CLOSEFN == nil {
		panic(errorhandler.ErrDriveNotExist)
	}
	return FILE, CLOSEFN
}

func Unload() {
	FILE = nil
	CLOSEFN = nil
}
