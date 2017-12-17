package main

import (
	"gopkg.in/mgo.v2"
	"github.com/go-redis/redis"
	"sync"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

func main() {

	var counter int64
	session, err := mgo.Dial("10.15.0.145")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("userlist").C("newuserdata")

	var result UserInfo

	//err = c.Find(nil).All(&results)
	//if err != nil {
	//	fmt.Println("Not able to get the records from the db",err)
	//}

	find := c.Find(bson.M{})
	items := find.Iter()
	for items.Next(&result) {
		counter++
		fmt.Println(result.UserData.UID)
		fmt.Println(result.UserData.Msisdn)
		GetRedisInstanceGCP().Set("um:"+result.UserData.UID,result.UserData.Msisdn,0)
		fmt.Println("Migrated records till now --- >",counter)
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