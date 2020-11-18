// @Author: Perry
// @Date  : 2020/10/11
// @Desc  : json相关工具

package common

import "encoding/json"

func IndentJson(obj interface{}) string {
	ret, err := json.MarshalIndent(obj, "", "\t")
	if err != nil {
		return err.Error()
	}
	return string(ret)
}
