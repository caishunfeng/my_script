package mongo

import (
	"fmt"
	mgo "gopkg.in/mgo.v2"
	"sync"
)

var mgoClient MgoClient

type MgoClient struct {
	session  *mgo.Session
	database string
	lock     *sync.RWMutex
}

func InitMgoClient(adds string, database string) (err error) {
	if mgoClient.session == nil {
		session, err := mgo.Dial(adds)
		if err != nil {
			return fmt.Errorf("mongo连接失败，adds:%s,err:%v", adds, err)
		}
		mgoClient = MgoClient{
			session:  session,
			database: database,
		}
	}
	return
}

func CloseMgoClient() {
	if mgoClient.session != nil {
		mgoClient.session.Close()
	}
}

//获取集合
func GetMongoCollection(collectionName string) (*mgo.Collection, error) {
	if mgoClient.session == nil {
		return nil, fmt.Errorf("mgo client is nil")
	}
	db := mgoClient.session.DB(mgoClient.database)
	collection := db.C(collectionName)
	return collection, nil
}

//执行方法
func DoMongoFun(collectionName string, f func(*mgo.Collection) error) error {
	collection, err := GetMongoCollection(collectionName)
	if err != nil {
		return fmt.Errorf("在Mongo[%s][%s] 执行方法错误,原因:%s", mgoClient.database, collectionName, err)
	}
	return f(collection)
}
