// times insertion of docs to mongo
//
// inserts 3 docs with set names, random values at a time
// times map reduce to count set names
//
// Output: 
//Map Reduce Run #: 2
//#####################
//Doc count: {CalvinandHobbes 9}
//Doc count: {Facebook 1530917}
//Doc count: {Instagram 1530917}
//Doc count: {Twitter 1530917}
//Min Insert TIme 0
//Max Insert Time 0.43230098
//Avg Insert Time 0.06490538
//Min MapReduce TIme 0
//Max MapReduce Time 686.86914
//Avg MapReduce Time 188.86433
//Query One Time 2.365244
//Create Index Time 40.45919
//Drop Index Time 0.922316
//{MasterConns:38 SlaveConns:-22 SentOps:3123 ReceivedOps:1622 ReceivedDocs:1622 SocketsAlive:16 SocketsInUse:1 SocketRefs:1}

package main

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"math/rand"
	"time"
)

const Insert_Count = 10
const Async_Count = 15
const MapReduce_Count = 2
const FindOne_Count = 10

var Insert_t = make([]float32, Insert_Count*Async_Count)
var Mp_t = make([]float32, MapReduce_Count)
var Fo_t float32
var Id_t float32
var Di_t float32

type Mongodoc struct {
	Uid    int
	Name   string
	Score  int64
	Pubkey int64
}
type Mongostats struct {
	Mininsert int
	Maxinsert int
	Avginsert int
	Minget    int
	Maxget    int
	Avgget    int
}

// create random int64
func random() int64 {
	rand.Seed(time.Now().Unix() + rand.Int63() + 1)
	return rand.Int63()
}

// defer starttimer
func starttimer() int64 {
	return time.Now().UnixNano()
}

// defer endtimer
func endtimer(startTime int64, i int) {
	endTime := time.Now().UnixNano()
	Insert_t = append(Insert_t, float32(endTime-startTime)/1E9)
}

// defer map reduce timer
func mrendtimer(startTime int64) {
	endTime := time.Now().UnixNano()
	Mp_t = append(Mp_t, float32(endTime-startTime)/1E9)
}

// find one timer
func foendtimer(startTime int64) {
	endTime := time.Now().UnixNano()
	Fo_t = float32(endTime-startTime) / 1E9
}

//index timer
func iendtimer(startTime int64) {
	endTime := time.Now().UnixNano()
	Id_t = float32(endTime-startTime) / 1E9
}

// delete index timer
func dendtimer(startTime int64) {
	endTime := time.Now().UnixNano()
	Di_t = float32(endTime-startTime) / 1E9
}

// insert one doc
func insertonedoc(c *mgo.Collection) (err error) {
	//defer nmendtimer(starttimer())
	err = c.Insert(&Mongodoc{123456789, "CalvinandHobbes", 123456789, 123456789})
	return
}

// insert docs into mongo
func insertdoc(c *mgo.Collection, i int) (err error) {
	defer endtimer(starttimer(), i)
	err = c.Insert(&Mongodoc{i, "Facebook", random(), random()}, &Mongodoc{i, "Twitter", random(), random()}, &Mongodoc{i, "Instagram", random(), random()})
	return
}
func findone(c *mgo.Collection) (err error) {
	result := Mongodoc{}
	defer foendtimer(starttimer())
	err = c.Find(bson.M{"name": "CalvinandHobbes"}).One(&result)
	if err != nil {
		fmt.Println("not found")
		panic(err)
	}
	return
}

// count all Monogodoc Name fields
func testmapreduce(c *mgo.Collection, cm chan int, b int) (err error) {
	defer mrendtimer(starttimer())
	job := &mgo.MapReduce{
		Map:    "function() { emit(this.name, 1) }",
		Reduce: "function(key, values) { return Array.sum(values) }",
	}
	var result []struct {
		Id    string "_id"
		Value int
	}
	_, err = c.Find(nil).MapReduce(job, &result)
	if err != nil {
		return err
	}
	fmt.Printf("Map Reduce Run #: %v\n", b)
	fmt.Printf("#####################\n")
	for _, item := range result {
		fmt.Printf("Doc count: %v\n", item)
	}
	cm <- 1
	return
}

// with a copy of mongo session, call insert doc for specified insert count
func insertworker(s *mgo.Session, err error, ch chan int) {
	c := s.DB("test").C("mongotest")
	defer s.Close()
	for i := 0; i < Insert_Count; i++ {
		err = insertdoc(c, i)
		//insertdoc(c, i)
	}
	if err != nil {
		panic(err)
	}
	ch <- 1
}

// create index
func MakeIndex(c *mgo.Collection) {
	defer iendtimer(starttimer())
	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     false,
		DropDups:   false,
		Background: false,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

// delete index
func DeleteIndex(c *mgo.Collection) {
	defer dendtimer(starttimer())
	err := c.DropIndex("name")
	if err != nil {
		panic(err)
	}
}

// find avg insert time
func AvgFloat(avg []float32) (r float32) {
	var sum float32
	for i := 0; i < len(avg); i++ {
		sum += avg[i]
	}
	r = sum / float32(len(avg))
	return
}

// find min insert time
func MinFloat(min []float32) (r float32) {
	if len(min) > 0 {
		r = min[0]
	}
	for i := 1; i < len(min); i++ {
		if min[i] < r {
			r = min[i]
		}
	}
	return
}

// find max insert time
func MaxFloat(max []float32) (r float32) {
	if len(max) > 0 {
		r = max[0]
	}
	for i := 1; i < len(max); i++ {
		if max[i] > r {
			r = max[i]
		}
	}
	return
}

func main() {
	mgo.SetStats(true)
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	// go routine to start inserts
	c := session.DB("test").C("mongotest")
	insertonedoc(c)
	ch := make(chan int)
	for a := 0; a < 10; a++ {

		for j := 0; j < Async_Count; j++ {
			go insertworker(session.Copy(), err, ch)
		}
		// drain the channel
		for i := 0; i < Async_Count; i++ {
			<-ch
		}
	}
	// try to do query as a test
	findone(c)
	cm := make(chan int)
	for b := 0; b < MapReduce_Count; b++ {
		go testmapreduce(c, cm, b)
	}
	for z := 0; z < MapReduce_Count; z++ {
		<-cm
	}
	MakeIndex(c)
	DeleteIndex(c)
	minInsert := MinFloat(Insert_t[:])
	maxInsert := MaxFloat(Insert_t[:])
	avgInsert := AvgFloat(Insert_t[:])
	minMap := MinFloat(Mp_t[:])
	maxMap := MaxFloat(Mp_t[:])
	avgMap := AvgFloat(Mp_t[:])
	fmt.Printf("Min Insert TIme %v\n", minInsert)
	fmt.Printf("Max Insert Time %v\n", maxInsert)
	fmt.Printf("Avg Insert Time %v\n", avgInsert)
	fmt.Printf("Min MapReduce TIme %v\n", minMap)
	fmt.Printf("Max MapReduce Time %v\n", maxMap)
	fmt.Printf("Avg MapReduce Time %v\n", avgMap)
	fmt.Printf("Query One Time %v\n", Fo_t)
	fmt.Printf("Create Index Time (Blocking) %v\n", Id_t)
	fmt.Printf("Drop Index Time %v\n", Di_t)
	fmt.Printf("%+v\n", mgo.GetStats())

}
