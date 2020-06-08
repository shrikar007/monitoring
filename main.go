package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/robfig/cron.v2"
	"log"
	"net/http"
	"net/http/pprof"
	"runtime"
	"sync"
	"time"
)
type Profile struct {
	Id     int       `gorm:"unique;not null;PRIMARY_KEY;AUTO_INCREMENT"`
	Name   string    `json:"name"`

}
type Profiles []Profile
var Flag bool


func main() {
	go memorycheck()
	r := mux.NewRouter()
	r.HandleFunc("/", rungoroutine)

	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/block", pprof.Handler("block"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	http.ListenAndServe(":8080",r )
}

func rungoroutine(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup

	db, err := gorm.Open("mysql", "root:root@tcp(192.168.43.106:3306)/channels")
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Profile{})
	off:=0
	Pro:=Profiles{}
	data:=make(chan Profiles,10)
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

	//for x:=range data{
	//	fmt.Println(x)
	//}

	for j:=0;j<=len(data)+1;j++{
		fmt.Println(<-data)
	}
	fmt.Fprintf(w,"Done")
}

func getdata(wg *sync.WaitGroup,db *gorm.DB, off int,d chan Profiles,pro Profiles){
	defer wg.Done()
	db.Limit(5).Offset(off).Find(&pro)
	if len(pro)==0{
		Flag=true
		return
	}
	d<-pro
	time.Sleep(time.Second)
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