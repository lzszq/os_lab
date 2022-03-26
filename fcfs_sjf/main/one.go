package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Items struct {
	Data []Item
}

// 定义每个进程执行过程所产生数据
type Item struct {
	P                     Process
	FinishTime            float64 // 完成时间
	TurnoverTime          float64 // 周转时间
	TurnoverWithRightTime float64 // 带权周转时间
}

// 初始化一个空Item
func InitItem(p Process) Item {
	t := Item{
		P:                     p,
		FinishTime:            0,
		TurnoverTime:          0,
		TurnoverWithRightTime: 0,
	}
	return t
}

type Processs struct {
	Data []Process
}

// 定义进程
type Process struct {
	Name       string  // 进程名
	ArriveTime float64 // 到达时间
	ServeTime  float64 // 服务时间
	CountTime  float64 // 剩下服务时间
}

// 实现结构体排序接口
type ByArriveTime []Process

func (a ByArriveTime) Len() int           { return len(a) }
func (a ByArriveTime) Less(i, j int) bool { return a[i].ArriveTime < a[j].ArriveTime }
func (a ByArriveTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// 实现结构体排序接口
type ByServeTime []Process

func (a ByServeTime) Len() int           { return len(a) }
func (a ByServeTime) Less(i, j int) bool { return a[i].ServeTime < a[j].ServeTime }
func (a ByServeTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// 实现结构体排序接口
type ByProcessName []Item

func (a ByProcessName) Len() int           { return len(a) }
func (a ByProcessName) Less(i, j int) bool { return a[i].P.Name < a[j].P.Name }
func (a ByProcessName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// 添加队列尾元素
func (p *Processs) Push(d Process) {
	p.Data = append(p.Data, d)
}

// 取队列首元素
func (p *Processs) Pop() Process {
	t := p.Data[0]
	p.Data = p.Data[1:]
	return t
}

// 判断队列是否为空
func (p *Processs) IsEmpty() bool {
	return len(p.Data) == 0
}

// 读取测试数据
func GetData() Processs {
	file, err := os.Open("./data.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	content, _ := ioutil.ReadAll(file)
	t := strings.Split(string(content), "\r\n")
	names := strings.Split(t[0], ",")
	arriveTimes := []float64{}
	for _, i := range strings.Split(t[1], ",") {
		f, err := strconv.ParseFloat(i, 32)
		if err != nil {
			panic(err)
		}
		arriveTimes = append(arriveTimes, f)
	}
	serveTimes := []float64{}
	for _, i := range strings.Split(t[2], ",") {
		f, err := strconv.ParseFloat(i, 32)
		if err != nil {
			panic(err)
		}
		serveTimes = append(serveTimes, f)
	}
	data := Processs{}
	data.InsertData(names, arriveTimes, serveTimes)
	return data
}

// 生成测试数据
func (p *Processs) InsertData(names []string, arriveTimes, serveTimes []float64) {
	for index := range names {
		process := Process{
			Name:       names[index],
			ArriveTime: arriveTimes[index],
			ServeTime:  serveTimes[index],
			CountTime:  serveTimes[index],
		}
		p.Data = append(p.Data, process)
	}
}

// 添加Item
func (items *Items) AddItem(time float64, process Process) {
	items.Data = append(items.Data, InitItem(process))
	l := len(items.Data) - 1
	items.Data[l].FinishTime = time
	items.Data[l].TurnoverTime = time - items.Data[l].P.ArriveTime
	items.Data[l].TurnoverWithRightTime = items.Data[l].TurnoverTime / items.Data[l].P.ServeTime
}

// 模拟进程调度过程，每次进出队列计算时间，并根据算法对进程进行排序
func Algorithm(data Processs, algorithm string) Items {
	sort.Sort(ByArriveTime(data.Data)) // 保证数据以到达时间升序排列
	processs := Processs{}
	items := Items{}
	time := 0.0

	for !data.IsEmpty() || !processs.IsEmpty() {
		if !data.IsEmpty() {
			if data.Data[0].ArriveTime == time { // 按时刻将进程加入待执行队列
				front := data.Pop()
				processs.Push(front)
				fmt.Printf("到达进程为：%s\n", front.Name)
				if algorithm == "SJF" {
					sort.Sort(ByServeTime(processs.Data[1:])) // 对当前不在运行的进程按服务时间进行升序排列
				}
			}
		}
		fmt.Printf("当前时刻为：%.f-%.f, 当前运行进程为：%s\n", time, time+1, processs.Data[0].Name)
		time += 1
		if !processs.IsEmpty() {
			if processs.Data[0].CountTime > 0 {
				processs.Data[0].CountTime -= 1
				if processs.Data[0].CountTime == 0 {
					items.AddItem(time, processs.Pop()) // 服务完后，往结果队列添加进程
				}
			}
		}
	}
	sort.Sort(ByProcessName(items.Data))

	return items
}

// 获取平均时间
func GetAverageTime(items Items) (float64, float64) {
	time1, time2 := 0.0, 0.0
	for _, item := range items.Data {
		time1 += item.TurnoverTime
		time2 += item.TurnoverWithRightTime
	}
	l := float64(len(items.Data))
	return time1 / l, time2 / l
}

// add Pause to windows exe
func Pause() error {
	buf := make([]byte, 1)
	fmt.Println("\n输入 c ，然后回车以退出。")
	for {
		c, e := os.Stdin.Read(buf)
		if c != 1 {
			return e
		}
		if buf[0] == 'c' {
			break
		}
	}
	return nil
}

// Print break-line
func PrintLine() {
	fmt.Println("------------------------------------------------------")
}

func main() {
	algorthmType := ""
	for !(algorthmType == "FCFS" || algorthmType == "SJF") {
		fmt.Printf("请输入调度算法（FCFS, SJF)：")
		fmt.Scanf("%s", &algorthmType)
	}
	PrintLine()
	items := Algorithm(GetData(), algorthmType)
	PrintLine()
	fmt.Println("进程名 到达时间 服务时间 完成时间 周转时间 带权周转时间")
	for _, item := range items.Data {
		fmt.Printf("%s     %.2f     %.2f     %.2f     %.2f     %.2f\n", item.P.Name,
			item.P.ArriveTime, item.P.ServeTime, item.FinishTime, item.TurnoverTime,
			item.TurnoverWithRightTime)
	}
	time1, time2 := GetAverageTime(items)
	fmt.Printf("平均周转时间：%.2f 平均带权周转时间：%.2f\n", time1, time2)
	Pause()
}
