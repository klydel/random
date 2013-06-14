package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

// how many times same src ip used for one worker conn
const connperworkerip = 1000
const numworkers = 100

func amiBorked(err error) {
	if err != nil {
		if err == io.EOF {
			return
		}
		fmt.Println("shit! ", err.Error())
			os.Exit(1)
	}
}

func httpDirectConnect(url *url.URL, cd chan int){
	client := &http.Client{}
	request, err := http.NewRequest("GET", url.String(), nil)
	amiBorked(err)
	dump, _ := httputil.DumpRequest(request, false)
	fmt.Println(string(dump))
	response, err := client.Do(request)
	amiBorked(err)
	if response.Status != "200 OK" {                                                                                                                                                                                                                                    
		fmt.Println(response.Status)
		fmt.Println(url)
		//os.Exit(2)
	}
		
	cd <- 1
}

func httpConnect(url *url.URL, srcurl *url.URL, c chan int) {
	transport := &http.Transport{Proxy: http.ProxyURL(srcurl)}
	client := &http.Client{Transport: transport}
	request, err := http.NewRequest("GET", url.String(), nil)
	amiBorked(err)
	dump, _ := httputil.DumpRequest(request, false)
	fmt.Println(string(dump))
	response, err := client.Do(request)
	amiBorked(err)
	if response.Status != "200 OK" {                                                                                                                                                                                                                                    
		fmt.Println(response.Status)
			fmt.Println(srcurl)
		os.Exit(2)
	}
		
	c <- 1
}

func httpConnectWorker(proxyURLS []string, url *url.URL, cg chan int) {
	ch := make(chan int)
	for i := range proxyURLS {
		srcURL := proxyURLS[i]
		srcurl, err := url.Parse(srcURL)
		amiBorked(err)
		for k := 0; k < connperworkerip; k++{
			go httpConnect(url, srcurl, ch)
		}
		for i := 0; i < connperworkerip; i++ {
                        <-ch
                }
		
	}
	cg <- 1
}
//http://64.78.30.180:80/
//bad: http://147.31.182.137:80/ http://64.78.30.180:80/ http://54.224.82.100:80/ http://147.31.182.137:80/ http://60.214.67.86:80/
func main() {
	proxyURLS := []string{"http://54.224.82.100:80/","http://64.78.30.180:80/","http://147.31.182.137:80/"}
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "http://www.offerpop.com/")
		os.Exit(1)
	}

	targetURL := os.Args[1]
	url, err := url.Parse(targetURL)
	amiBorked(err)
	cd := make(chan int)
	cg := make(chan int)
	for i := 0 ; i < numworkers; i++ {
		go httpDirectConnect(url, cd)
		go httpConnectWorker(proxyURLS[:], url, cg)
	}
	
	for i :=0; i < numworkers; i++ {
		<-cd
		<-cg
	}

	fmt.Println("Proxys: ", proxyURLS)	

	os.Exit(0)
}

