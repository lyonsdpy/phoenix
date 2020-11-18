// @Author : DAIPENGYUAN
// @File : common
// @Time : 2020/10/19 16:58
// @Description : 用于测试的接口，打印相关的请求信息
// 此接口可以测试相应callback的返回值

package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"phoenix/mod"
	"phoenix/utils/common"
	"phoenix/utils/encrypt"
	"phoenix/utils/log"
	"strings"

	"github.com/gin-gonic/gin"
)

var CommonRoute = []mod.Route{
	{Method: "POST", Path: "/common", Handler: Common, Comment: "common测试接口"},
	{Method: "GET", Path: "/common-resp", Handler: CommonResp, Comment: "common测试接口"},
	{Method: "POST", Path: "/encrypt_word", Handler: Encrypt, Comment: "获取密码的加密格式"},
}

// @Tags Common
// @Summery 获取body内容的加密格式
// @Description 通常用户获取密码的加密形式，避免在网络中传输密文
// @Accept  json
// @Produce json
// @Param plaintext body string true "明文信息"
// @Success 200 {string} string "密文信息"
// @Router /encrypt_word [post]
func Encrypt(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	enbody, err := encrypt.DesEncrypt(string(body))
	if err != nil {
		c.Error(err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, enbody)
}

func Common(c *gin.Context) {
	log.Logger.Info("进入Common方法")
	commonHeader(c)    // 解析并打印http-header
	commonMultipart(c) // 解析并打印multipart body
	commonJson(c)      // 解析并打印json
	// 解析multipar/form-data中的数据,包含上传的文件
	commonBody(c)
	c.String(200, "ok")
}

func CommonResp(c *gin.Context) {
	fname := "test.txt"
	fbytes := []byte("this is a test file")
	hd := c.Writer.Header()
	hd.Set("Content-Type", "application/octet-stream")
	hd.Set("test-header", "header-content")
	hd.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fname))
	hd.Set("Pragma", "No-cache")
	hd.Set("Cache-Control", "no-cache")
	hd.Set("Expires", "0")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(fbytes)
}

// 解析http-header
func commonHeader(c *gin.Context) {
	if len(c.Request.Header) != 0 {
		log.Logger.Info(strings.Repeat("##", 10) + " 解析到header，打印header")
		for k, v := range c.Request.Header {
			log.Logger.Infof("key= %s len=%d value= %s", k, len(v), v)
		}
	}
}

// 解析multipar/form-data中的数据,包含上传的文件
func commonMultipart(c *gin.Context) {
	mf, err := c.MultipartForm()
	if mf != nil && err == nil {
		log.Logger.Infof(strings.Repeat("##", 10) + " 解析到multipart，打印multipart内容")
		for fieldName, v := range mf.File {
			log.Logger.Infof("fieldName=%s\nfieldType=file", fieldName)
			for _, v2 := range v {
				f, _ := v2.Open()
				fb, _ := ioutil.ReadAll(f)
				log.Logger.Infof("filename=%s\t\tfilectt=%s", v2.Filename, fb)
			}
		}
		for k, fieldName := range mf.Value {
			log.Logger.Infof("fieldName=%s\nfieldType=value\n"+
				"valueLength=%d,value=%+v", k, len(fieldName), fieldName)
		}
	}
}

// 解析json
func commonJson(c *gin.Context) {
	// 解析content-type为json时解析json数据
	if c.ContentType() == "application/json" {
		log.Logger.Info(strings.Repeat("##", 10) + " 解析到json，打印json内容")
		jdata := new(interface{})
		bd, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Logger.Error(err)
			return
		}
		err = json.Unmarshal(bd, jdata)
		if err != nil {
			log.Logger.Error(err)
			return
		}
		log.Logger.Info(common.IndentJson(jdata))
	}
}

// 解析纯body
func commonBody(c *gin.Context) {
	if c.Request.Body == nil {
		log.Logger.Error(errors.New("no body found"))
		return
	}
	bd, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Logger.Error(err)
		return
	}
	log.Logger.Info(strings.Repeat("##", 10) + " 解析到header，打印body")
	log.Logger.Info(string(bd))
}
