package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	args := os.Args
	var url = ""
	var ts = []string{}

	if len(args) == 1 {
		usage()
	} else if len(args) == 2 {
		if args[1] == "-h" || args[1] == "--help" {
			usage()
		} else {
			usage()
		}
	} else if len(args) == 3 {
		if args[1] == "-url" {
			url = args[2]

			if !strings.HasSuffix(url, ".m3u8") {
				fmt.Println("please input correct m3u8 file url!!")
				os.Exit(1)
			}

			filename := downloadFile(url)

			// 解析m3u8文件
			ts = parseM3u8(filename)

			outfile := string(strings.Split(filename, ".")[0]) + ".mp4"
			fd, _ := os.OpenFile(outfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
			defer fd.Close()
			for _, v := range ts {
				fmt.Println("download", v)
				rs, err := http.Get(v)
				if err != nil {
					panic(err)
				}
				io.Copy(fd, rs.Body)
			}
		} else {
			usage()
		}
	} else {
		usage()
	}
}

// 解析m3u8文件
func parseM3u8(path string) []string {
	var ret = []string{}

	handler, err := os.Open(path)
	if err != nil {
		fmt.Println("Parse m3u8 file error ", err)
	}
	defer handler.Close()

	buffer := bufio.NewReader(handler)

	for {
		line, _, err := buffer.ReadLine()
		if err == io.EOF {
			break
		}

		if strings.HasPrefix(string(line), "http") {
			ret = append(ret, string(line))
		}
	}

	return ret
}

// show usage
func usage() {
	var rt = "Usage:m3u8.exe options [option_value]\n\nOptions:\n-h\t\tget help\n-url\t\tset the m3u8 file url\n"

	fmt.Print(rt)
	os.Exit(1)
}

// 下载ts文件
func downloadFile(path string) string {
	var pathinfo = strings.Split(path, "/")
	var filename = pathinfo[len(pathinfo)-1]
	res, err := http.Get(path)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	io.Copy(f, res.Body)
	return filename
}
