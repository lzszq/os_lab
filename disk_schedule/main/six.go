package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Track struct {
	TrackSize        int
	TrackAccessQueue []int
	TrackStartNum    int
}

// 读取测试数据
func GetData(trackStartNum int) Track {
	var (
		file    *os.File
		err     error
		content []byte
	)

	track := Track{
		TrackSize:        0,
		TrackAccessQueue: nil,
		TrackStartNum:    trackStartNum,
	}

	file, err = os.Open("./data.txt")
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
		track.TrackAccessQueue = append(track.TrackAccessQueue, int(f))
	}
	track.TrackSize = len(track.TrackAccessQueue)
	return track
}

func Sum(data []int) int {
	result := int(0)
	for _, i := range data {
		result += i
	}
	return result
}

func Reverse(input []int) {
	inputLen := len(input)
	inputMid := inputLen / 2
	for i := 0; i < inputMid; i++ {
		j := inputLen - i - 1
		input[i], input[j] = input[j], input[i]
	}
}

func FCFS(track Track) {
	lastTrackNum := int(track.TrackStartNum)
	var moveDistance []int
	for _, i := range track.TrackAccessQueue {
		moveDistance = append(moveDistance, int(math.Abs(float64(i-int(lastTrackNum)))))
		lastTrackNum = int(i)
		fmt.Printf("被访问的下一磁道号为：%d，移动距离：%d\n", lastTrackNum, moveDistance[len(moveDistance)-1])
	}
	fmt.Printf("平均移动距离：%v\n\n", (float64(Sum(moveDistance)) / float64(track.TrackSize)))
}

func SSTF(track Track) {
	minElement := track.TrackStartNum
	elementIndex := -1
	var moveDistance []int
	for cnt := 0; cnt < track.TrackSize; cnt++ {
		minDistance := math.Inf(0)
	loop:
		for index, i := range track.TrackAccessQueue {
			if i == int(-1) {
				continue loop
			}
			t := math.Abs(float64(minElement - int(i)))
			if minDistance > t {
				minDistance = t
				minElement = int(i)
				elementIndex = index
			}
		}
		track.TrackAccessQueue[elementIndex] = -1
		moveDistance = append(moveDistance, int(minDistance))
		fmt.Printf("被访问的下一磁道号为：%d，移动距离：%d\n", minElement, moveDistance[len(moveDistance)-1])
	}
	fmt.Printf("平均移动距离：%v\n\n", (float64(Sum(moveDistance)) / float64(track.TrackSize)))
}

func SCAN(track Track, trackDirection int) {
	var newTrackAccessQueue []int
	sort.Ints(track.TrackAccessQueue)
	if trackDirection == 0 {
		Reverse(track.TrackAccessQueue)
	}

	elementIndex := -1
loop:
	for index, i := range track.TrackAccessQueue {
		if i > track.TrackStartNum && trackDirection == 1 {
			elementIndex = index
			break loop
		} else if i < track.TrackStartNum && trackDirection == 0 {
			elementIndex = index
			break loop
		}
	}
	t := track.TrackAccessQueue[:elementIndex]
	Reverse(t)
	newTrackAccessQueue = append(track.TrackAccessQueue[elementIndex:], t...)
	lastTrackNum := int(track.TrackStartNum)
	var moveDistance []int
	for _, i := range newTrackAccessQueue {
		moveDistance = append(moveDistance, int(math.Abs(float64(i-int(lastTrackNum)))))
		lastTrackNum = int(i)
		fmt.Printf("被访问的下一磁道号为：%d，移动距离：%d\n", lastTrackNum, moveDistance[len(moveDistance)-1])
	}
	fmt.Printf("平均移动距离：%v\n\n", (float64(Sum(moveDistance)) / float64(track.TrackSize)))
}

func CSCAN(track Track, trackDirection int) {
	var newTrackAccessQueue []int
	sort.Ints(track.TrackAccessQueue)
	if trackDirection == 0 {
		Reverse(track.TrackAccessQueue)
	}

	elementIndex := -1
loop:
	for index, i := range track.TrackAccessQueue {
		if i > track.TrackStartNum && trackDirection == 1 {
			elementIndex = index
			break loop
		} else if i < track.TrackStartNum && trackDirection == 0 {
			elementIndex = index
			break loop
		}
	}
	t := track.TrackAccessQueue[:elementIndex]
	newTrackAccessQueue = append(track.TrackAccessQueue[elementIndex:], t...)

	lastTrackNum := int(track.TrackStartNum)
	var moveDistance []int
	for _, i := range newTrackAccessQueue {
		moveDistance = append(moveDistance, int(math.Abs(float64(i-int(lastTrackNum)))))
		lastTrackNum = int(i)
		fmt.Printf("被访问的下一磁道号为：%d，移动距离：%d\n", lastTrackNum, moveDistance[len(moveDistance)-1])
	}
	fmt.Printf("平均移动距离：%v\n\n", (float64(Sum(moveDistance)) / float64(track.TrackSize)))
}

// 获取输入
func GetInput() string {
	data := GetData(0)
	fmt.Printf("访问序列：")
	fmt.Println(data.TrackAccessQueue)
	fmt.Printf("请输入调度算法（1-FCFS，2-SSTF，3-SCAN，4-CSCAN，q-quit)：")
	var algorthmType string
	fmt.Scanf("%s\n", &algorthmType)
	return algorthmType
}

// Print break-line
func PrintLine() {
	fmt.Println("------------------------------------------------------")
}

func main() {
loop:
	for {
		var trackStartNum int
		var trackDirection int
		switch algorthmType := GetInput(); algorthmType {
		case "1":
			fmt.Printf("请输入起始磁道号：")
			fmt.Scanf("%d\n", &trackStartNum)
			FCFS(GetData(trackStartNum))
		case "2":
			fmt.Printf("请输入起始磁道号：")
			fmt.Scanf("%d\n", &trackStartNum)
			SSTF(GetData(trackStartNum))
		case "3":
			fmt.Printf("请输入起始磁道号：")
			fmt.Scanf("%d\n", &trackStartNum)
			fmt.Printf("请输入方向（1-up，0-down）：")
			fmt.Scanf("%d\n", &trackDirection)
			SCAN(GetData(trackStartNum), trackDirection)
		case "4":
			fmt.Printf("请输入起始磁道号：")
			fmt.Scanf("%d\n", &trackStartNum)
			fmt.Printf("请输入方向（1-up，0-down）：")
			fmt.Scanf("%d\n", &trackDirection)
			CSCAN(GetData(trackStartNum), trackDirection)
		case "q":
			break loop
		}
		PrintLine()
	}
}
