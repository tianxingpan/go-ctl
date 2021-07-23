// Package cmd provides add custom CTL command
package cmd

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

func RemoteExecute(ipPort, user, password, opm string) {
	var authMethods []ssh.AuthMethod

	fmt.Printf("\033[1;33m")
	fmt.Printf("[%s]\n", ipPort)
	fmt.Printf("\033[m")

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

	client, err := ssh.Dial("tcp", ipPort, sshConfig)
	if err != nil {
		fmt.Printf("\033[0;32;31m")
		fmt.Printf("%s\n", err)
		fmt.Printf("\033[m")
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("\033[0;32;31m")
		fmt.Printf("%s\n", err)
		fmt.Printf("\033[m")
		return
	}
	defer session.Close()

	reply, err := session.Output(opm)
	if err != nil {
		if err.Error() == "Process exited with status 255" {
			fmt.Printf("[\033[0;32;31mERROR\033[m] command `%s` not exists\n", opm)
		} else if err.Error() == "Process exited with status 1" {
			fmt.Printf("[\033[0;32;31mERROR\033[m] command `%s` return 1\n", opm)
		} else {
			fmt.Printf("[\033[0;32;31mERROR\033[m] %s\n", err)
		}
	} else {
		fmt.Printf("%s\n", reply)
		fmt.Printf("[\033[1;33mOK\033[m][%s]\n\n", ipPort)
	}
}

func ExecutePush(ipPort, user, password, sources, destination string) {
	var authMethods []ssh.AuthMethod

	fmt.Printf("\033[1;33m")
	fmt.Printf("[%s]\n", ipPort)
	fmt.Printf("\033[m")

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
		Timeout: 1 * time.Second,
	}

	client, err := ssh.Dial("tcp", ipPort,  sshConfig)
	if err != nil {
		fmt.Printf("\033[0;32;31m");
		fmt.Printf("%s\n", err)
		fmt.Printf("\033[m")
		return
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		fmt.Printf("\033[0;32;31m");
		fmt.Printf("%s\n", err)
		fmt.Printf("\033[m")
		return
	}
	defer sftpClient.Close()

	//fmt.Println(source)
	filepathArray := strings.Split(sources, ",")
	for _, filepath:=range filepathArray {
		st, err := os.Stat(filepath)
		if err != nil {
			fmt.Printf("\033[0;32;31m");
			fmt.Printf("%s\n", err)
			fmt.Printf("\033[m")
			continue
		}
		if st.IsDir() {
			filename := path.Base(filepath)
			remotePath := path.Join(destination, filename)
			err = pushDirectory(sftpClient, filepath, remotePath)
		} else {
			err = pushFile(sftpClient, filepath, destination, st.Mode())
		}
		if err != nil {
			fmt.Printf("[\033[0;32;31mERROR\033[m] push %s error: %s\n", filepath, err.Error())
		} else {
			fmt.Printf("[\033[1;33mOK\033[m] push %s to %s#%s\n", filepath, ipPort, destination);
		}
	}
}

func ExecutePull(ipPort, user, password, localPath, remotePath string) {
	var authMethods []ssh.AuthMethod

	fmt.Printf("\033[1;33m")
	fmt.Printf("[%s]\n", ipPort)
	fmt.Printf("\033[m")

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

	client, err := ssh.Dial("tcp", ipPort,  sshConfig)
	if err != nil {
		fmt.Printf("\033[0;32;31m");
		fmt.Printf("%s\n", err)
		fmt.Printf("\033[m")
		return
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		fmt.Printf("\033[0;32;31m");
		fmt.Printf("%s\n", err)
		fmt.Printf("\033[m")
		return
	}
	defer sftpClient.Close()

	fInfo, err := sftpClient.Stat(remotePath)
	if err != nil {
		return
	}
	if fInfo.IsDir() {
		err = pullDirectory(sftpClient, localPath, remotePath)
		if err != nil {
			fmt.Printf("[\033[0;32;31mERROR\033[m] pull dir %s error: %s\n", remotePath, err.Error())
		} else {
			fmt.Printf("[\033[1;33mOK\033[m] pull dir %s#%s localhost#%s\n", ipPort, remotePath, localPath);
		}
	} else {
		filePath := path.Join(localPath, fInfo.Name())
		err = pullFile(sftpClient, filePath, remotePath, fInfo)
		if err != nil {
			fmt.Printf("[\033[0;32;31mERROR\033[m] pull %s error: %s\n", remotePath, err.Error())
		} else {
			fmt.Printf("[\033[1;33mOK\033[m] pull %s to %s#%s\n", filePath, ipPort, remotePath);
		}
	}
}

func pushFile(client *sftp.Client, localFilePath, remotePath string, mode os.FileMode) error {
	fmt.Println("push file")
	// 打开本地文件
	locF, err := os.Open(localFilePath)
	if err != nil {
		return fmt.Errorf("os.Open error, file:%s, %s\n", localFilePath, err)
	}
	defer locF.Close()

	// 上传到远端服务器的文件名,与本地路径末尾相同
	var remoteFileName = path.Base(localFilePath)
	//打开远程文件,如果不存在就创建一个
	remoteFP := path.Join(remotePath, remoteFileName)
	dstFile, err := client.Create(remoteFP)
	if err != nil {
		return fmt.Errorf("Create sftp session error: %s", path.Join(remotePath, remoteFileName))
	}
	//关闭远程文件
	defer dstFile.Close()
	//读取本地文件,写入到远程文件中(这里没有分快穿,自己写的话可以改一下,防止内存溢出)
	ff, err := ioutil.ReadAll(locF)
	if err != nil {
		return fmt.Errorf("ReadAll error: %s", localFilePath)
	}
	_, _ = dstFile.Write(ff)

	err = client.Chmod(remoteFP, mode)
	if err != nil {
		return fmt.Errorf("remote Chmod %s error:%s", remoteFP, err.Error())
	}
	return nil
}

func pushDirectory(client *sftp.Client, localPath, remotePath string) error {
	fmt.Println("push directory")
	//打开本地文件夹流
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		return fmt.Errorf("Path error:%s", err.Error())
	}
	//先创建最外层文件夹
	_ = client.Mkdir(remotePath)
	//遍历文件夹内容
	for _, backupDir := range localFiles {
		localFilePath := path.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		//判断是否是文件,是文件直接上传.是文件夹,先远程创建文件夹,再递归复制内部文件
		if backupDir.IsDir() {
			_ = client.Mkdir(remoteFilePath)

			err = pushDirectory(client, localFilePath, remoteFilePath)
			if err != nil {
				fmt.Printf("[\033[0;32;31mERROR\033[m] push %s error: %s\n", localFilePath, err.Error())
			}
		} else {
			err = pushFile(client, path.Join(localPath, backupDir.Name()), remotePath, backupDir.Mode())
			if err != nil {
				fmt.Printf("[\033[0;32;31mERROR\033[m] push %s error: %s\n", localFiles, err.Error())
			}
		}
	}

	fmt.Println(localPath + "  copy directory to remote server finished!")
	return nil
}

func pullFile(client *sftp.Client, localFile, remoteFile string, fInfo os.FileInfo) error {
	sf, err := client.Open(remoteFile)
	if err != nil {
		return fmt.Errorf("sftp.Client.Open error, file:%s, %s\n", remoteFile, err)
	}
	defer sf.Close()

	lf, err := os.Create(localFile)
	if err != nil {
		return fmt.Errorf("os.Create error: file=%s, errMsg=%s", localFile, err.Error())
	}
	defer lf.Close()

	_, err = sf.WriteTo(lf)
	if err != nil {
		return fmt.Errorf("sftp wrtieTo error: %s", err.Error())
	}
	err = lf.Chmod(fInfo.Mode())
	if err != nil {
		return fmt.Errorf("chmod %s error:%s", localFile, err.Error())
	}

	return nil
}

func pullDirectory(client *sftp.Client, localPath, remotePath string) error {
	remoteFiles, err := client.ReadDir(remotePath)
	if err != nil {
		return fmt.Errorf("Path error:%s", err.Error())
	}
	for _, subDir := range remoteFiles {
		subLocalPath := path.Join(localPath, subDir.Name())
		subRemotePath := path.Join(remotePath, subDir.Name())
		if subDir.IsDir() {
			// 递归执行
			// 创建本地目录
			_ = os.Mkdir(subLocalPath, subDir.Mode())
			err = pullDirectory(client, subLocalPath, subRemotePath)
		} else {
			err = pullFile(client, path.Join(localPath, subDir.Name()), remotePath, subDir)
		}
	}
	return nil
}