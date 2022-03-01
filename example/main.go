package main

import (
	"github.com/gin-gonic/gin"
	"github.com/justdomepaul/gin-storage"
	"github.com/justdomepaul/gin-storage/pkg/panicerrorhandler"
	_ "github.com/justdomepaul/gin-storage/storage/cloud"
	"net/http"
)

func main() {
	srv := gin.New()
	srv.MaxMultipartMemory = 8 << 20
	srv.Use(gin.Logger(), panicerrorhandler.GinPanicErrorHandler("system", "gin server error"))
	srv.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	closeFn := gin_storage.Register(srv)
	defer closeFn()

	srv.Run()
}
