package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/gin-gonic/gin"
	echov3 "github.com/labstack/echo"
	"github.com/tockins/fresh"
)

var port = 8080
var sleepTime = 0
var sleepTimeDuration time.Duration
var message = []byte("hello world")
var messageStr = "hello world"
var samplingPoint = 20 //seconds

// server [default] [10] [8080]
func main() {
	args := os.Args
	argsLen := len(args)
	webFramework := "default"
	if argsLen > 1 {
		webFramework = args[1]
	}
	if argsLen > 2 {
		sleepTime, _ = strconv.Atoi(args[2])
	}
	if argsLen > 3 {
		port, _ = strconv.Atoi(args[3])
	}
	if argsLen > 4 {
		samplingPoint, _ = strconv.Atoi(args[4])
	}
	sleepTimeDuration = time.Duration(sleepTime) * time.Millisecond
	samplingPointDuration := time.Duration(samplingPoint) * time.Second

	go func() {
		time.Sleep(samplingPointDuration)
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		var u uint64 = 1024 * 1024
		fmt.Printf("TotalAlloc: %d\n", mem.TotalAlloc/u)
		fmt.Printf("Alloc: %d\n", mem.Alloc/u)
		fmt.Printf("HeapAlloc: %d\n", mem.HeapAlloc/u)
		fmt.Printf("HeapSys: %d\n", mem.HeapSys/u)
	}()

	switch webFramework {

	case "beego":
		startBeego()
	case "echov3":
		startEchoV3()
	case "fresh":
		startFresh()
	case "gin":
		startGin()
	}
}

//beego
func beegoHandler(ctx *context.Context) {
	if sleepTime > 0 {
		time.Sleep(sleepTimeDuration)
	} else {
		runtime.Gosched()
	}
	ctx.WriteString(messageStr)
}
func startBeego() {
	beego.BConfig.RunMode = beego.PROD
	beego.BeeLogger.Close()
	mux := beego.NewControllerRegister()
	mux.Get("/hello", beegoHandler)
	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}

// echov3-standard
func echov3Handler(c echov3.Context) error {
	if sleepTime > 0 {
		time.Sleep(sleepTimeDuration)
	} else {
		runtime.Gosched()
	}
	c.Response().Write(message)
	return nil
}
func startEchoV3() {
	e := echov3.New()
	e.GET("/hello", echov3Handler)

	e.Start(":" + strconv.Itoa(port))
}

//fresh
func freshHandler(c fresh.Context) error {
	if sleepTime > 0 {
		time.Sleep(sleepTimeDuration)
	} else {
		runtime.Gosched()
	}
	c.Response().Text(http.StatusOK, messageStr)
	return nil
}

func startFresh() {
	f := fresh.New()
	f.Config().SetPort(port)
	f.GET("/hello", freshHandler)
	f.Run()
}

// gin
func ginHandler(c *gin.Context) {
	if sleepTime > 0 {
		time.Sleep(sleepTimeDuration)
	} else {
		runtime.Gosched()
	}
	c.Writer.Write(message)
}
func startGin() {
	gin.SetMode(gin.ReleaseMode)
	mux := gin.New()
	mux.GET("/hello", ginHandler)
	mux.Run(":" + strconv.Itoa(port))
}

// mock
type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}
