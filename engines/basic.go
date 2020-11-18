// @Author : DAIPENGYUAN
// @File : basic
// @Time : 2020/9/11 15:36 
// @Description : 

package engines

import (
	"bytes"
	"errors"
	"io"
	"phoenix/conf"
	"sync"
	"time"
)

const msgSeperator = " ##  ## " // cli执行命令时的分隔符

func cliWrite(t io.Writer, data []byte) {
	t.Write(data)
	if (len(data)+len(msgSeperator))%4096 < 6 {
		t.Write([]byte("      "))
	}
	t.Write([]byte(msgSeperator))
	t.Write([]byte("\n"))
	return
}

// 交互式命令行读取信息时使用的方法
// 读取到分隔符退出
// io读取超时退出
func cliRead(rd io.Reader) ([]byte, error) {
	var (
		out bytes.Buffer
		err error
		n   int
		sep = []byte(msgSeperator)
		buf = make([]byte, 4096)
		pos = 0
	)
	for {
		n, err = rd.Read(buf[pos : pos+(len(buf)/2)])
		if n > 0 {
			if end := bytes.Index(buf[0:pos+n], sep); end > -1 {
				out.Write(buf[0:end])
				break
			}
			out.Write(buf[pos : pos+n])
			if pos > 0 {
				copy(buf, buf[pos:pos+n])
			}
			pos = n
		}
		if err != nil {
			break
		}
	}
	return out.Bytes(), err
}

// 第二种read，暂不删除
//func cliRead(rd io.Reader) ([]byte, error) {
//	var (
//		out []byte
//		err error
//		n   int
//		sep = []byte(msgSeperator)
//		buf = make([]byte, 4096)
//	)
//	for {
//		n, err = rd.Read(buf)
//		if n > 0 {
//			out = append(out, buf[0:n]...)
//			if end := bytes.Index(out, sep); end > -1 {
//				out = out[0:end]
//				break
//			}
//		}
//		if err != nil {
//			break
//		}
//	}
//	return out, err
//}

//// 带超时时间的wg等待
//// NOTE:此方法为上层读取超时等待,优先使用
func waitTimeout(wg *sync.WaitGroup) error {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		wg.Wait()
	}()
	select {
	case <-ch:
		return nil
	case <-time.After(conf.ReadTimeout):
		return errors.New("read failed,io timeout")
	}
}
