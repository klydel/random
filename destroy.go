package main
//	"net/http/httputil"
import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// how many times same src ip used for one worker conn
const connperworkerip = 10
const numworkers = 10
var numget int64
//var lastget = make([]float32, numworkers*connperworkerip)
var lastget = make([]float32, 0)
func amiBorked(err error) bool{
	if err != nil {
		if err == io.EOF {
			return true
		}
		fmt.Println("shit! ", err.Error())
			//os.Exit(1)
		return true
	}
	return false
}

func starttimer() int64 {
        return time.Now().UnixNano()
}
func endtimer(startTime int64) {
        endTime := time.Now().UnixNano()
	lastget = append(lastget, float32(endTime-startTime)/1E9)
}
func httpDirectConnect(url *url.URL, cd chan int){
	for i := 0; i < connperworkerip; i++{
		numget += 1
		client := &http.Client{}
		request, err := http.NewRequest("GET", url.String(), nil)
		//amiBorked(err)
		defer endtimer(starttimer())
		response, err := client.Do(request)
		if amiBorked(err){
			return
		}
		if response.Status != "200 OK" {                                                                                                                                                                                                                                    			fmt.Println("Non 200 Status from : ", url)
				//os.Exit(2)
		}
	}	
	cd <- 1
}
func httpConnect(url *url.URL, srcurl *url.URL) {
	transport := &http.Transport{Proxy: http.ProxyURL(srcurl)}
	client := http.Client{Transport: transport}
	request, err := http.NewRequest("GET", url.String(), nil)
	amiBorked(err)
	response, err := client.Do(request)
	if amiBorked(err){
		return
	}
	if response.Status != "200 OK" {                                                                                                                                                                                                                                    
		fmt.Println(response.Status)
		fmt.Println("Non 200 Status from : ", srcurl)
	}
	response.Body.Close()

}
func httpConnectWorker(proxyURLS []string, url *url.URL, cc chan int) {
	for i := range proxyURLS {
		fmt.Println("Using Proxy: ",proxyURLS[i])
		srcURL := proxyURLS[i]
		srcurl, err := url.Parse(srcURL)
		amiBorked(err)
		for k := 0; k < connperworkerip; k++{
			httpConnect(url, srcurl)
			//fmt.Println(srcurl)
		}
		
	}
	cc <- 1
}
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "test")
	}

//http://64.78.30.180:80/
//bad: http://147.31.182.137:80/ http://64.78.30.180:80/ http://54.224.82.100:80/ http://147.31.182.137:80/ http://60.214.67.86:80/
//		"http://54.224.82.100:80/", "http://54.224.152.175:80/", "http://173.163.42.57:80/", "http://64.78.30.180:80/", "http://54.224.82.100:80/"}
func main() {
	http.HandleFunc("/", handler)
	go http.ListenAndServe(":8080", nil)
	proxyURLS := []string{"http://54.224.152.175:80/", "http://173.163.42.57:80/", "http://54.224.82.100:80/"}
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "localhost/test")
		os.Exit(1)
	}

	targetURL := os.Args[1]
	url, err := url.Parse(targetURL)
	amiBorked(err)
	cd := make(chan int)
	//cc := make(chan int)
	for i := 0 ; i < numworkers; i++ {
		//go httpConnectWorker(proxyURLS[:], url, cc)
		go httpDirectConnect(url, cd)
	}
	
	for i :=0; i < numworkers ; i++ {
		//<-cc
		<-cd
	}
	fmt.Println("Proxys: ", proxyURLS)	
	fmt.Printf("Times: %+v", lastget)
	fmt.Printf("Gets: %+v", numget)
	os.Exit(0)
	//select{}
}

