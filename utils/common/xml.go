// @Author: Perry
// @Date  : 2020/10/11
// @Desc  : 

package common

import (
	"encoding/xml"
)

func IndentXml(v interface{}) string {
	r, err := xml.MarshalIndent(v, "", "\t")
	if err != nil {
		return IndentXml(err)
	}
	return string(r)
}

func IsValidXML(data []byte) bool {
	return xml.Unmarshal(data, new(interface{})) == nil
}
