package main

import (
	"github.com/codeskyblue/go-sh"
	"github.com/russross/blackfriday"

	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	// "os/exec"
	"regexp"
)

func main() {

	var filename string // 参数文件名
	var Tmp string      // 临时html文件
	var Mobi string     // mobi名称
	var Title string
	var Author string
	var Comment string
	var Lang string

	flag.StringVar(&Title, "t", "", "标题")
	flag.StringVar(&Author, "a", "", "作者")
	flag.StringVar(&filename, "f", "", "文件")
	flag.StringVar(&Comment, "c", "", "简介")
	flag.StringVar(&Lang, "l", "zh-CN", "语言")
	flag.Parse()

	if filename == "" {
		if len(os.Args) > 1 {
			filename = os.Args[1]
			if len(os.Args) == 3 {
				Lang = os.Args[2]
			}
			ft := regexp.MustCompile(`^(.*).md|txt`).FindStringSubmatch(filename)
			if len(ft) > 1 {
				Tmp = ft[1] + ".html"
				Mobi = ft[1] + ".mobi"
				fts := strings.Split(ft[1], "-")
				if len(fts) == 2 {
					Title = strings.TrimSpace(fts[1])
					Author = strings.TrimSpace(fts[0])
				} else {
					Title = strings.TrimSpace(ft[1])
					Author = "UNKNOW"
				}
			} else {
				fmt.Println("err: use -f filename or input filename of type is md|text")
				os.Exit(-1)
			}
		} else {
			fmt.Println("err: use -f filename or input filename of type is md|text")
			os.Exit(-1)
		}
	}

	info, err := os.Stat(filename)
	if err != nil {
		fmt.Println("error: check your file is exist !")
		os.Exit(-1)
	}

	if Tmp == "" {
		Tmp = info.Name() + ".html"
	}
	if Mobi == "" {
		Mobi = info.Name() + ".mobi"
	}

	fmt.Println("read source file : ", filename)

	b, _ := ioutil.ReadFile(filename)
	re := blackfriday.HtmlRenderer(1, "title", "")
	Md := blackfriday.Markdown(b, re, blackfriday.EXTENSION_TABLES+blackfriday.EXTENSION_FENCED_CODE)

	if Title == "" {
		ft := regexp.MustCompile(`^(.*).md|txt`).FindStringSubmatch(filename)
		if len(ft) > 1 {
			fts := strings.Split(ft[1], "-")
			if len(fts) == 2 {
				Title = strings.TrimSpace(fts[1])
				Author = strings.TrimSpace(fts[0])
			} else {
				Title = strings.TrimSpace(ft[1])
				Author = "UNKNOW"
			}
		} else {
			// 读取文件获取标题
		}
	}

	fmt.Printf("Title: %s\nAuthor: %s \n", Title, Author)

	fmt.Println("create html file ...")
	// 创建历史文件
	tmpfile, _ := os.Create(Tmp)
	defer os.Remove(Tmp)
	defer fmt.Println("remove html file...")
	defer tmpfile.Close()

	tmpfile.WriteString(fmt.Sprintf("<html><head><meta http-equiv='content-language' content='zh-CN' /><meta http-equiv='Content-type' content='text/html; charset=utf-8'><meta name='Author' content='%s'><title>%s</title></head><body>%s</body></html>", Author, Title, regexp.MustCompile(`\n`).ReplaceAllString(string(Md), "")))

	fmt.Println("use ebook-convert create mobi...")

	fmt.Print(fmt.Sprintf("ebook-convert %s %s --authors %s --comments '%s' --level1-toc '//h:h1' --level2-toc '//h:h2' --language '%s'\n", Tmp, Mobi, Author, Comment, Lang))

	sh.Command("ebook-convert", Tmp, Mobi, "--authors", Author, "--comments", Comment, "--level1-toc", "//h:h1", "--level2-toc", "//h:h2", "--language", Lang).Run()

	fmt.Println("complete!")
}
