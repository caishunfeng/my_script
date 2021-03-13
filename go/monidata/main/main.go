package main

/**
1、定时检测输入的服务器列表的/home/dsc/log/*monitor.log.*
2、若存在该文件,则将文件改名为*monidata.log.*,则读取该文件的数据,存入influxdb
3、当文件读取完则删掉,等待下次任务执行
4、监控日志大小切割为5m,本机需初始化/home/dsc/log/路径
**/

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/robfig/cron"

	log "code.google.com/p/log4go"
)

var servers string
var userName string
var password string
var influxAdds string
var influxUserName string
var influxPwd string

/**
nohup ./monidata  172.16.3.53,172.16.3.52,172.16.3.42,172.16.3.48,172.16.3.21,172.16.3.22,172.16.3.18,172.16.3.19,172.16.3.29,172.16.3.30,172.16.3.27 name 'password' http://127.0.0.1:8086 user pass &
nohup ./monidata  114.55.146.52,114.55.253.132,120.26.75.230,120.27.250.224,121.41.40.91,121.41.41.168,121.43.60.224 name 'password' http://127.0.0.1:8086 user pass &
*/
func main() {

	openLog()

	args := os.Args
	if len(args) <= 1 {
		log.Info("请在启动时输入命令参数,参数一为监控服务器列表,参数二为登录名,参数三为登录密码")
		return
	}

	servers = args[1]
	userName = args[2]
	password = args[3]
	influxAdds = args[4]
	influxUserName = args[5]
	influxPwd = args[6]

	log.Info("servers:%s", servers)
	log.Info("userName:%s,pwd:%s", userName, password)
	log.Info("influxdb,adds:%s,username:%s,pwd:%s", influxAdds, influxUserName, influxPwd)

	InitInfluxdb(influxAdds, influxUserName, influxPwd)

	log.Info("influxdb 初始化结束")

	getRemoteFileCron := cron.New()
	getRemoteFileCron.AddFunc("0 0/10 * * * ?", startGetRemoteFile)
	getRemoteFileCron.Start()
	defer getRemoteFileCron.Stop()

	// startGetRemoteFile()

	for {
		time.Sleep(60 * time.Second)
		log.Info("主线程保持存活")
	}
}

//开启日志记录
func openLog() {
	flw := log.NewFileLogWriter("monidata.log", false)
	flw.SetFormat("[%D %T] [%L] (%S) %M")
	flw.SetRotate(true)
	flw.SetRotateSize(0)
	flw.SetRotateLines(100000)
	flw.SetRotateDaily(false)
	log.AddFilter("file", log.FINE, flw)
}

func startGetRemoteFile() {
	log.Info("获取远程监控日志文件任务开始")
	serverList := strings.Split(servers, ",")
	for i := 0; i < len(serverList); i++ {
		log.Info("开始获取,服务器地址:%s", serverList[i])
		go getRemoteFile(serverList[i], userName, password, 22)
	}
}

//获取远程服务器的日志文件，将其抓取回本地
//删除远程服务器日志文件
//处理抓取回来的本地日志文件
func getRemoteFile(host, user, password string, port int) {
	client, err := Connect(host, user, password, port)
	if err != nil {
		log.Error("远程连接SSH异常,原因:%v", err)
		return
	}
	defer client.Close()

	lsCmd := "ls /home/dsc/log/"

	lsRes, err := RunCmd(client, lsCmd)
	if err != nil {
		log.Error("执行远程命令异常,命令:%s,原因:%v", lsCmd, err)
		return
	}

	files := strings.Split(lsRes, "\n")
	newFiles := make(map[string]string)
	newCount := 0
	for i := 0; i < len(files); i++ {
		oldfileName := files[i]
		if strings.Contains(oldfileName, "monitor.log.") {
			newFileName := host + "-" + strings.Replace(oldfileName, "monitor.log.", "monidata.log.", -1)
			mvCmd := fmt.Sprintf("mv /home/dsc/log/%s /home/dsc/log/%s", oldfileName, newFileName)
			_, err := RunCmd(client, mvCmd)
			if err != nil {
				log.Error("执行远程命令异常,命令:%s,原因:%v", lsCmd, err)
				return
			}
			newFiles[newFileName] = newFileName
			newCount++
		}
	}

	if newCount > 0 {
		sftpClient, err := ConnectSFTP(host, user, password, port)
		if err != nil {
			log.Error("远程连接SFTP异常,原因:%v", err)
		}
		defer sftpClient.Close()

		//创建本地文件
		for fileName, _ := range newFiles {
			localFileDir := "/home/dsc/log"
			localFileName := fileName + "_" + GetRandomString(3)
			localFilePath := localFileDir + "/" + localFileName
			remoteFileDir := "/home/dsc/log"
			remoteFilePath := remoteFileDir + "/" + fileName
			isExist, _ := PathExists(remoteFileDir)
			if !isExist {
				//创建目录
				err := os.Mkdir(localFileDir, os.ModePerm)
				if err != nil {
					fmt.Printf("文件夹创建失败,错误原因:%s", err)
					return
				}
			}
			isFileExist, _ := PathExists(localFilePath)
			var targetFile *os.File
			if isFileExist {
				targetFile, err = os.Open(localFilePath)
				if err != nil {
					log.Error("打开本地文件失败,path:%s,原因:%v", localFilePath, err)
					return
				}
			} else {
				targetFile, err = os.Create(localFilePath)
				if err != nil {
					log.Error("打开本地文件失败,path:%s,原因:%v", localFilePath, err)
					return
				}
			}
			defer targetFile.Close()

			srcFile, err := sftpClient.Open(remoteFilePath)
			if err != nil {
				log.Error("打开远程文件失败,path:%s,原因:%v", remoteFilePath, err)
				return
			}
			defer srcFile.Close()

			buf := make([]byte, 1024)
			for {
				n, _ := srcFile.Read(buf)
				if n == 0 {
					break
				}
				log.Info("读取了%d个字节", n)

				targetFile.Write(buf)
			}

			log.Info("远程传输文件结束,文件名:%s", fileName)

			go handleFile(localFilePath)

			rmCmd := "rm -rf /home/dsc/log/" + fileName
			rmRes, err := RunCmd(client, rmCmd)
			if err != nil {
				log.Error("删除文件异常,文件名:%s,原因:%v", fileName, err)
				return
			}
			log.Info("删除文件%s成功,返回结果:%s", fileName, rmRes)
		}
	}
}

//处理日志文件
func handleFile(filePath string) {
	log.Info("开始处理监控日志文件,path:%s", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		log.Error("打开本地文件失败,path:%s,原因:%v", filePath, err)
		return
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	reader := bufio.NewReader(file)
	count := 0

	bp, err := CreateBatchPoints()
	if err != nil {
		log.Error("生成influxdb批量节点失败,原因:%v", err)
		return
	}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Error("读取文件失败,path:%s,原因:%v", filePath, err)
			return
		}
		if line == "" {
			break
		}

		index := strings.Index(line, "customerId")
		if index >= 0 {
			newLine := line[index:]

			log.Info("读取到内容:%s", newLine)

			kvs := strings.Split(newLine, ",")

			maps := make(map[string]string)
			for j := 0; j < len(kvs); j++ {
				kv := strings.Split(kvs[j], ":")
				key := kv[0]
				value := kv[1]
				maps[key] = value
			}

			measurement := ""
			if strings.Contains(filePath, "qm") {
				measurement = "qm"
			} else if strings.Contains(filePath, "top") {
				measurement = "top"
			} else if strings.Contains(filePath, "lppz") {
				measurement = "lppz"
			}

			if count == 1000 {
				//满一千条就写入
				log.Info("满一千条,开始批量写入influxDB")
				oldbp := bp
				SendInfluxDB(oldbp)
				bp, err = CreateBatchPoints()
				if err != nil {
					log.Error("生成influxdb批量节点失败,原因:%v", err)
					return
				}
				//清零
				count = 0
			}
			err := AddInfluxQueue(bp, measurement, maps["customerId"], maps["busiType"], maps["success"], maps["fail"], maps["timeStamp"])
			if err != nil {
				log.Error("添加influxdb队列失败,原因:%v", err)
				return
			}
			count++
		}
	}

	//将最后一次的数据写入
	if count > 0 {
		log.Info("最后一次count>0,开始批量写入influxDB")
		SendInfluxDB(bp)
	}

	err = os.Remove(filePath)
	if err != nil {
		log.Error("删除本地文件失败,path:%s,原因:%v", filePath, err)
		return
	}
}
