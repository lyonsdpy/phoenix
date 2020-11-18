// @Author : DAIPENGYUAN
// @File : handler
// @Time : 2020/10/26 11:52 
// @Description : 

package mod

import "github.com/gin-gonic/gin"

type Route struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
	Comment string
}