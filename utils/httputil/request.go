// @Author: Perry
// @Date  : 2020/10/3
// @Desc  : http请求方法,使用fast http封装

package httputil

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

const (
	TMFORM  = "multipart/form-data"
	TJSON   = "application/json"
	TTEXT   = "text/plain"
	TXML    = "text/xml"
	TSTREAM = "application/octet-stream"
)

type Client struct {
	URI         string
	Method      string
	Headers     map[string]string
	ContentType string
	Timeout     time.Duration
	req         *fasthttp.Request
	resp        *fasthttp.Response
}

// RAW 发起常规请求
func (s *Client) RAW(body ...[]byte) ([]byte, error) {
	s.init()
	defer s.close()
	s.setHeaders()
	s.req.SetRequestURI(s.URI)
	s.req.Header.SetMethod(s.Method)
	s.req.Header.SetContentType(s.ContentType)
	if len(body) != 0 && strings.ToUpper(s.Method) != "GET" {
		s.req.SetBody(body[0])
	}
	return s.request()
}

// FILE 发起文件请求
func (s *Client) FILE(fname string, fctt []byte) ([]byte, error) {
	s.init()
	defer s.close()
	s.setHeaders()
	s.req.SetRequestURI(s.URI)
	s.req.Header.SetContentType(TSTREAM)
	s.req.Header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fname))
	s.req.Header.Set("Pragma", "No-cache")
	s.req.Header.Set("Cache-Control", "no-cache")
	s.req.Header.Set("Expires", "0")
	s.req.SetBody(fctt)
	return s.request()
}

type RequestForm struct {
	FieldName    string
	IsFile       bool
	FileName     string
	FieldContent []byte
}

// FORM 上传多键值对以及文件
func (s *Client) FORM(formList []RequestForm) ([]byte, error) {
	var (
		buf        bytes.Buffer
		bodyWriter = multipart.NewWriter(&buf)
	)
	if len(formList) == 0 {
		return nil, errors.New("form request failed,no form value given")
	}
	s.init()
	defer s.close()
	s.setHeaders()
	s.req.SetRequestURI(s.URI)
	s.req.Header.SetMethod(s.Method)
	for _, v := range formList {
		if !v.IsFile {
			bodyWriter.WriteField(v.FieldName, string(v.FieldContent))
			continue
		}
		writer, _ := bodyWriter.CreateFormFile(v.FieldName, v.FileName)
		writer.Write(v.FieldContent)
	}
	defer bodyWriter.Close()
	s.req.Header.SetContentType(bodyWriter.FormDataContentType())
	s.req.SetBody(buf.Bytes())
	return s.request()
}

func (s *Client) request() ([]byte, error) {
	var err error
	if s.Timeout != 0 {
		err = fasthttp.DoTimeout(s.req, s.resp, s.Timeout)
	} else {
		err = fasthttp.Do(s.req, s.resp)
	}
	if err != nil {
		return nil, err
	}
	if s.resp.StatusCode() != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("response data code:%d", s.resp.StatusCode()))
	}
	respBody := s.resp.Body()
	return respBody, nil
}

func (s *Client) init() {
	if s.ContentType == "" {
		s.ContentType = TTEXT
	}
	s.req = fasthttp.AcquireRequest()
	s.resp = fasthttp.AcquireResponse()
}

func (s *Client) close() {
	s.req.ConnectionClose()
	s.resp.ConnectionClose()
}

func (s *Client) setHeaders() {
	if s.Headers != nil && len(s.Headers) != 0 {
		for k, v := range s.Headers {
			s.req.Header.Set(k, v)
		}
	}
}
