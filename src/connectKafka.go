//连接，接受 kafka 数据
package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Shopify/sarama"
	"github.com/wvanbergen/kafka/consumergroup"
	kazoo "github.com/wvanbergen/kazoo-go"
)

type kafkaConfig struct {
	topicNames        []string
	consumerGroupName string
}

var (
	zookeeper      = zookeeperAddr()
	zookeeperNodes []string
)

// 从每个topic 接受数据
func receiveTopicData(k *kafkaConfig, ch chan string) {
	config := consumergroup.NewConfig()
	config.Offsets.Initial = sarama.OffsetNewest
	config.Offsets.ProcessingTimeout = 30 * time.Second
	zookeeperNodes, config.Zookeeper.Chroot = kazoo.ParseConnectionString(*zookeeper)
	fmt.Println(k.consumerGroupName, k.topicNames, zookeeperNodes)
	consumer, consumerErr := consumergroup.JoinConsumerGroup(k.consumerGroupName, k.topicNames, zookeeperNodes, config)
	if consumerErr != nil {
		gologer.Fatalln(consumerErr)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		if err := consumer.Close(); err != nil {
			sarama.Logger.Println("Error closing the consumer", err)
		}
	}()
	go func() {
		for err := range consumer.Errors() {
			gologer.Println(err)
		}
	}()
	offsets := make(map[string]map[int32]int64)
	gologer.Println(consumer.InstanceRegistered())
	for message := range consumer.Messages() {
		if offsets[message.Topic] == nil {
			offsets[message.Topic] = make(map[int32]int64)
		}
		ch <- string(message.Value)
		if offsets[message.Topic][message.Partition] != 0 && offsets[message.Topic][message.Partition] != message.Offset-1 {
			gologer.Printf("Unexpected offset on %s:%d. Expected %d, found %d, diff %d.\n", message.Topic, message.Partition, offsets[message.Topic][message.Partition]+1, message.Offset, message.Offset-offsets[message.Topic][message.Partition]+1)
		}
		offsets[message.Topic][message.Partition] = message.Offset
		consumer.CommitUpto(message)
	}
	gologer.Printf("%+v", offsets)
}
