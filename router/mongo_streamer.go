package router

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type EventLog struct {
	Id        bson.ObjectId `bson:"_id,omitempty"`
	Source    string        `bson:"source"`
	EventType string        `bson:"event_type"`
	When      time.Time     `bson:"when"`
	Whom      string        `bson:"whom"`
	Contents  string        `bson:"contents"`
}

func mongoStreamer(target Target, types []string, logstream chan *Log) {
	println("mongo stream line")
	session, err := mgo.Dial("qa1-mongo-radioedit101.qa.cloud.ihr:37017,qa1-mongo-radioedit102.qa.cloud.ihr:37017,qa1-mongo-radioedit103.qa.cloud.ihr:37017/infrastructure")
	session.SetMode(mgo.Monotonic, true)
	if err != nil {
		fmt.Print("MONGO DOWN, MONGO DOWN!!")
		panic(err)
	}
	defer session.Close()
	c := session.DB("infrastructure").C("event_log")

	for logline := range logstream {
		spew.Dump(logline)
		el := EventLog{
			Source:    "i-" + logline.ID[:8],
			EventType: "instance",
			Contents:  logline.Data,
			When:      time.Now(),
		}
		c.Insert(&el)
	}
}
