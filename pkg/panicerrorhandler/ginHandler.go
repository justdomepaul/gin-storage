package panicerrorhandler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type IGinErrorReport interface {
	IErrorReport
	GinReport(c *gin.Context)
}

func GinPanicErrorHandler(system, prefixMessage string) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if errReportInstance, okIErrorReport := err.(IGinErrorReport); okIErrorReport {
					errReportInstance.SetSystem(system).Report(prefixMessage)
					errReportInstance.GinReport(c)
					return
				}
				if instance, ok := err.(error); ok {
					c.AbortWithError(http.StatusInternalServerError, instance.(error))
					return
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
