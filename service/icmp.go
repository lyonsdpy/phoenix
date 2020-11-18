// @Author : DAIPENGYUAN
// @File : icmp
// @Time : 2020/9/28 13:44
// @Description : ping服务

package service

import (
	"phoenix/engines"
	"phoenix/mod"
	"time"
)

func (s Service) Ping(size, count int, maxrtt int64, useudp bool) (*mod.ICMPResponse, error) {
	ng := &engines.Pinger{
		Addr:   s.Target,
		Size:   size,
		Count:  count,
		MaxRTT: time.Duration(maxrtt) * time.Millisecond,
		UseUdp: useudp,
	}
	r, err := ng.Ping()
	if err != nil {
		return nil, err
	}
	var rtts []string
	for _, v := range r.Rtts {
		rtts = append(rtts, v.String())
	}
	rst := &mod.ICMPResponse{
		PacketsRecv:  r.PacketsRecv,
		PacketsSent:  r.PacketsSent,
		PacketLoss:   r.PacketLoss,
		RemoteIpAddr: r.IPAddr.String(),
		RemoteAddr:   r.Addr,
		Rtts:         rtts,
		MinRtt:       r.MinRtt.String(),
		MaxRtt:       r.MaxRtt.String(),
		AvgRtt:       r.AvgRtt.String(),
		StdDevRtt:    r.StdDevRtt.String(),
	}
	return rst, nil
}
