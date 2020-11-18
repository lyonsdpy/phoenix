// @Author : DAIPENGYUAN
// @File : logger_gin
// @Time : 2020/10/21 14:47 
// @Description : gin框架插件

package log

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"time"
)

func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		url := c.Request.URL
		clientIP := c.Request.RemoteAddr // 请求IP
		method := c.Request.Method       // 请求方法
		statusCode := c.Writer.Status()  // 请求返回code
		latency := time.Since(start)
		if statusCode >= 400 {
			Logger.Infof("[%s] %3d| %4v | %9s | %s", method, statusCode, latency, clientIP, url)
			body, _ := ioutil.ReadAll(c.Request.Body)
			Logger.Debugf("%s\n", body)
		} else {
			Logger.Infof("[%s] %3d| %4v | %9s | %s", method, statusCode, latency, clientIP, url)
		}
		if len(c.Errors) > 0 {
			for _, v := range c.Errors {
				Logger.Errorf("[%s] %3d| %4v | %9s | %s | error=%s",
					method, statusCode, latency, clientIP, url, v.Error())
			}
		}
	}
}
