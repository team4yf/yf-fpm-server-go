//Package fpm the core fpm
package fpm

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/gorilla/mux"
	"github.com/team4yf/yf-fpm-server-go/config"
	"github.com/team4yf/yf-fpm-server-go/ctx"
	"github.com/team4yf/yf-fpm-server-go/middleware"
	"github.com/team4yf/yf-fpm-server-go/pkg/cache"
	"github.com/team4yf/yf-fpm-server-go/pkg/log"
	"github.com/team4yf/yf-fpm-server-go/pkg/utils"
	"github.com/team4yf/yf-fpm-server-go/version"
)

var (
	registerEvents []HookHandler

	errNoMethod = errors.New("No method defined")

	defaultInstance *Fpm

	tempMapData map[string]interface{}
)

func init() {
	tempMapData = make(map[string]interface{})
}

//Register register some plugin
func Register(event HookHandler) {
	if len(registerEvents) < 1 {
		registerEvents = make([]HookHandler, 0)
	}
	registerEvents = append(registerEvents, event)
}

//Fpm the core type defination
type Fpm struct {
	// the start time of the instance
	starttime time.Time

	// the version of the core
	v string

	// the build time
	buildAt string

	// the routers, include the api, health, something else
	routers *mux.Router

	// the message queue, for pub and sub
	mq map[string][]MessageHandler

	// the lifecycle hooks for
	hooks map[string][]*Hook

	// the biz filters for
	filters map[string][]*Filter

	// the biz modules
	modules map[string]*BizModule

	// the logger
	Logger log.Logger

	// the cache instance
	cacher cache.Cache
	//appInfo
	appInfo *AppInfo
}

//HookHandler the hook handler
type HookHandler func(*Fpm)

//FilterHandler the hook handler
type FilterHandler func(app *Fpm, biz string, args *BizParam) (bool, error)

//MessageHandler message handler
type MessageHandler func(topic string, data interface{})

//AppInfo the basic config of the app
//"mode": "release",
// "domain": "",
// "version": "",
// "addr": ":9090",
// "name": "fpm-server",
type AppInfo struct {
	Mode    string
	Domain  string
	Version string
	Addr    string
	Name    string
}

//Hook the hook handler
type Hook struct {
	f HookHandler
	p int
}

//Filter the filter handler
type Filter struct {
	f FilterHandler
	p int
}

//NewHook create a new hook
func NewHook(f HookHandler, p int) *Hook {
	return &Hook{
		f: f,
		p: p,
	}
}

//NewFilter create a new filter
func NewFilter(f FilterHandler, p int) *Filter {
	return &Filter{
		f: f,
		p: p,
	}
}

//Handler the bizHandler
type Handler func(*ctx.Ctx, *Fpm)

//New 使用默认配置的构造函数
func New() *Fpm {
	return NewWithConfig("")
}

//NewWithConfig 初始化函数
//路由加载
//插件加载
//加载中间件
//执行init钩子函数
// BEFORE_INIT -> AFTER_INIT -> BEFORE_START -> BEFORE_SHUTDOWN(not sure) -> AFTER_SHUTDOWN(not sure)
func NewWithConfig(configFile string) *Fpm {
	if defaultInstance != nil {
		return defaultInstance
	}
	//加载配置文件
	if err := config.Init(configFile); err != nil {
		fmt.Println("Init config file error: ", configFile)
		panic(err)
	}

	fpm := &Fpm{}
	fpm.Logger = log.GetLogger()
	fpm.v = version.Version
	fpm.buildAt = version.BuildAt

	fpm.mq = make(map[string][]MessageHandler)
	fpm.routers = mux.NewRouter()
	fpm.hooks = make(map[string][]*Hook, 0)
	fpm.filters = make(map[string][]*Filter, 0)
	fpm.modules = make(map[string]*BizModule, 0)
	fpm.appInfo = &AppInfo{}

	if err := viper.Unmarshal(&(fpm.appInfo)); err != nil {
		panic(err)
	}

	fpm.loadPlugin()
	defaultInstance = fpm
	return fpm
}

//Default 获取默认的实例，通常可以避免不断传递 fpm 实例的引用
func Default() *Fpm {
	return defaultInstance
}

//Init run the init
func (fpm *Fpm) Init() {
	fpm.runHook("BEFORE_INIT")

	fpm.BindHandler("/health", func(c *ctx.Ctx, _ *Fpm) {
		c.JSON(map[string]interface{}{"status": "UP", "startAt": fpm.starttime, "version": fpm.v, "buildAt": fpm.buildAt})
	}).Methods("GET")

	fpm.BindHandler("/ping", func(c *ctx.Ctx, _ *Fpm) {
		c.JSON(map[string]interface{}{"status": "UP", "timestamp": time.Now().Unix()})
	}).Methods("GET")

	fpm.Use(middleware.Recover)
	fpm.BindHandler("/api", api).Methods("POST")
	fpm.BindHandler("/webhook/{upstream}/{event}/{method}", webhook).Methods("POST")
	fpm.routers.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	fpm.runHook("AFTER_INIT")
}

//GetAppInfo get the basic info of the app from the config
func (fpm *Fpm) GetAppInfo() *AppInfo {
	return fpm.appInfo
}

func api(c *ctx.Ctx, fpm *Fpm) {
	var data APIReq
	var rsp APIRsp
	rsp.Timestamp = time.Now().Unix()
	if err := c.ParseBody(&data); err != nil {
		rsp.Message = err.Error()
		rsp.Errno = -1
		rsp.Error = err
		c.Fail(rsp)
		return
	}
	method := data.Method

	result, err := fpm.Execute(method, data.Param)
	if err != nil {
		rsp.Message = err.Error()
		rsp.Errno = -1
		rsp.Error = err
		c.Fail(rsp)
		return
	}
	rsp.Errno = 0

	rsp.Data = result
	c.JSON(rsp)
}

func webhook(c *ctx.Ctx, fpm *Fpm) {

	upstream := c.Param("upstream")
	event := c.Param("event")
	method := c.Param("method")

	body := make(map[string]interface{})
	if err := c.ParseBody(&body); err != nil {
		c.Fail(err)
		return
	}

	go func() {
		body["url_data"] = method
		fpm.Publish(fmt.Sprintf("#webhook/%s/%s", upstream, event), body)
	}()

	c.JSON(map[string]interface{}{
		"errno": 0,
	})
}

//Publish publish a message
func (fpm *Fpm) Publish(topic string, data interface{}) {
	handlers, ok := fpm.mq[topic]
	if !ok {
		return
	}
	go func() {
		for _, handler := range handlers {
			handler(topic, data)
		}
	}()
}

//Subscribe subscribe a topic
func (fpm *Fpm) Subscribe(topic string, f MessageHandler) {
	handlers, ok := fpm.mq[topic]
	if !ok {
		handlers = make([]MessageHandler, 0)
	}
	handlers = append(handlers, f)
	fpm.mq[topic] = handlers
}

//SetCacher set the instance of cache
func (fpm *Fpm) SetCacher(c cache.Cache) {
	fpm.cacher = c
}

//GetCacher get the instance of cache
func (fpm *Fpm) GetCacher() (cache.Cache, bool) {
	if fpm.cacher == nil {
		return nil, false
	}
	return fpm.cacher, true
}

//Get get some key/val from the context
func Get(key string, dftVal interface{}) interface{} {
	val, ok := tempMapData[key]
	if !ok {
		return dftVal
	}
	return val
}

//Set set a key/val item into the context
func Set(key string, value interface{}) {
	tempMapData[key] = value
}

//loadPlugin load the plugins
func (fpm *Fpm) loadPlugin() {
	if len(registerEvents) < 1 {
		return
	}
	for _, event := range registerEvents {
		event(fpm)
	}
}

//HasConfig return true if config in the configfile
func (fpm *Fpm) HasConfig(key string) bool {
	return viper.IsSet(key)
}

//GetConfig get the config from the configfile
func (fpm *Fpm) GetConfig(key string) interface{} {
	return viper.Get(key)
}

//FetchConfig fetch config to the c
func (fpm *Fpm) FetchConfig(key string, c interface{}) error {
	if !viper.IsSet(key) {
		return errors.New("config: " + key + " not defined")
	}

	if err := utils.Interface2Struct(viper.Get(key), &c); err != nil {
		return err
	}
	return nil
}

//runHook 执行钩子函数
func (fpm *Fpm) runHook(hookName string) {
	hooks, exists := fpm.hooks[hookName]
	if !exists || len(hooks) < 1 {
		//No hooks
		return
	}
	for _, hook := range hooks {
		//TODO: sort by the priority desc
		hook.f(fpm)
	}
}

//runFilter 执行过滤器函数
func (fpm *Fpm) runFilter(filterName string, biz string, args *BizParam) (bool, error) {
	filters, exists := fpm.filters[filterName]
	if !exists || len(filters) < 1 {
		//No filters
		return true, nil
	}
	for _, filter := range filters {
		if ok, err := filter.f(fpm, biz, args); !ok {
			return false, err
		}
	}
	return true, nil
}

//AddFilter 增加一个钩子函数
func (fpm *Fpm) AddFilter(filterName, event string, handler FilterHandler, priority int) {
	filters, exists := fpm.filters["_"+filterName+"_"+event]
	if !exists {
		filters = make([]*Filter, 0)
	}
	filters = append(filters, NewFilter(handler, priority))
	fpm.filters["_"+filterName+"_"+event] = filters
}

//AddHook 增加一个钩子函数
func (fpm *Fpm) AddHook(hookName string, handler HookHandler, priority int) {
	hooks, exists := fpm.hooks[hookName]
	if !exists {
		hooks = make([]*Hook, 0)
	}
	hooks = append(hooks, NewHook(handler, priority))
	fpm.hooks[hookName] = hooks
}

//Execute 执行具体的业务函数
func (fpm *Fpm) Execute(biz string, args *BizParam) (interface{}, error) {
	if ok, err := fpm.runFilter("_"+biz+"_before", biz, args); !ok {
		return nil, err
	}
	bizPath := strings.Split(biz, ".")
	moduleName := bizPath[0]
	module, exists := fpm.modules[moduleName]
	if !exists {
		return nil, errNoMethod
	}
	bizName := strings.Join(bizPath[1:], ".")
	handler, exists := (*module)[bizName]
	if !exists {
		return nil, errNoMethod
	}
	result, err := handler(args)
	if err != nil {
		return nil, err
	}
	(*args)["__result__"] = result
	if ok, err := fpm.runFilter("_"+biz+"_after", biz, args); !ok {
		log.Errorf("run _%s_after error: %v", biz, err)
	}
	return result, nil
}

//Use add some middleware
func (fpm *Fpm) Use(mws ...func(next http.Handler) http.Handler) {
	for _, mw := range mws {
		fpm.routers.Use(mw)
	}
}

//AddBizModule 添加业务函数组
func (fpm *Fpm) AddBizModule(name string, module *BizModule) {
	fpm.modules[name] = module
}

//BindHandler 绑定接口路由
func (fpm *Fpm) BindHandler(url string, handler Handler) *mux.Route {
	f := func(w http.ResponseWriter, r *http.Request) {
		handler(ctx.WrapCtx(w, r), fpm)
	}
	return fpm.routers.HandleFunc(url, f)

}

//Run 启动程序
func (fpm *Fpm) Run() {
	fpm.runHook("BEFORE_START")
	fpm.starttime = time.Now()
	addr := fpm.appInfo.Addr
	if addr == "" {
		addr = ":9090"
	} else if !strings.HasPrefix(addr, ":") {
		addr = ":" + addr
	}
	srv := &http.Server{
		Handler: fpm.routers,
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  time.Second * 30,
	}

	wait := time.Second * 30

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	fpm.runHook("AFTER_START")
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Info("shutting down")
	os.Exit(0)

}
