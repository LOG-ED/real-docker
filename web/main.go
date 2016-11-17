package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Ping struct {
	Id   bson.ObjectId `bson:"_id"`
	Time time.Time     `bson:"time"`
}

func (p Ping) String() string {
	// layouts must use the reference time Mon Jan 2 15:04:05 MST 2006
	return fmt.Sprintf("Document inserted at %v\n", p.Time.Format("3:04:05 PM"))
}

func main() {
	http.HandleFunc("/", list)
	http.HandleFunc("/new", add)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mongoConnect() (session *mgo.Session) {
	// get the session using information from the environment
	session, err := mgo.Dial(os.Getenv("MONGODB_ADDRESS"))

	// panics if connection error occurs
	if err != nil {
		panic(err)
	}

	return session
}

func list(w http.ResponseWriter, r *http.Request) {
	session := mongoConnect()
	defer session.Close()

	// get all documents
	pings := []Ping{}
	db := mongoConnect().DB(os.Getenv("MONGODB_NAME"))
	db.C("pings").Find(nil).All(&pings)

	fmt.Fprint(w, pings)
}

func add(w http.ResponseWriter, r *http.Request) {
	session := mongoConnect()
	defer session.Close()

	// insert new document
	ping := Ping{
		Id:   bson.NewObjectId(),
		Time: time.Now(),
	}

	db := mongoConnect().DB(os.Getenv("MONGODB_NAME"))
	db.C("pings").Insert(ping)

	fmt.Fprint(w, ping)
}
