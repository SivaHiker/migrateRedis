package main

import (
	"gopkg.in/mgo.v2"
	"github.com/go-redis/redis"
	"sync"
	"fmt"
	//"gopkg.in/mgo.v2/bson"
)

func main() {

	var counter int64
	var activeCounter int64
	var inactiveCounter int64
	session, err := mgo.Dial("10.15.0.145")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("userlist").C("newuserdata")


	//fromsession, err := mgo.Dial("10.15.0.75")
	//if err != nil {
	//	panic(err)
	//}
	//defer fromsession.Close()
	//
	//c1 := fromsession.DB("userdb").C("users")

	var result []UserInfo

	err = c.Find(nil).All(&result)
	if err != nil {
		fmt.Println("Not able to get the records from the db",err)
	}
	//var resultSet UserRecord

	//find := c.Find(bson.M{})
	//items := find.Iter()
	//for items.Next(&result) {
	fmt.Println("Results count",len(result))
	for i:=0;i<len(result);i++ {
		counter++
		if(result[i].Active) {
			GetRedisInstanceGCP().Set("um:"+result[i].UserData.UID, result[i].UserData.Msisdn, 0)
			activeCounter++
		} else {
			GetRedisInstanceGCP().Set("ud:"+result[i].UserData.UID, result[i].UserData.Msisdn, 0)
			GetRedisInstanceGCP().Set("md:"+result[i].UserData.Msisdn,result[i].UserData.UID , 0)
			inactiveCounter++
		}

		//fmt.Println(resultSet.Status)
		//msisdnFilter := resultSet.Msisdn[0]
		//if(resultSet.Status==1 || resultSet.Status==2){
		//	erro:=c.Update(bson.M{"userdata.msisdn": msisdnFilter}, bson.M{"$set": bson.M{"active": false}})
		//	if(erro!=nil){
		//		fmt.Println("Not able to update the record",msisdnFilter)
		//	}
		//} else {
		//	erro:=c.Update(bson.M{"userdata.msisdn": msisdnFilter}, bson.M{"$set": bson.M{"active": true}})
		//	if(erro!=nil){
		//		fmt.Println("Not able to update the record",msisdnFilter)
		//	}
		//}
		//GetRedisInstanceGCP().Set("um:"+result.UserData.UID,result.UserData.Msisdn,0)
		fmt.Println("Migrated records till now --- >",counter)
	}
	fmt.Println("Total Active User records  --- >",activeCounter)
	fmt.Println("Total InActive User records  --- >",inactiveCounter)

}


func GetRedisInstanceGCP() *redis.Client {
	var onceGCP sync.Once
	var instanceGCP *redis.Client
	onceGCP.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     "redis-13065.c2.asia-southeast-1-1.gce.cloud.redislabs.com:13065",
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
	Active bool `json:"active"`
}

type UserData struct {
	Msisdn string `json:"msisdn"`
	Token  string `json:"token"`
	UID    string `json:"uid"`
	PlatformUID string `json:"platformuid"`
	PlatformToken string `json:"platformtoken"`
}



type UserRecord struct {
	_id           string        `json:"_id"`
	Addressbook   struct{}      `json:"addressbook"`
	BackupToken   string        `json:"backup_token"`
	Connect       int           `json:"connect"`
	Country       string        `json:"country"`
	Devices       []interface{} `json:"devices"`
	Gender        string        `json:"gender"`
	Icon          string        `json:"icon"`
	InvitedJoined []interface{} `json:"invited_joined"`
	Invitetoken   string        `json:"invitetoken"`
	Locale        string        `json:"locale"`
	Msisdn        []string      `json:"msisdn"`
	Name          string        `json:"name"`
	PaUID         string        `json:"pa_uid"`
	Referredby    []string      `json:"referredby"`
	RewardToken   string        `json:"reward_token"`
	Status        int           `json:"status"`
	Sus           int           `json:"sus"`
	Uls           int           `json:"uls"`
	Version       int           `json:"version"`
}
