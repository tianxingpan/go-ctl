// Package controller provides SSH control method
package controller

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type SFTP struct {
	client *ssh.Client
	sftpClient *sftp.Client
	addr string
	isContinue bool		// 在下载目录或者批量下载文件时，中间遇到错误是否继续下载
}

// 初始化SFTP
func (s *SFTP) Init(addr, user, password string) error {
	s.addr = addr
	var err error
	s.client, err = NewSSHClient(addr, user, password)
	if err != nil {
		fmt.Printf("\033[0;32;31m");
		fmt.Printf("%s\n", err)
		fmt.Printf("\033[m")
		return err
	}
	s.sftpClient, err = sftp.NewClient(s.client)
	return err
}

// 设置是否继续执行标志
func (s *SFTP) SetContinue(c bool) {
	s.isContinue = c
}

func (s *SFTP) GetSFTP() *sftp.Client {
	return s.sftpClient
}

func (s *SFTP) Close() {
	if s.sftpClient != nil {
		s.sftpClient.Close()
	}
	if s.client != nil {
		s.client.Close()
	}
}

// 从远程服务器下载
func (s *SFTP) Pull(localPath, remotePath string) error {
	if s.sftpClient == nil {
		err := fmt.Errorf("")
		fmt.Printf("\033[0;32;31m")
		fmt.Printf("%s\n", err)
		fmt.Printf("\033[m")
		return err
	}
	fInfo, err := s.sftpClient.Stat(remotePath)
	if err != nil {
		fmt.Printf("\033[0;32;31m")
		fmt.Printf("%s\n", err)
		fmt.Printf("\033[m")
		return err
	}
	if fInfo.IsDir() {
		err = s.pullDirectory(localPath, remotePath)
		if err != nil {
			fmt.Printf("[\033[0;32;31mERROR\033[m] pull dir %s error: %s\n", remotePath, err.Error())
		} else {
			fmt.Printf("[\033[1;33mOK\033[m] pull dir %s#%s localhost#%s\n", s.addr, localPath, remotePath);
		}
	} else {
		filePath := path.Join(localPath, fInfo.Name())
		err = s.pullFile(filePath, remotePath, fInfo)
		if err != nil {
			fmt.Printf("[\033[0;32;31mERROR\033[m] pull %s error: %s\n", remotePath, err.Error())
		} else {
			fmt.Printf("[\033[1;33mOK\033[m] pull %s to %s#%s\n", filePath, s.addr, remotePath);
		}
	}
	return nil
}

// 上传文件或目录到远程服务器
func (s *SFTP) Push(localPath, remotePath string) error {
	if s.sftpClient == nil {
		return fmt.Errorf("SFTP Client not created")
	}

	filepathArray := strings.Split(localPath, ",")
	for _, filepath:=range filepathArray {
		st, err := os.Stat(filepath)
		if err != nil {
			continue
		}
		if st.IsDir() {
			filename := path.Base(filepath)
			destination := path.Join(remotePath, filename)
			err = s.pushDirectory(filepath, destination)
		} else {
			err = s.pushFile(filepath, remotePath, st.Mode())
		}
		if err != nil {
			fmt.Printf("[\033[0;32;31mERROR\033[m] push %s error: %s\n", filepath, err.Error())
		} else {
			fmt.Printf("[\033[1;33mOK\033[m] push %s to %s#%s\n", filepath, s.addr, remotePath);
		}
	}

	return nil
}

// 从远程服务器下拉文件
func (s *SFTP) pullFile(localFile, remoteFile string, fInfo os.FileInfo) error {
	sf, err := s.sftpClient.Open(remoteFile)
	if err != nil {
		return fmt.Errorf("Open sftp failed, errMsg:%s\n", err.Error())
	}
	defer sf.Close()

	lf, err := os.Create(localFile)
	if err != nil {
		return fmt.Errorf("os.Create error, file:%s, %s\n", localFile, err.Error())
	}
	defer lf.Close()

	_, err = sf.WriteTo(lf)
	if err != nil {
		return fmt.Errorf("remote to local failed, local:%s, remoteIP:%s/%s, errMsg:%s\n", localFile, s.addr, remoteFile, err.Error())
	}
	err = lf.Chmod(fInfo.Mode())
	if err != nil {
		return fmt.Errorf("Chmod %s failed, errMsg:%s\n", localFile, err.Error())
	}

	return nil
}

// 从远程服务器下拉目录
func (s *SFTP) pullDirectory(localPath, remotePath string) error {
	remoteFiles, err := s.sftpClient.ReadDir(remotePath)
	if err != nil {
		return fmt.Errorf("Path error:%s\n", err.Error())
	}
	for _, subDir := range remoteFiles {
		subLocalPath := path.Join(localPath, subDir.Name())
		subRemotePath := path.Join(remotePath, subDir.Name())
		if subDir.IsDir() {
			// 创建本地目录
			_ = os.Mkdir(subLocalPath, subDir.Mode())
			err = s.pullDirectory(subLocalPath, subRemotePath)
		} else {
			err = s.pullFile(path.Join(localPath, subDir.Name()), remotePath, subDir)
		}
	}
	return nil
}

// 将本地文件上传到远程服务器
func (s *SFTP) pushFile(localFilePath, remotePath string, mode os.FileMode) error {
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
	dstFile, err := s.sftpClient.Create(remoteFP)
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

	err = s.sftpClient.Chmod(remoteFP, mode)
	if err != nil {
		return fmt.Errorf("remote Chmod %s error:%s", remoteFP, err.Error())
	}
	return nil
}

// 将本地目录上传到远程服务器
func (s *SFTP) pushDirectory(localPath, remotePath string) error {
	//打开本地文件夹流
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		return fmt.Errorf("Path error:%s\n", err.Error())
	}
	//先创建最外层文件夹
	_ = s.sftpClient.Mkdir(remotePath)
	//遍历文件夹内容
	for _, backupDir := range localFiles {
		localFilePath := path.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		//判断是否是文件,是文件直接上传.是文件夹,先远程创建文件夹,再递归复制内部文件
		if backupDir.IsDir() {
			_ = s.sftpClient.Mkdir(remoteFilePath)

			err = s.pushDirectory(localFilePath, remoteFilePath)
			if err != nil {
				fmt.Printf("[\033[0;32;31mERROR\033[m] push %s error: %s\n", localFilePath, err.Error())
			}
		} else {
			err = s.pushFile(localFilePath, remotePath, backupDir.Mode())
			if err != nil {
				fmt.Printf("[\033[0;32;31mERROR\033[m] push %s error: %s\n", localFiles, err.Error())
			}
		}
	}

	//fmt.Println(localPath + "  copy directory to remote server finished!")
	return nil
}
