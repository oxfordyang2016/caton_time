package main

//about this server dev ,i want to say that
//ts is a central mannage unit,and it control database,tn and other config
//tn's reponsibility is to download and upload
//use comment to this file to make greate influence
import (
	_ "./api/"
	//i guess it is another_name,you can derectly use api package's stuff
	//if ti is . "api",that is  load all things that donot need somethings sunch as api.something
	"./db" //current dir
	"./pkg"
	pkg_model "./pkg/models" //this another name
	trans "./transfer"
	trans_model "./transfer/models"
	"./transfer/task"
	"github.com/astaxie/beego"
	clog "github.com/cihub/seelog"
)

//============================log block=================================================================>
func initLog() {
	cfgfiles := []string{
		"/opt/cydex/etc/ts_seelog.xml", //it work at after make,in linux mkdir
		"seelog.xml",
	} //slice

	//deal with log block.but  i have no time to figure out it
	for _, file := range cfgfiles {
		logger, err := clog.LoggerFromConfigAsFile(file)
		if err != nil {
			println(err.Error())
			continue
		}
		clog.ReplaceLogger(logger)
		break
	}
}

//============================================database=================================================>
// const DB = ":memory:"

const DB = "/tmp/cydex.sqlite3"

func initDB() (err error) {
	// make a databse
	clog.Info("init db")
	//create a database
	if err = db.CreateEngine("sqlite3", DB, false); err != nil {
		return //if expression meet standard,function stop
	}
	// 同步表结构
	if err = pkg_model.SyncTables(); err != nil { //from pkg/models
		return //have no return value,do nothing
	}
	if err = trans_model.SyncTables(); err != nil {
		return
	}
	return
}

//======================================start server=====================================================>
func start() {
	initLog()
	err := initDB()
	if err != nil {
		clog.Criticalf("Init DB failed: %s", err)
		panic("Shutdown")
	}
	// 设置拆包器
	pkg.SetUnpacker(pkg.NewDefaultUnpacker(50*1024*1024, 25))
	// 从数据库导入track
	pkg.JobMgr.LoadTracks()
	// set scheduler
	task.TaskMgr.SetScheduler(task.NewDefaultScheduler())
	// listen task state
	task.TaskMgr.ListenTaskState()
	// add job listen
	task.TaskMgr.AddObserver(pkg.JobMgr)
}

//================================ws server setting============================================================================
func run_ws() { //run  ws server
	ws_service := trans.NewWSServer("/ts", 12345, nil) //trans "./transfer"
	ws_service.SetVersion("1.0.0")
	ws_service.Serve()
}

//=====================================================main block============================================>
/*=====goroutines====
package main
import "fmt"
func f(from string) {
    for i := 0; i < 3; i++ {
        fmt.Println(from, ":", i)
    }
}

func main() {
f("direct")

go f("goroutine")
go func(msg string) {
        fmt.Println(msg)
    }("going")//pass agrv

    var input string
    fmt.Scanln(&input)
    fmt.Println("done")
}

$ go run goroutines.go
direct : 0
direct : 1
direct : 2
goroutine : 0
going
goroutine : 1
goroutine : 2
<enter>
done



*/
func main() {
	start()
	clog.Info("start")
	go run_ws()
	beego.Run(":8088") //this from beego project
}