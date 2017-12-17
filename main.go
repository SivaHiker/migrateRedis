package main

import (
	"gopkg.in/mgo.v2"
	"github.com/go-redis/redis"
	"sync"
	"fmt"
)

func main() {

	session, err := mgo.Dial("10.15.0.145")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("userlist").C("newuserdata")

	var results []UserInfo

	err = c.Find(nil).All(&results)
	if err != nil {
		fmt.Println("Not able to get the records from the db",err)
	}

	for i := 0; i <=len(results);i++ {
		fmt.Println(results[i].UserData.UID)
		fmt.Println(results[i].UserData.Msisdn)
		GetRedisInstanceGCP().Set("um:"+results[i].UserData.UID,results[i].UserData.Msisdn,0)
		fmt.Println("Migrated records till now --- >",i)
	}
}


func GetRedisInstanceGCP() *redis.Client {
	var onceGCP sync.Once
	var instanceGCP *redis.Client
	onceGCP.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     "redis-13511.c2.asia-southeast-1-1.gce.cloud.redislabs.com:13511",
			//Addr:     "localhost:6379",
			Password: "",
			DB:       0,
			PoolSize: 100,
			/*
					PoolTimeout:  10 * time.Minute,
					IdleTimeout:  5 * time.Minute,
					ReadTimeout:  2 * time.Second,
					WriteTimeout: 10 * time.Second,
			*/
		})
		instanceGCP = client
	})
	return instanceGCP
}

type UserInfo struct {
	UserData UserData `json:"UserData"`
	Flag bool `json:"flag"`
}

type UserData struct {
	Msisdn string `json:"msisdn"`
	Token  string `json:"token"`
	UID    string `json:"uid"`
	PlatformUID string `json:"platformuid"`
	PlatformToken string `json:"platformtoken"`
}