// @Author : DAIPENGYUAN
// @File : ssh_run
// @Time : 2020/9/24 15:49 
// @Description : 

package engines

import (
	"strings"
)

// 单次运行一条或多条命令
func SSHRun(target, user, password string, cmds ...string) ([]byte, error) {
	client, err := newSSHTimeoutClient(target, user, password)
	if err != nil {
		return nil, err
	}
	session, err := client.NewSession()
	session.Wait()
	if err != nil {
		closeCLient(client)
		return nil, err
	}
	return session.CombinedOutput(strings.Join(cmds, "\n") + "\n")
}
