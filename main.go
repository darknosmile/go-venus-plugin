package main

import (
	"fmt"
	"os"
	"io"
	"bufio"
	"strings"
	"flag"
)

func main(){
	fileName := flag.String("file", "init.ini", "file path! default is init.ini")
	flag.Parse()
	file, err := os.OpenFile(*fileName, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	var size = stat.Size()
	fmt.Println("file size=", size)

	buf := bufio.NewReader(file)
	param := make(map[string]string)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		fmt.Println(line)
		tmp :=strings.Split(line,"=")
		if(len(tmp)==2){
			param[tmp[0]]=tmp[1]
		}
		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return
			}
		}
	}

	venus := Conncect(param["venus_address"])
	venus.AuthByDummy("venus")
	var json = param["json"]
	venus.Request(param["interface"], param["version"], json)
	fmt.Println("")
}
