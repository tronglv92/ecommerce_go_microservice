package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tronglv92/accounts/common"
)

// func Recover(ac appctx.AppContext) gin.HandlerFunc {
func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Header("Content-Type", "application/json")
				if appErr, ok := err.(*common.AppError); ok {
					c.AbortWithStatusJSON(appErr.StatusCode, appErr)
					panic(err)
					return
				}

				fmt.Printf("Recover err %v", err)
				appErr := common.ErrInternal(err.(error))
				c.AbortWithStatusJSON(appErr.StatusCode, appErr)
				panic(err)
				return
			}
		}()

		c.Next()
	}
}
