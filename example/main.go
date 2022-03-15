package main

import (
	"github.com/gin-gonic/gin"
	"github.com/justdomepaul/gin-storage"
	_ "github.com/justdomepaul/gin-storage/storage/cloud"
	"github.com/justdomepaul/toolbox/errorhandler"
	"net/http"
)

func main() {
	srv := gin.New()
	srv.MaxMultipartMemory = 8 << 20
	srv.Use(gin.Logger(), errorhandler.GinPanicErrorHandler("system", "gin server error"))
	srv.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	closeFn := gin_storage.Register(srv)
	defer closeFn()

	srv.Run()
}
