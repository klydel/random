// times insertion of docs to mongo
//
// inserts 3 docs with set names, random values at a time
// times map reduce to count set names
//
// Output: 
//Doc count: {Facebook 2210473}
//Doc count: {Instagram 2210473}
//Doc Count: {Twitter 2210473}
//Min Insert TIme 0
//Max Insert Time 1.02038
//Avg Insert Time 0.00026085915
//Map Reduce Time 104.97748
//Query One Time 0
//Result: {0 Facebook 331870806444823685 9107841974541800341}
//{MasterConns:167 SlaveConns:-16 SentOps:3000221 ReceivedOps:1500221 ReceivedDocs:1500221 SocketsAlive:151 SocketsInUse:1 SocketRefs:1}

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
const MapReduce_Count = 10
const FindOne_Count = 10
var Insert_t = make([]float32, Insert_Count * Async_Count)
var Mp_t = make([]float32, MapReduce_Count)
var Fo_t float32
type Mongodoc struct {
	Uid int
        Name string
        Score int64
	Pubkey int64
}
type Mongostats struct {
	Mininsert int
	Maxinsert int
	Avginsert int
	Minget int
	Maxget int
	Avgget int
}
// create random int64
func random() (int64) {
	rand.Seed(time.Now().Unix() + rand.Int63() + 1)
	return rand.Int63()
}
// defer starttimer
func starttimer() (int64) {
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
	//Mp_t[] = float32(endTime-startTime)/1E9
	Mp_t = append(Mp_t, float32(endTime-startTime)/1E9)
}
// find one timer
func foendtimer(startTime int64) {
	endTime := time.Now().UnixNano()
	Fo_t = float32(endTime-startTime)/1E9
}
// insert one doc
func insertonedoc(c *mgo.Collection) (err error){
	//defer nmendtimer(starttimer())
	err = c.Insert(&Mongodoc{123456789, "CalvinandHobbes", 123456789, 123456789})
	return
}
// insert docs into mongo
func insertdoc(c *mgo.Collection, i int) (err error){
	defer endtimer(starttimer(), i)
	err = c.Insert(&Mongodoc{i, "Facebook", random(), random()},&Mongodoc{i, "Twitter", random(), random()},&Mongodoc{i, "Instagram", random(), random()})
	return
}
func findone(c *mgo.Collection) (err error){
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
func testmapreduce(c *mgo.Collection, cm chan int, b int) (err error){
	defer mrendtimer(starttimer())
	job := &mgo.MapReduce{
        Map:      "function() { emit(this.name, 1) }",
        Reduce:   "function(key, values) { return Array.sum(values) }",
	}
	var result []struct { Id string "_id"; Value int}
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
func insertworker(s *mgo.Session, err error, ch chan int) (){
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
// find avg insert time
func AvgFloat(avg []float32) (r float32) {
	var sum float32 
	for i := 0; i < len(avg); i++ {
		sum += avg[i]
	}
	r = sum/float32(len(avg))
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
	for a :=0; a < 10; a++ {
		
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
	for b:=0; b < MapReduce_Count; b++{
		go testmapreduce(c, cm, b)
	}
	for c := 0; c < MapReduce_Count; c++{
		<-cm
	}

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
	fmt.Printf("%+v\n", mgo.GetStats())

}
