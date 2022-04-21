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

func FIFO(blocksNum int, pagesQueue PagesQueue) float64 {
	var blocks Blocks
	blocks.Init(blocksNum, pagesQueue)
	// 当未命中时，进入loop2，否则持续执行loop1
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
	fmt.Printf("缺页率：%.1f%%\n\n", 100.0*(float64(blocks.MissingPages)/float64(blocks.Pages)))
	return 100.0 * (float64(blocks.MissingPages) / float64(blocks.Pages))
}

// 向后查找离得最远的页面序号，如不存在，则选择最后一个不存在的页面序号
func GetOPIIndex(data []Block, pagesQueue PagesQueue, index_i int) int {
	flag := false
	result := 0
	for index_j, j := range data {
	loop:
		for index_k, k := range pagesQueue.Data[index_i:] {
			if k == j.PageNum {
				data[index_j].Count = int64(index_k)
				break loop
			}
		}
		flag = true
		result = index_j
	}
	if flag {
		return result
	}

	max := int64(0)
	for index, i := range data {
		if i.Count > max {
			result = index
			max = i.Count
		}
	}
	for item := range data {
		data[item].Count = 0
	}
	return result
}

func OPI(blocksNum int, pagesQueue PagesQueue) float64 {
	var blocks Blocks
	blocks.Init(blocksNum, pagesQueue)
	// 当未命中时，进入loop2，否则持续执行loop1
loop1:
	for index_i, i := range pagesQueue.Data {
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
		for index_j, j := range blocks.Data {
			if j.PageNum == 0 && !IsFilled(blocks.Data) {
				blocks.Data[index_j] = Block{
					Count:   0,
					PageNum: int64(i),
				}
				break loop2
			} else if IsFilled(blocks.Data) {
				blocks.Data[GetOPIIndex(blocks.Data, pagesQueue, index_i)] = Block{
					Count:   0,
					PageNum: int64(i),
				}
				break loop2
			}
		}
	}
	fmt.Printf("缺页率：%.1f%%\n\n", 100.0*(float64(blocks.MissingPages)/float64(blocks.Pages)))
	return 100.0 * (float64(blocks.MissingPages) / float64(blocks.Pages))
}

// 获取count最大的块内索引
func GetLRUIndex(data []Block) int {
	result := 0
	max := int64(0)
	for index, i := range data {
		if i.Count > max {
			result = index
			max = i.Count
		}
	}
	return result
}

func LRU(blocksNum int, pagesQueue PagesQueue) float64 {
	var blocks Blocks
	blocks.Init(blocksNum, pagesQueue)
	// 当未命中时，进入loop2，否则持续执行loop1，并对块内页面作计数
loop1:
	for _, i := range pagesQueue.Data {
		fmt.Printf("当前访问页面序列：%d", i)
		fmt.Printf(",物理块内情况：")
		for _, tt := range blocks.Data {
			fmt.Printf("%d ", tt.PageNum)
		}
		fmt.Println()
		flag := false
		for index_k, k := range blocks.Data {
			if i == k.PageNum {
				blocks.Data[index_k].Count = 0
				flag = true
			} else {
				blocks.Data[index_k].Count += 1
			}
		}
		if flag {
			continue loop1
		}
		blocks.MissingPages += 1
	loop2:
		for index_j, j := range blocks.Data {
			if j.PageNum == 0 && !IsFilled(blocks.Data) {
				blocks.Data[index_j] = Block{
					Count:   0,
					PageNum: int64(i),
				}
				break loop2
			} else if IsFilled(blocks.Data) {
				blocks.Data[GetLRUIndex(blocks.Data)] = Block{
					Count:   0,
					PageNum: int64(i),
				}
				break loop2
			}
		}
	}
	fmt.Printf("缺页率：%.1f%%\n\n", 100.0*(float64(blocks.MissingPages)/float64(blocks.Pages)))
	return 100.0 * (float64(blocks.MissingPages) / float64(blocks.Pages))
}

// 获取输入
func GetInput() string {
	data := GetData().Data
	fmt.Printf("访问序列：")
	fmt.Println(data)
	fmt.Printf("请输入置换算法（1-FIFO，2-OPI，3-LRU，q-quit)：")
	var algorthmType string
	fmt.Scanf("%s\n", &algorthmType)
	return algorthmType
}

func main() {

loop:
	for {
		var blocksNum int
		switch algorthmType := GetInput(); algorthmType {
		case "1":
			fmt.Printf("请输入最小物理块数：")
			fmt.Scanf("%d\n", &blocksNum)
			FIFO(blocksNum, GetData())
		case "2":
			fmt.Printf("请输入最小物理块数：")
			fmt.Scanf("%d\n", &blocksNum)
			OPI(blocksNum, GetData())
		case "3":
			fmt.Printf("请输入最小物理块数：")
			fmt.Scanf("%d\n", &blocksNum)
			LRU(blocksNum, GetData())
		case "q":
			break loop
		}
		PrintLine()
	}
}
