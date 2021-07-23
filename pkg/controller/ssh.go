// Package controller provides SSH control method
package controller

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
)

// 创建SSH客户端。创建出的客户端需要用户自己主动关闭。
func NewSSHClient(addr, user, password string) (*ssh.Client, error) {
	var authMethods []ssh.AuthMethod

	keyboardInteractiveChallenge := func(
		user,
		instruction string,
		questions []string,
		echos []bool,
	) (answers []string, err error) {
		if len(questions) == 0 {
			return []string{}, nil
		}
		return []string{password}, nil
	}

	authMethods = append(authMethods, ssh.KeyboardInteractive(keyboardInteractiveChallenge))
	authMethods = append(authMethods, ssh.Password(password))
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: authMethods,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	return ssh.Dial("tcp", addr,  sshConfig)
}

// SSH操作对象结构体
type SSH struct {
	client *ssh.Client
	session *ssh.Session
}

func (s *SSH) Init(addr, user, password string) error {
	var authMethods []ssh.AuthMethod

	keyboardInteractiveChallenge := func(
		user,
		instruction string,
		questions []string,
		echos []bool,
	) (answers []string, err error) {
		if len(questions) == 0 {
			return []string{}, nil
		}
		return []string{password}, nil
	}

	authMethods = append(authMethods, ssh.KeyboardInteractive(keyboardInteractiveChallenge))
	authMethods = append(authMethods, ssh.Password(password))
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: authMethods,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	var err error
	s.client, err = ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		fmt.Printf("\033[0;32;31m");
		fmt.Printf("%s\n", err)
		fmt.Printf("\033[m")
		return err
	}

	s.session, err = s.client.NewSession()
	if err != nil {
		fmt.Printf("\033[0;32;31m")
		fmt.Printf("%s\n", err)
		fmt.Printf("\033[m")
		return err
	}
	return err
}

// 关闭
func (s *SSH) Close() {
	if s.session != nil {
		s.session.Close()
	}
	if s.client != nil {
		s.client.Close()
	}
}

// 获取SSH客户端
func (s *SSH) GetClient() *ssh.Client {
	return s.client
}

// 远程执行命令
func (s *SSH) Execute(opm string) error {
	reply, err := s.session.Output(opm)
	if err != nil {
		if err.Error() == "Process exited with status 255" {
			fmt.Printf("[\033[0;32;31mERROR\033[m] command `%s` not exists\n", opm)
		} else if err.Error() == "Process exited with status 1" {
			fmt.Printf("[\033[0;32;31mERROR\033[m] command `%s` return 1\n", opm)
		} else {
			fmt.Printf("[\033[0;32;31mERROR\033[m] %s\n", err)
		}
		return err
	} else {
		fmt.Printf("%s\n", reply)
	}
	return nil
}