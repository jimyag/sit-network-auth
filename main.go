package main

import (
	sitauth "github.com/sit-network-auth-go/lib"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	defer func(logFileRaw *os.File) {
		err := sitauth.LogFileRaw.Close()
		if err != nil {
		}
	}(sitauth.LogFileRaw)
	data, _ := sitauth.ReadCSV(sitauth.AvailableStuDataFileName)
	source := rand.NewSource(time.Now().Unix())
	myRand := rand.New(source) // initialize local pseudorandom generator
	randNum := myRand.Intn(len(data))
	var user []string
	for {
		if !sitauth.CheckNetwork() {
			randNum = myRand.Intn(len(data))
			user = data[randNum]
			if !sitauth.Login(user[0], user[1]) {
				log.Println(user[0] + " error")
			}
		}
		time.Sleep(time.Second * 3)
	}
}
