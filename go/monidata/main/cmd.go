package main

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/confd/log"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

//远程连接
func Connect(host, user, password string, port int) (*ssh.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		err          error
	)
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
	}
	addr = fmt.Sprintf("%s:%d", host, port)
	client, err = ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, err
	}
	return client, nil
}

//远程连接SFTP
func ConnectSFTP(host, user, password string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
	}

	addr = fmt.Sprintf("%s:%d", host, port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}

//远程执行命令
func RunCmd(client *ssh.Client, cmd string) (string, error) {

	log.Info("准备执行命令%s", cmd)

	session, err := client.NewSession()
	if err != nil {
		log.Error("创建session异常,原因:%v", err)
		return "", err
	}
	defer session.Close()

	bytes, err := session.Output(cmd)
	if err != nil {
		log.Error("执行命令%s失败,异常:%v", cmd, err)
		return "", err
	}

	log.Info("执行命令%s成功,返回结果:%s", cmd, string(bytes))

	return string(bytes), nil
}
