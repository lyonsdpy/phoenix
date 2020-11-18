// @Author : DAIPENGYUAN
// @File : icmp
// @Time : 2020/9/8 11:29
// @Description : icmp引擎

package engines

import (
	"github.com/tatsushid/go-fastping"
	"math"
	"net"
	"time"
)

type Pinger struct {
	Addr   string // 目标ip或域名
	Size   int
	Count  int
	MaxRTT time.Duration
	UseUdp bool
	engine *fastping.Pinger
	ipAddr *net.IPAddr
}

func (p *Pinger) Ping() (*PingStatistics, error) {
	addr, err := p.checkArgs()
	if err != nil {
		return nil, err
	}
	eg := fastping.NewPinger()
	eg.Size = p.Size
	eg.MaxRTT = p.MaxRTT
	eg.AddIPAddr(addr)
	if p.UseUdp {
		eg.Network("udp")
	}
	p.ipAddr = addr
	p.engine = eg
	return p.exec()
}

func (p *Pinger) exec() (*PingStatistics, error) {
	var rst = new(PingStatistics)
	rst.Addr = p.Addr
	rst.IPAddr = *p.ipAddr
	p.engine.OnRecv = func(addr *net.IPAddr, duration time.Duration) {
		rst.PacketsRecv++
		rst.Rtts = append(rst.Rtts, duration)
	}
	for i := 0; i < p.Count; i++ {
		p.engine.Run()
		rst.PacketsSent++
		if err := p.engine.Err(); err != nil {
			return nil, err
		}
		p.engine.Stop()
	}
	// 计算ping结果
	rst = p.countResult(rst)
	return rst, nil
}

func (p *Pinger) countResult(rst *PingStatistics) *PingStatistics {
	var min, max, total time.Duration
	rst.PacketLoss = float64(rst.PacketsSent-rst.PacketsRecv) / float64(rst.PacketsSent) * 100
	// 计算延迟
	if len(rst.Rtts) > 0 {
		min = rst.Rtts[0]
		max = rst.Rtts[0]
	}
	for _, rtt := range rst.Rtts {
		if rtt < min {
			min = rtt
		}
		if rtt > max {
			max = rtt
		}
		total += rtt
	}
	rst.MinRtt = min
	rst.MaxRtt = max
	if len(rst.Rtts) > 0 {
		rst.AvgRtt = total / time.Duration(len(rst.Rtts))
		var sumsquares time.Duration
		for _, rtt := range rst.Rtts {
			sumsquares += (rtt - rst.AvgRtt) * (rtt - rst.AvgRtt)
		}
		rst.StdDevRtt = time.Duration(math.Sqrt(
			float64(sumsquares / time.Duration(len(rst.Rtts)))))
	}
	return rst
}

func (p *Pinger) checkArgs() (*net.IPAddr, error) {
	ra, err := net.ResolveIPAddr("ip4:icmp", p.Addr)
	if err != nil {
		return nil, err
	}
	if p.Size == 0 {
		p.Size = 64
	}
	if p.Count == 0 {
		p.Count = 1
	}
	if p.MaxRTT == 0 {
		p.MaxRTT = 1 * time.Second
	}
	return ra, nil
}

type PingStatistics struct {
	PacketsRecv int             `json:"packets_recv"`
	PacketsSent int             `json:"packets_sent"`
	PacketLoss  float64         `json:"packet_loss"`
	IPAddr      net.IPAddr      `json:"ip_addr"` //返回地址
	Addr        string          `json:"addr"`    // 目标地址
	Rtts        []time.Duration `json:"rtts"`
	MinRtt      time.Duration   `json:"min_rtt"`
	MaxRtt      time.Duration   `json:"max_rtt"`
	AvgRtt      time.Duration   `json:"avg_rtt"`
	StdDevRtt   time.Duration   `json:"std_dev_rtt"`
}

type OneRecv struct {
	Addr string
	Rtt  int64
}
