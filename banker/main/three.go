package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Data struct {
	Available  []int64   // 可利用资源向量
	Allocation [][]int64 // 分配矩阵
	Need       [][]int64 // 需求矩阵
	Max        [][]int64 // 最大需求矩阵
}

// 初始化数据
func InitData() Data {
	available, allocation, need := GetData()
	max := make([][]int64, len(allocation))
	for i := range allocation {
		max[i] = make([]int64, len(allocation[i]))
	}
	for i := range allocation {
		for j := range allocation[i] {
			max[i][j] = allocation[i][j] + need[i][j]
		}
	}
	data := Data{
		Available:  available,
		Allocation: allocation,
		Need:       need,
		Max:        max,
	}
	return data
}

// 读取测试数据
func GetData() ([]int64, [][]int64, [][]int64) {
	var t []string
	var content []byte
	var available []int64
	var allocation, need [][]int64
	var err error
	var file *os.File

	file, err = os.Open("./available.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, _ = ioutil.ReadAll(file)
	for _, item_j := range strings.Split(string(content), ",") {
		f, err := strconv.ParseInt(item_j, 10, 64)
		if err != nil {
			panic(err)
		}
		available = append(available, f)
	}

	file, err = os.Open("./allocation.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, _ = ioutil.ReadAll(file)
	t = strings.Split(string(content), "\r\n")
	allocation = make([][]int64, len(t))
	for i, item_i := range t {
		for _, item_j := range strings.Split(item_i, ",") {
			f, err := strconv.ParseInt(item_j, 10, 64)
			if err != nil {
				panic(err)
			}
			allocation[i] = append(allocation[i], f)
		}
	}

	file, err = os.Open("./need.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, _ = ioutil.ReadAll(file)
	t = strings.Split(string(content), "\r\n")
	need = make([][]int64, len(t))
	for i, item_i := range t {
		for _, item_j := range strings.Split(item_i, ",") {
			f, err := strconv.ParseInt(item_j, 10, 64)
			if err != nil {
				panic(err)
			}
			need[i] = append(need[i], f)
		}
	}

	return available, allocation, need
}

// 安全性判断函数，可以判断当前时刻是否安全，以及判断输入请求向量后是否安全
func IsSafe(data Data, p int64, request []int64, useRequest bool) ([]int64, string, bool) {
	var safeSequence []int64
	t := 0
	finish := make([]bool, len(data.Allocation))
	// 初始化finish数组，已经执行完成的进程置为true
	for i := range finish {
		finish[i] = true
		for j := range data.Need[i] {
			if data.Need[i][j] == 0 {
				finish[i] = finish[i] && true
			} else {
				finish[i] = finish[i] && false
			}
		}
	}
	// 是否使用请求向量
	if useRequest {
		// 判断分配资源是否足够
		for i := range data.Available {
			if data.Available[i] < request[i] {
				return safeSequence, "可分配资源不足，请等待", false
			}
		}
		if data.ApplyRequest(p, request) {
			safeSequence = append(safeSequence, p)
			finish[p] = true
		}
	} else {
		p = 0
	}
	for i := p; true && t != len(data.Allocation); i++ {
		if IsAllTrue(finish) { // 当所有进程执行完break
			break
		}
		data.Collect(i)     // 当可以分配资源时，分配资源
		if data.IsDone(i) { // 当资源使用完毕，执行以下步骤
			for k := range data.Allocation[i] {
				data.Available[k] += data.Allocation[i][k]
				data.Allocation[i][k] = 0
				data.Need[i][k] = 0
			}
			safeSequence = append(safeSequence, i)
			finish[i] = true
			// 每当回收一次资源，将i和t置为最初始状态
			i = -1
			t = -1
		}
		if int(i) == len(data.Allocation)-1 {
			i = 0
		}
		t += 1
	}
	// 当遍历一遍数组，仍无法分配资源，则表明已经造成死锁
	if t == len(data.Allocation) {
		return safeSequence, "不安全, 会导致死锁", false
	}
	return safeSequence, "安全", true
}

// 判断数组是否全为true
func IsAllTrue(d []bool) bool {
	for _, i := range d {
		if !i {
			return false
		}
	}
	return true
}

// 判断进程是否得到足够资源
func (data *Data) IsDone(p int64) bool {
	for j := range data.Allocation[p] {
		if data.Allocation[p][j] == data.Max[p][j] {
			continue
		} else {
			return false
		}
	}
	return true
}

// 足够分配资源时，分配资源
func (data *Data) Collect(p int64) {
	flag := true
	for j := range data.Available {
		if data.Need[p][j] <= data.Available[j] {
			continue
		} else {
			flag = false
			break
		}
	}
	if flag {
		for j := range data.Available {
			data.Allocation[p][j] += data.Need[p][j]
			data.Available[j] -= data.Need[p][j]
			data.Need[p][j] = 0
		}
	}
}

// 往数据中加入请求变量
func (data *Data) ApplyRequest(p int64, request []int64) bool {
	for i := range request {
		data.Available[i] -= request[i]
		data.Allocation[p][i] += request[i]
		data.Need[p][i] -= request[i]
	}
	if data.IsDone(p) {
		for k := range data.Allocation[p] {
			data.Available[k] += data.Allocation[p][k]
			data.Allocation[p][k] = 0
			data.Need[p][k] = 0
		}
		return true
	}
	return false
}

// 对数据进行深拷贝
func deepCopy(src Data) Data {
	var dst = new(Data)
	b, _ := json.Marshal(src)
	json.Unmarshal(b, &dst)
	return Data{
		Available:  dst.Available,
		Allocation: dst.Allocation,
		Need:       dst.Need,
		Max:        dst.Max,
	}
}

// 将字符数组转为int数组
func stoi(request string) []int64 {
	r := []int64{}
	for _, item := range strings.Split(request, ",") {
		f, _ := strconv.ParseInt(item, 10, 64)
		r = append(r, f)
	}
	return r
}

// 打印table
func (data *Data) PrintTable() {
	fmt.Println("-------------------------------------------------------------")
	fmt.Printf("Allocation: ")
	fmt.Println(data.Allocation)
	fmt.Printf("Need:       ")
	fmt.Println(data.Need)
	fmt.Printf("Available:   ")
	fmt.Println(data.Available)
	fmt.Println("输入 -1 1 以测试当前时刻安全性")
	fmt.Printf("共有%d个进程，%d类资源\n输入进程号(%d-%d)，请求向量%d维(1,1,...)：", len(data.Allocation), len(data.Available), 0, len(data.Allocation), len(data.Available))
}

// 执行主函数
func main() {
	data := InitData()
	var num int64
	var request string
	var t bool
	for {
		data.PrintTable()
		fmt.Scanf("%d %s\n", &num, &request)
		if num == -1 {
			t = false
		} else {
			t = true
		}
		safeSequence, content, flag := IsSafe(deepCopy(data), num, stoi(request), t)
		fmt.Println(content)
		if flag {
			if t {
				data.ApplyRequest(num, stoi(request))
			}
			fmt.Printf("其中一种安全序列是 ")
			fmt.Println(safeSequence)
		}
		fmt.Printf("\n\n")
	}
}
