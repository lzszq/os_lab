package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Blocks struct {
	MissingPages int64
	Pages        int64
	Data         []Block
}

type Block struct {
	Count   int64
	PageNum int64
}

type PagesQueue struct {
	Data []int64
}

// 读取测试数据
func GetData() PagesQueue {
	var (
		file    *os.File
		err     error
		content []byte
	)

	pagesQueue := PagesQueue{
		Data: nil,
	}

	file, err = os.Open("./pages_queue.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, _ = ioutil.ReadAll(file)
	for _, i := range strings.Split(string(content), ",") {
		f, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			panic(err)
		}
		pagesQueue.Data = append(pagesQueue.Data, f)
	}
	return pagesQueue
}

// 初始化Blocks
func (blocks *Blocks) Init(blocksNum int, pagesQueue PagesQueue) {
	blocks.MissingPages = 0
	blocks.Pages = int64(len(pagesQueue.Data))
	blocks.Data = make([]Block, blocksNum)
}

// Print break-line
func PrintLine() {
	fmt.Println("------------------------------------------------------")
}

func IsFilled(data []Block) bool {
	flag := true
	for _, i := range data {
		if i.PageNum == 0 {
			flag = false
		}
	}
	return flag
}

func FIFO(blocksNum int, pagesQueue PagesQueue) {
	var blocks Blocks
	blocks.Init(blocksNum, pagesQueue)
loop1:
	for _, i := range pagesQueue.Data {
		fmt.Printf("当前访问页面序列：%d", i)
		fmt.Printf(",物理块内情况：")
		for _, tt := range blocks.Data {
			fmt.Printf("%d ", tt.PageNum)
		}
		fmt.Println()
		for _, k := range blocks.Data {
			if i == k.PageNum {
				continue loop1
			}
		}
		blocks.MissingPages += 1
	loop2:
		for index, j := range blocks.Data {
			if j.PageNum == 0 && !IsFilled(blocks.Data) {
				blocks.Data[index] = Block{
					Count:   0,
					PageNum: int64(i),
				}
				break loop2
			} else if IsFilled(blocks.Data) {
				t := Block{
					Count:   0,
					PageNum: int64(i),
				}
				blocks.Data = append(blocks.Data[1:], t)
				break loop2
			}
		}
	}
	fmt.Printf("缺页率：%.1f%%", 100.0*(float64(blocks.MissingPages)/float64(blocks.Pages)))
}

func OPI(blocksNum int, pagesQueue PagesQueue) {

}

func LRU(blocksNum int, pagesQueue PagesQueue) {

}

func main() {
	FIFO(20, GetData())
}
