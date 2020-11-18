// @Author : DAIPENGYUAN
// @File : ssh_basic
// @Time : 2020/9/9 9:50 
// @Description : ssh基础配置

package engines

import (
	"golang.org/x/crypto/ssh"
	"phoenix/conf"
	"phoenix/utils/netutil"
)

// 获取正常的SSH客户端
func newSSHClient(target, user, password string) (*ssh.Client, error) {
	sshCfg := getSSHConfig(user, password)
	conn, err := netutil.DialConn(target, "tcp", conf.NetTimeout)
	if err != nil {
		return nil, err
	}
	c, chans, reqs, err := ssh.NewClientConn(conn, target, sshCfg)
	if err != nil {
		return nil, err
	}
	client := ssh.NewClient(c, chans, reqs)
	return client, nil
}

// 获取带超时连接的SSH客户端
func newSSHTimeoutClient(target, user, password string) (*ssh.Client, error) {
	sshCfg := getSSHConfig(user, password)
	conn, err := netutil.DialTimeoutConn(target, "tcp",
		conf.NetTimeout, conf.ReadTimeout, conf.WriteTimeout)
	if err != nil {
		return nil, err
	}
	c, chans, reqs, err := ssh.NewClientConn(conn, target, sshCfg)
	if err != nil {
		return nil, err
	}
	client := ssh.NewClient(c, chans, reqs)
	return client, nil
}

// 获取ssh配置
func getSSHConfig(user, password string) *ssh.ClientConfig {
	var (
		sshConfig ssh.Config
	)
	sshConfig.SetDefaults()
	sshConfig.Ciphers = append(sshConfig.Ciphers, "aes128-ctr", "aes192-ctr", "aes256-ctr",
		"aes128-cbc", "aes256-cbc", "3des-cbc", "des-cbc")
	clientConfig := &ssh.ClientConfig{
		Config:          sshConfig,
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         conf.NetTimeout,
	}
	return clientConfig
}

// 关闭ssh连接
func closeCLient(client *ssh.Client) {
	if client != nil {
		client.Close()
	}
	return
}
