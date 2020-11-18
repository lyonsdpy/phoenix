// @Author : DAIPENGYUAN
// @File : config
// @Time : 2020/9/9 10:12 
// @Description : 变量配置,用于控制应用中的各个部分

package conf

import "time"

var (
	// 配置参数
	Port     int
	LogLevel string
	Logpath  string // 日志默认保存目录，默认不在文件保存日志
	Pidpath  string // PID文件的保存目录，默认在当前目录下

	// 相关时间参数
	NetTimeout      = 2 * time.Second        // 建立TCP网络连接的超时时间
	WriteTimeout    = 2 * time.Second        // 网络写入超时时间
	ReadTimeout     = 5 * time.Second        // 独立的SSH读取时间
	ReadSleep       = 200 * time.Millisecond // 协商通过之后,需要等待一个时间,否则对方读取通道可能未开启
	ShutdownTimeout = 5 * time.Second        // 关闭程序时，如果5秒还无法清理所有资源则强制关闭
	// pool连接超时
	PoolTimeout   = 300 * time.Second // 连接池的超时时间
	LockerTimeout = 2 * time.Second
	// locker 连接数锁信息
	SNMPLocker = 1
	VTYLocker  = 4
	// SNMP配置
	SNMPMaxRepitition uint8 = 20
)
