package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type docJSON struct {
	Dt      string `json:"time"`
	LossUDP int    `json:"UDP LOSS"`
}

var nodeUDP = make(map[string]int)
var gologer = sunLog()
var pushBuff int = 5

func main() {
	fmt.Println(time.Now(), "-- START")
	fmt.Println("version: beta-0.6.0")
	runtime.GOMAXPROCS(useableCPUNum()) //配置程序可用cpu数量
	var dataChannel = make(chan string, 1000)
	setModel := configDataSource()
	go runningModel(setModel, dataChannel)
	time.Sleep(10e8)
	dealData(dataChannel)

}
func runningModel(setModel string, ch chan string) {
	if setModel == "kafka" {
		topics := topicInfo()
		for topicName, settingInfo := range topics {
			fmt.Println(topicName, settingInfo.PushBuffNum)
			pushBuff = settingInfo.PushBuffNum
			k := analysisTopicGroup(topicName)
			receiveTopicData(k, ch)
		}

	} else {
		fileModel(setModel, "test", ch)
	}
}

func dealData(ch chan string) {
	var indexData toESData
	var docData docJSON
	var nodeName string
	for data := range ch {
		nodeName, docData.Dt, docData.LossUDP = analysisNodeUDP(data)
		if nodeName != "" {
			indexData.types, indexData.doc = nodeName, docData
			indexData.writeToESData()
		}
	}
}
func analysisNodeUDP(s string) (string, string, int) {
	var nodeName, dt string
	var udpNum int
	if strings.Contains(s, "InDatagrams") == false {
		sl := strings.Fields(s)
		nodeName = sl[0]
		dataDate, dataTime, UDPCount := sl[1], sl[2], sl[6]
		// fmt.Println(nodeName, dataDate, dataTime, UDPCount)
		udpNum = udpLoss(nodeName, UDPCount)
		dt = dataDate + "T" + dataTime + "+0800"
	}
	return nodeName, dt, udpNum
}
func udpLoss(nodeName, str string) int {
	i, err := strconv.Atoi(str)
	checkErr("string to int", err)
	j := nodeUDP[nodeName] //获取节点历史数据
	v := i - j
	nodeUDP[nodeName] = i //更新节点最新数据
	return v
}

// 从配置文件中读取kafka消费者信息
func analysisTopicGroup(topicName string) *kafkaConfig {
	k := &kafkaConfig{[]string{topicName}, topicName}
	return k
}
