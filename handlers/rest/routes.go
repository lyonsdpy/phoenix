// @Author : DAIPENGYUAN
// @File : routes
// @Time : 2020/9/30 11:09 
// @Description : 

package rest

import (
	_ "phoenix/docs"
	"phoenix/mod"
	"phoenix/utils/log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var apiVer = "v1"

func SetupRoutes(ng *gin.Engine) {
	// 注册路由
	routeList := makeRouteList(CommonRoute, IcmpRoute, ScpRoute,
		SftpRoute, SSHRunRoute, SSHShellRoute, SnmpRoute,
		SnetconfRoute, JobRoute)

	// 将汇总的路由注册到GIN中
	rg := ng.Group(apiVer)
	rg.Use(log.GinMiddleware())
	for _, v := range routeList {
		rg.Handle(v.Method, v.Path, v.Handler)
	}

	ng.GET("/swagger/*any",
		ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "NAME_OF_ENV_VARIABLE"))

	ng.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
}

func makeRouteList(lists ...[]mod.Route) []mod.Route {
	var routeList []mod.Route
	for _, v := range lists {
		routeList = append(routeList, v...)
	}
	return routeList
}
