package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Memo struct {
	TimeStr string
	Content string
	Time    time.Time
}

func main() {
	inputFile := "江楠大盗的笔记.html"
	outputFile := "output.txt"

	// 读取整个文件内容
	content, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	htmlStr := string(content)

	// 正则匹配 memo 块
	memoRegex := regexp.MustCompile(`(?s)<div class="memo">.*?<div class="time">(.*?)</div>.*?<div class="content">(.*?)</div>.*?</div>`)
	matches := memoRegex.FindAllStringSubmatch(htmlStr, -1)

	if matches == nil {
		fmt.Println("No memos found.")
		return
	}

	var memos []Memo

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}
		timeStr := strings.TrimSpace(match[1])
		contentHtml := match[2]

		// 处理 content
		// 1. 将 </p> 替换为换行符，以便分段
		contentHtml = strings.ReplaceAll(contentHtml, "</p>", "\n")
		contentHtml = strings.ReplaceAll(contentHtml, "<br>", "\n")
		contentHtml = strings.ReplaceAll(contentHtml, "<br/>", "\n")

		// 2. 去除所有 HTML 标签
		reTag := regexp.MustCompile(`<[^>]+>`)
		contentStr := reTag.ReplaceAllString(contentHtml, "")

		// 3. HTML 解码
		contentStr = html.UnescapeString(contentStr)

		// 4. 按行处理：去除首尾空格，过滤空行
		lines := strings.Split(contentStr, "\n")
		var cleanLines []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			// 过滤以 # 开头的行
			if strings.HasPrefix(line, "#") {
				continue
			}
			cleanLines = append(cleanLines, line)
		}

		finalContent := strings.Join(cleanLines, "\n")

		// 如果内容为空，则跳过该条笔记
		if finalContent == "" {
			continue
		}

		// 解析时间以便排序
		// 假设时间格式为 "2006-01-02 15:04:05"
		parsedTime, err := time.Parse("2006-01-02 15:04:05", timeStr)
		if err != nil {
			// 如果解析失败，尝试只保留数字和分隔符后再解析，或者直接使用当前时间作为默认值
			// 这里简单处理：如果解析失败，打印警告并按字符串排序
			fmt.Printf("Warning: could not parse time '%s': %v\n", timeStr, err)
		}

		memos = append(memos, Memo{
			TimeStr: timeStr,
			Content: finalContent,
			Time:    parsedTime,
		})
	}

	// 按时间排序（正序：旧的在前）
	sort.Slice(memos, func(i, j int) bool {
		return memos[i].Time.Before(memos[j].Time)
	})

	var outputBuilder strings.Builder
	for _, memo := range memos {
		outputBuilder.WriteString(memo.TimeStr + "\n")
		outputBuilder.WriteString(memo.Content + "\n\n")
	}

	// 写入输出文件
	err = ioutil.WriteFile(outputFile, []byte(outputBuilder.String()), 0644)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		return
	}

	fmt.Printf("Successfully processed %d memos. Output written to %s\n", len(memos), outputFile)
}
