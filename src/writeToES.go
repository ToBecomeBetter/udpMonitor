package main

import (
	"context"
	"fmt"
	"os"
	"regexp"

	elastic "gopkg.in/olivere/elastic.v5"
)

type toESData struct {
	types string
	doc   docJSON
}

var bulkRequest = connetEs() // es 连接地址
var ctx = context.Background()

func (d toESData) writeToESData() {
	if smallLetter(d.types) {
		fmt.Println("contain big letter !!!")
		os.Exit(1)
	}
	indexReq := elastic.NewBulkIndexRequest().Index("udp-monitor").Type(d.types).Doc(d.doc)
	bulkWriteToES(indexReq, bulkRequest)
}
func bulkWriteToES(indexReq *elastic.BulkIndexRequest, bulkRequest *elastic.BulkService) {
	bulkRequest = bulkRequest.Add(indexReq)
	if bulkRequest.NumberOfActions() >= pushBuff {
		_, err := bulkRequest.Do(ctx)
		checkErr("bulk write to es !", err)

	}
}

func smallLetter(s string) bool {
	reg := regexp.MustCompile(`[A-Z]`)
	bigLetter := reg.MatchString(s)
	return bigLetter
}
func connetEs() *elastic.BulkService {
	gologer.Println("连接elasticsearch地址为： ", esAddr())
	client, err := elastic.NewClient(elastic.SetURL(esAddr()))
	if err != nil {
		gologer.Println(err)
		// panic(err)
	}
	bulkRequest := client.Bulk()
	return bulkRequest

}
