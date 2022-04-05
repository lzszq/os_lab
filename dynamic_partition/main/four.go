package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

// 使用结构体
type System struct {
	P Processs
	S Partitions
}

// 进程数组
type Processs struct {
	Data []Process
}

// 进程结构体
type Process struct {
	Name             string
	Need             int
	AllocatePosition int
}

// 空间
type Partitions struct {
	Data []Partition
}

// 分区
type Partition struct {
	Sequence  int
	Partition int
}

// 实现结构体排序接口
type ByPartition []Partition

func (a ByPartition) Len() int { return len(a) }
func (a ByPartition) Less(i, j int) bool {
	if a[i].Partition == a[j].Partition {
		return a[i].Sequence < a[j].Sequence
	}
	return a[i].Partition < a[j].Partition
}
func (a ByPartition) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// 实现结构体排序接口
type BySequence []Partition

func (a BySequence) Len() int { return len(a) }
func (a BySequence) Less(i, j int) bool {
	return a[i].Sequence < a[j].Sequence
}
func (a BySequence) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// 首次适应算法
func (s *System) FirstFit() (System, string) {
	system := s.deepCopy()
	for i := range system.P.Data {
		for j := range system.S.Data {
			if system.P.Data[i].Need < system.S.Data[j].Partition {
				system.P.Data[i].AllocatePosition = system.S.Data[j].Sequence
				system.S.Data[j].Partition -= system.P.Data[i].Need
				break
			}
		}
		if system.P.Data[i].AllocatePosition == -1 {
			return system, "fail"
		}
	}
	return system, "success"
}

// 循环首次适应算法
func (s *System) NextFit() (System, string) {
	system := s.deepCopy()
	j := 0
	cnt := 0
	for i := range system.P.Data {
		for {
			if system.P.Data[i].Need < system.S.Data[j].Partition {
				system.P.Data[i].AllocatePosition = system.S.Data[j].Sequence
				system.S.Data[j].Partition -= system.P.Data[i].Need
				break
			}
			if cnt == len(system.S.Data) {
				break
			}
			j += 1
			cnt += 1
		}
		if system.P.Data[i].AllocatePosition == -1 {
			return system, "fail"
		}
	}
	return system, "success"
}

// 最佳适应算法
func (s *System) BestFit() (System, string) {
	system := s.deepCopy()
	sort.Sort(ByPartition(system.S.Data))
	system, flag := system.FirstFit()
	return system, flag
}

// 最坏适应算法
func (s *System) WorstFit() (System, string) {
	system := s.deepCopy()
	sort.Sort(sort.Reverse(ByPartition(system.S.Data)))
	system, flag := system.FirstFit()
	return system, flag
}

// 读取测试数据
func (s *System) GetData() {
	var (
		file    *os.File
		err     error
		content []byte
	)

	file, err = os.Open("./free.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, _ = ioutil.ReadAll(file)
	for index, i := range strings.Split(string(content), ",") {
		f, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			panic(err)
		}
		partition := Partition{
			Sequence:  index,
			Partition: int(f),
		}
		s.S.Data = append(s.S.Data, partition)
	}

	file, err = os.Open("./process.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, _ = ioutil.ReadAll(file)
	for index, i := range strings.Split(string(content), ",") {
		f, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			panic(err)
		}
		process := Process{
			Name:             string(rune(int(index) + 65)),
			Need:             int(f),
			AllocatePosition: -1,
		}
		s.P.Data = append(s.P.Data, process)
	}
}

// 对数据进行深拷贝
func (src *System) deepCopy() System {
	var dst = new(System)
	b, _ := json.Marshal(src)
	json.Unmarshal(b, &dst)
	return System{
		P: dst.P,
		S: dst.S,
	}
}

// 输出数据
func (system *System) PrintData() {
	fmt.Println("待分配进程：")
	for _, i := range system.P.Data {
		fmt.Printf("进程名：%s，所需空间：%d\n", i.Name, i.Need)
	}
	fmt.Println()
	fmt.Println("待分配空间：")
	for _, i := range system.S.Data {
		fmt.Printf("序号：%d，空间：%d\n", i.Sequence, i.Partition)
	}
	fmt.Println()
}

// 输出结果
func Result(s System, msg string) {
	fmt.Printf("分配状态：%s\n\n", msg)
	fmt.Println("分配进程：")
	for _, i := range s.P.Data {
		fmt.Printf("进程名：%s，所需空间：%d，分得位置：%d\n", i.Name, i.Need, i.AllocatePosition)
	}
	fmt.Println()
	fmt.Println("分配空间：")
	sort.Sort(BySequence(s.S.Data))
	for _, i := range s.S.Data {
		fmt.Printf("序号：%d，剩余空间：%d\n", i.Sequence, i.Partition)
	}
	fmt.Println()
}

// Print break-line
func PrintLine() {
	fmt.Println("------------------------------------------------------")
}

// 获取输入
func GetInput(system System) string {
	system.PrintData()
	fmt.Printf("请输入调度算法（1-FirstFit，2-NextFit，3-BestFit，4-WorstFit，q-quit)：")
	var algorthmType string
	fmt.Scanf("%s\n", &algorthmType)
	return algorthmType
}

func main() {
	var system, s System
	var msg string
	system.GetData()
loop:
	for {
		switch algorthmType := GetInput(system); algorthmType {
		case "1":
			s, msg = system.FirstFit()
		case "2":
			s, msg = system.NextFit()
		case "3":
			s, msg = system.BestFit()
		case "4":
			s, msg = system.WorstFit()
		case "q":
			break loop
		}
		Result(s, msg)
		PrintLine()
	}
}
