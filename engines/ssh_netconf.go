// @Author : DAIPENGYUAN
// @File : snetconf
// @Time : 2020/9/8 11:28 
// @Description : netconf引擎

package engines

import (
	"github.com/lyonsdpy/go-netconf/netconf"
)

func NewNetconf(target, user, password string, capabilities ...string) (*SSHNetconf, error) {
	client, err := newSSHClient(target, user, password)
	if err != nil {
		return nil, err
	}
	s, err := netconf.NewSSHSession2(client, capabilities...)
	if err != nil {
		closeCLient(client)
		return nil, err
	}
	return &SSHNetconf{netconf: s}, nil
}

type SSHNetconf struct {
	netconf *netconf.Session
}

func (s *SSHNetconf) Exec(xmlbyte string) (ret *netconf.RPCReply, err error) {
	ret, err = s.netconf.Exec(NCRule(xmlbyte))
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *SSHNetconf) Capabilities() []string {
	if s.netconf != nil {
		return s.netconf.ServerCapabilities
	} else {
		return nil
	}
}

func (s *SSHNetconf) Close() error {
	var err error
	if s.netconf != nil {
		if _, err = s.netconf.Exec(NCRule("<close-session/>")); err != nil {
			return err
		}
		if err = s.netconf.Close(); err != nil {
			return err
		}
	}
	return nil
}

type NCRule string

func (r NCRule) MarshalMethod() string {
	return string(r)
}
