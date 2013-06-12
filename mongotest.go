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
//Min Create Index Time (Blocking) 0
//Max Create Index Time (Blocking) 212.57533
//Avg Create Index Time (Blocking) 103.97096
//Min Drop Index Time 0
//Max Drop Index Time 0.74884
//Avg Drop Index Time 0.24233687

//{MasterConns:38 SlaveConns:-22 SentOps:3123 ReceivedOps:1622 ReceivedDocs:1622 SocketsAlive:16 SocketsInUse:1 SocketRefs:1}

package main

import (
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"math/rand"
	"time"
)
//how many times to run insert test
const Test_Runs = 2

//how many inserts a worker should do
const Insert_Count = 10

//how many concurrent workers should insert
const Async_Count = 20

//how many map reduce calls
const MapReduce_Count = 2

//not used
const FindOne_Count = 10

//how many times to create/drop index
const IndexOp_Count = 4

var Insert_t = make([]float32, Insert_Count*Async_Count)
var Mp_t = make([]float32, MapReduce_Count)
var Id_t = make([]float32, IndexOp_Count)
var Di_t = make([]float32, IndexOp_Count)
var Fo_t float32

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

// random creates random int64
func random() int64 {
	rand.Seed(time.Now().Unix() + rand.Int63() + 1)
	return rand.Int63()
}

// starttimer is a general starttimer
func starttimer() int64 {
	return time.Now().UnixNano()
}

// endtimer is an general endtimer
func endtimer(startTime int64, i int) {
	endTime := time.Now().UnixNano()
	Insert_t = append(Insert_t, float32(endTime-startTime)/1E9)
}

// mrendtimer is a map reduce timer
func mrendtimer(startTime int64) {
	endTime := time.Now().UnixNano()
	Mp_t = append(Mp_t, float32(endTime-startTime)/1E9)
}

// foendtimer is findone func timer
func foendtimer(startTime int64) {
	endTime := time.Now().UnixNano()
	Fo_t = float32(endTime-startTime) / 1E9
}

//iendtimer is a create index timer
func iendtimer(startTime int64) {
	endTime := time.Now().UnixNano()
	Id_t = append(Id_t, float32(endTime-startTime) / 1E9)
}

// dendtimer is a delete index timer
func dendtimer(startTime int64) {
	endTime := time.Now().UnixNano()
	Di_t = append(Di_t, float32(endTime-startTime) / 1E9)
}

// insertonedoc inserts one doc
func insertonedoc(c *mgo.Collection) error {
	//defer nmendtimer(starttimer())
	err := c.Insert(&Mongodoc{123456789, "CalvinandHobbes", 123456789, 123456789})
	return err
}

// insertdoc inserts docs into mongo, called from insertworker
func insertdoc(c *mgo.Collection, i int) error {
	defer endtimer(starttimer(), i)
	err := c.Insert(&Mongodoc{i, "Facebook", random(), random()}, &Mongodoc{i, "Twitter", random(), random()}, &Mongodoc{i, "Instagram", random(), random()})
	return err
}

// findone calls an mgo Find for one doc
func findone(c *mgo.Collection) error {
	result := Mongodoc{}
	defer foendtimer(starttimer())
	err := c.Find(bson.M{"name": "CalvinandHobbes"}).One(&result)
	if err != nil {
		fmt.Println("not found")
		panic(err)
	}
	return err
}

// testmapreduce counts all Monogodoc name fields
func testmapreduce(c *mgo.Collection, cm chan int, b int) error {
	defer mrendtimer(starttimer())
	job := &mgo.MapReduce{
		Map:    "function() { emit(this.name, 1) }",
		Reduce: "function(key, values) { return Array.sum(values) }",
	}
	var result []struct {
		Id    string "_id"
		Value int
	}
	_, err := c.Find(nil).MapReduce(job, &result)
	if err != nil {
		return err
	}
	fmt.Printf("Map Reduce Run #: %v\n", b)
	fmt.Printf("#####################\n")
	for _, item := range result {
		fmt.Printf("Doc count: %v\n", item)
	}
	cm <- 1
	return err
}

// insertworker calls func insertdoc for specified insert count 
func insertworker(s *mgo.Session, ch chan int) {
	c := s.DB("test").C("mongotest")
	defer s.Close()
	for i := 0; i < Insert_Count; i++ {
		err := insertdoc(c, i)
		if err != nil {
			panic(err)
		}
	}
	ch <- 1
}

// MakeIndex creates index on given collection
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

// DeleteIndex deletes a collections index
func DeleteIndex(c *mgo.Collection) {
	defer dendtimer(starttimer())
	err := c.DropIndex("name")
	if err != nil {
		panic(err)
	}
}

// AvgFloat finds avg insert time
func AvgFloat(avg []float32) float32 {
	var sum float32
	for i := 0; i < len(avg); i++ {
		sum += avg[i]
	}
	r := sum / float32(len(avg))
	return r
}

// MinFloat finds min insert time
func MinFloat(min []float32) float32 {
	r := min[0]
	for i := 1; i < len(min); i++ {
		if min[i] < r {
			r = min[i]
		}
	}
	return r
}

// MaxFloat finds max insert time
func MaxFloat(max []float32) float32 {
	r := max[0]
	for i := 1; i < len(max); i++ {
		if max[i] > r {
			r = max[i]
		}
	}
	return r
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
	for a := 0; a < Test_Runs; a++ {

		for j := 0; j < Async_Count; j++ {
			go insertworker(session.Copy(), ch)
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
	for x := 0; x < IndexOp_Count; x++ {
		MakeIndex(c)
		DeleteIndex(c)
	}
	minInsert := MinFloat(Insert_t[:])
	maxInsert := MaxFloat(Insert_t[:])
	avgInsert := AvgFloat(Insert_t[:])
	minMap := MinFloat(Mp_t[:])
	maxMap := MaxFloat(Mp_t[:])
	avgMap := AvgFloat(Mp_t[:])
	minIndexC := MinFloat(Id_t[:])
	maxIndexC := MaxFloat(Id_t[:])
	avgIndexC := AvgFloat(Id_t[:])
	minIndexD := MinFloat(Di_t[:])
	maxIndexD := MaxFloat(Di_t[:])
	avgIndexD := AvgFloat(Di_t[:])

	fmt.Printf("Min Insert TIme %v\n", minInsert)
	fmt.Printf("Max Insert Time %v\n", maxInsert)
	fmt.Printf("Avg Insert Time %v\n", avgInsert)
	fmt.Printf("Min MapReduce TIme %v\n", minMap)
	fmt.Printf("Max MapReduce Time %v\n", maxMap)
	fmt.Printf("Avg MapReduce Time %v\n", avgMap)
	fmt.Printf("Query One Time %v\n", Fo_t)
	fmt.Printf("Min Create Index Time (Blocking) %v\n", minIndexC)
	fmt.Printf("Max Create Index Time (Blocking) %v\n", maxIndexC)
	fmt.Printf("Avg Create Index Time (Blocking) %v\n", avgIndexC)
	fmt.Printf("Min Drop Index Time %v\n", minIndexD)
	fmt.Printf("Max Drop Index Time %v\n", maxIndexD)
	fmt.Printf("Avg Drop Index Time %v\n", avgIndexD)
	fmt.Printf("%+v\n", mgo.GetStats())

}
