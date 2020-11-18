// @Author : DAIPENGYUAN
// @File : snmp
// @Time : 2020/9/28 15:02
// @Description : snmp服务

package service

import (
	"github.com/lyonsdpy/gosnmp"
	"phoenix/conf"
	"phoenix/engines"
	"phoenix/locker"
)

func (s Service) SnmpTypeList() map[gosnmp.Asn1BER]string {
	return engines.SNMPTypeList()
}

func (s Service) SnmpBulkGet(odilist []string) ([]gosnmp.SnmpPDU, error) {
	err := locker.SnmpGet(s.Target)
	if err != nil {
		return nil, err
	}
	defer locker.SnmpRelease(s.Target)
	return engines.SNMPBulkGet(s.Target, s.Community, odilist...)
}

func (s Service) SnmpWalk(oid string) ([]gosnmp.SnmpPDU, error) {
	err := locker.SnmpGet(s.Target)
	if err != nil {
		return nil, err
	}
	defer locker.SnmpRelease(s.Target)
	return engines.SNMPWalk(s.Target, s.Community, conf.SNMPMaxRepitition, oid)
}
