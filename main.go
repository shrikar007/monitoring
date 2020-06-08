package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/robfig/cron.v2"
	"log"
	"net/http"
	"runtime"
	"sync"
)
type Profile struct {
	Id     int       `gorm:"unique;not null;PRIMARY_KEY;AUTO_INCREMENT"`
	Name   string    `json:"name"`

}
type Profiles []Profile
var Flag bool


func main() {
	go memorycheck()
	http.HandleFunc("/", rungoroutine)
	http.ListenAndServe(":8080", nil)
}

func rungoroutine(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	off:=0
	Pro:=Profiles{}
	var data=make(chan interface{},10)
	db, err := gorm.Open("mysql", "root:root@tcp(192.168.43.106:3306)/channels")
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Profile{})

	for i:=0;i<=10;i++  {
		wg.Add(1)
		go getdata(&wg,db,off,data,Pro)
		off=off+5
		if Flag{
			break
		}
	}
	defer db.Close()
	wg.Wait()
	for i:=0;i<len(data);i++{
		fmt.Println(<-data)
	}
	fmt.Fprintf(w,"Done")
}

func getdata(wg *sync.WaitGroup,db *gorm.DB, off int,d chan interface{},pro Profiles){
	defer wg.Done()
	db.Limit(5).Offset(off).Find(&pro)
	if len(pro)==0{
		Flag=true
		return
	}
	d<-pro
}
func memorycheck(){
	c := cron.New()
	c.AddFunc("@every 0h0m20s", executememory)
	c.Start()
}
func executememory(){
	var mem runtime.MemStats
	fmt.Println("===============================================================")
	n:=runtime.NumGoroutine()
	fmt.Println("Number of Goroutines:",n)
	fmt.Println("memory ...")

	runtime.ReadMemStats(&mem)
	fmt.Println("Alloc memory:",mem.Alloc)
	fmt.Println("Total Alloc memory:",mem.TotalAlloc)
	fmt.Println("HeapAlloc memory:",mem.HeapAlloc)
	fmt.Println("HeapSys memory:",mem.HeapSys)

	cpu:=runtime.NumCPU()
	fmt.Println("CPU:",cpu)

}