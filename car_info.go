package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

type tCarInfo struct {
	Car    []string
	Mobile []string
}

var (
	reCar,
	reMobile *regexp.Regexp
	msgFile     []string
	carInfo     [][]string
	msgFileName = "0102.txt"
)

func init() {
	// 0. 基本规则
	// 第 1 位为省份简称 (汉字)
	// 第 2 位为发牌机关代号 (A-Z)
	re0 := `[京津沪渝冀豫云辽黑湘皖鲁新苏浙赣鄂桂甘晋蒙陕吉闽贵粤青藏川宁琼使领][A-Z]`

	// 1. 新能源车牌
	// 第 3-8 位为序号 (比传统车牌多一位) 规则如下:
	// 大型车: 第 1-5 位必须是数字, 第 6 位只能是字母 D 或 F
	// 小型车: 第 1 位只能是字母 D 或 F, 第 2 位可以是数字或字母, 第 3-6 位必须是数字
	re1 := `(([0-9]{5}[DF])|([DF]([A-HJ-NP-Z0-9])[0-9]{4}))`

	// 2. 传统车牌
	// 第 3-7 位为序号 (字母或数字组成, 但无字母 I 和 O
	// 另外最后一位可能是 (挂学警港澳使领) 中的一个汉字
	re2 := `([A-HJ-NP-Z0-9]{4}[A-HJ-NP-Z0-9挂学警港澳使领])`

	reCar = regexp.MustCompile(fmt.Sprintf("(%s%s)|(%s%s)", re0, re1, re0, re2))

	// 手机号
	reMobile = regexp.MustCompile(`1\d{10}`)
}

// 执行文件提取
func run() {
	if scanMsgFile() {
		// 写入结果文件
		// fmt.Println(carInfo)
		if err := saveCarInfo(); err != nil {
			fmt.Printf("\n提取结果:\n------\n%v\n------\n", carInfo)
			log.Println("xxx.文件保存失败:", err)
		} else {
			log.Println("OK")
		}
	} else {
		log.Println("xxx.提取消息结果为空")
	}

	fmt.Println("\n请关闭窗口或 Ctrl+C 结束程序(10 秒后自动退出)")
	time.Sleep(10 * time.Second)
}

// 按文件提取信息
func scanMsgFile() bool {
	for _, f := range msgFile {
		log.Println("....正在读取文件:", f)

		file, err := os.Open(f)
		if err != nil {
			log.Println("xxx.读取文件失败:", f)
			continue
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if info := findCarInfo(scanner.Text()); len(info.Car) > 0 {
				carInfo = append(carInfo, []string{
					strings.Join(info.Car, "|"),
					strings.Join(info.Mobile, "|"),
				})
			}
		}
	}

	return len(carInfo) > 0
}

// 按行提取内容信息
func scanMsg(s string) string {
	var res strings.Builder

	for _, line := range strings.Split(s, "\n") {
		if info := findCarInfo(line); len(info.Car) > 0 {
			res.WriteString(strings.Join(info.Car, "|"))
			res.WriteString("\t")
			res.WriteString(strings.Join(info.Mobile, "|"))
			res.WriteString("\n")
		}
	}

	return res.String()
}

// 查找车牌和手机信息
func findCarInfo(s string) *tCarInfo {
	return &tCarInfo{
		Car:    reCar.FindAllString(s, -1),
		Mobile: reMobile.FindAllString(reCar.ReplaceAllString(s, " "), -1),
	}
}

// 保存结果
func saveCarInfo() error {
	f, err := os.Create(path.Base(msgFile[0]) + ".csv")
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	// 写入UTF-8 BOM
	_, _ = f.WriteString("\xEF\xBB\xBF")

	w := csv.NewWriter(f)
	if err := w.WriteAll(carInfo); err != nil {
		return err
	}

	w.Flush()

	return nil
}
