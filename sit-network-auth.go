package sit_network_auth_go

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	Client                   = &http.Client{}
	LoginFileName            = "log.log"
	StudentDataFileName      = "NetworkLoginUsers.csv"
	AvailableStuDataFileName = "availableData.csv"
	logFileRaw               *os.File
)

func init() {
	logFileRaw, err := os.OpenFile(LoginFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFileRaw)
	log.SetOutput(mw)
}

func splitString(r rune) bool {
	return r == '(' || r == ')'
}

func login(user string, passwd string) bool {
	var loginUrl = "http://172.16.8.70/drcom/login?callback=dr1003&" +
		"DDDDD=" + user +
		"&upass=" + passwd +
		"&0MKKey=123456&R1=0&R2=&R3=0&R6=0&para=00&v6ip=&terminal_type=1&lang=zh-cn&jsVersion=4.1.3&v=9779&lang=zh"
	// get response
	_, err := Client.Get(loginUrl)
	if err != nil {

		log.Println(err)
		return false
	}
	log.Println("id:" + user + " network connect")

	return true

}

func checkNetwork() bool {
	var statusUrl = "http://172.16.8.70/drcom/chkstatus?callback=dr1002&jsVersion=4.1&v=7808&lang=zh"
	// get response
	resp, _ := Client.Get(statusUrl)
	// parse body
	body, _ := ioutil.ReadAll(resp.Body)
	var str = string(body)
	a := strings.FieldsFunc(str, splitString)
	myMap := make(map[string]interface{})
	// convert json
	err := json.Unmarshal([]byte(a[1]), &myMap)
	if err != nil {
		return false
	}
	// judge result==1
	var result, _ = strconv.Atoi(fmt.Sprintf("%v", myMap["result"]))
	return result == 1
}

func ReadCSV(fileName string) ([][]string, error) {
	opencast, err := os.Open(fileName)
	if err != nil {

		log.Println(StudentDataFileName + "打开失败")
		return nil, err
	}
	defer func(opencast *os.File) {
		err := opencast.Close()
		if err != nil {
			//log.Println(err)
		}
	}(opencast)

	Read := csv.NewReader(opencast)
	read, _ := Read.ReadAll() //返回切片类型：[chen  hai wei]
	log.Println("load " + StudentDataFileName + "success")

	return read, nil
}

func AddAvailableUse(user string, passwd string) error {
	opencast, err := os.OpenFile(AvailableStuDataFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {

		log.Println("load" + AvailableStuDataFileName + " fail")
		return err
	}
	defer func(opencast *os.File) {
		err := opencast.Close()
		if err != nil {
			//log.Println(err)
		}
	}(opencast)
	_, _ = opencast.Seek(0, io.SeekEnd)
	write := csv.NewWriter(opencast)
	row := []string{user, passwd}
	err = write.Write(row)
	if err != nil {

		log.Println("write new student data error")
		return err
	}
	write.Flush()
	fmt.Println("rtsp://user:6ad77fc8@10.1.160.241:{58000,50555}")
	return nil

}
