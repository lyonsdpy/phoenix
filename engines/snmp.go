// @Author : DAIPENGYUAN
// @File : snmp
// @Time : 2020/9/26 14:50 
// @Description : 

package engines

import (
	"github.com/lyonsdpy/gosnmp"
	"phoenix/conf"
)

func SNMPTypeList() map[gosnmp.Asn1BER]string {
	return DataTypeStrings
}

func SNMPBulkGet(target, community string, oids ...string) (pdus []gosnmp.SnmpPDU, err error) {
	ng, err := gosnmp.NewGoSNMP(target, community, gosnmp.Version2c, int64(conf.NetTimeout))
	if err != nil {
		return
	}
	defer ng.Close()
	pkg, err := ng.GetBulk(uint8(len(oids)), 1, oids...)
	if err != nil {
		return
	}
	return pkg.Variables, nil
}

func SNMPWalk(target, community string, maxRepetition uint8, oid string) (pdus []gosnmp.SnmpPDU, err error) {
	ng, err := gosnmp.NewGoSNMP(target, community, gosnmp.Version2c, int64(conf.NetTimeout))
	if err != nil {
		return
	}
	defer ng.Close()
	pdus, err = ng.BulkWalk(maxRepetition, oid)
	return
}

const (
	Integer          gosnmp.Asn1BER = 0x02
	BitString                       = 0x03
	OctetString                     = 0x04
	Null                            = 0x05
	ObjectIdentifier                = 0x06
	Sequence                        = 0x30
	IpAddress                       = 0x40
	Counter32                       = 0x41
	Gauge32                         = 0x42
	TimeTicks                       = 0x43
	Opaque                          = 0x44
	NsapAddress                     = 0x45
	Counter64                       = 0x46
	Uinteger32                      = 0x47
	NoSuchObject                    = 0x80
	NoSuchInstance                  = 0x81
	GetRequest                      = 0xa0
	GetNextRequest                  = 0xa1
	GetResponse                     = 0xa2
	SetRequest                      = 0xa3
	Trap                            = 0xa4
	GetBulkRequest                  = 0xa5
	EndOfMibView                    = 0x82
)

// String representations of each SNMP Data Type
var DataTypeStrings = map[gosnmp.Asn1BER]string{
	Integer:          "Integer",
	BitString:        "BitString",
	OctetString:      "OctetString",
	Null:             "Null",
	ObjectIdentifier: "ObjectIdentifier",
	IpAddress:        "IpAddress",
	Sequence:         "Sequence",
	Counter32:        "Counter32",
	Gauge32:          "Gauge32",
	TimeTicks:        "TimeTicks",
	Opaque:           "Opaque",
	NsapAddress:      "NsapAddress",
	Counter64:        "Counter64",
	Uinteger32:       "Uinteger32",
	NoSuchObject:     "NoSuchObject",
	NoSuchInstance:   "NoSuchInstance",
	GetRequest:       "GetRequest",
	GetNextRequest:   "GetNextRequest",
	GetResponse:      "GetResponse",
	SetRequest:       "SetRequest",
	Trap:             "Trap",
	GetBulkRequest:   "GetBulkRequest",
	EndOfMibView:     "endOfMib",
}
