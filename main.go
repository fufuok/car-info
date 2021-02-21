package main

import (
	"log"
	"os"
	"strconv"
	"time"
)

func init() {
	// 消息文件名, 默认为日期.txt (0220.txt)
	msgFile = os.Args[1:]
	if len(msgFile) < 1 {
		msgFile = []string{time.Now().Format(msgFileName)}
	} else {
		for k, v := range msgFile {
			// 处理 -1 表示昨天日期
			if i, err := strconv.Atoi(v); err == nil {
				msgFile[k] = time.Now().AddDate(0, 0, i).Format(msgFileName)
			}
		}
	}
}

func main() {
	if msgFile[0] == "web" {
		if err := initWebServer(); err != nil {
			log.Println("Web 服务启动失败", err)
			time.Sleep(10 * time.Second)
			os.Exit(1)
		}
	}

	// 命令模式
	run()
}
