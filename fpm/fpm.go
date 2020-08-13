//Package fpm the core fpm
package fpm

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/justinas/alice"
	"github.com/spf13/viper"

	"github.com/gorilla/mux"
	"github.com/team4yf/yf-fpm-server-go/config"
	"github.com/team4yf/yf-fpm-server-go/ctx"
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
	mq chan map[string]string

	// the lifecycle hooks for
	hooks map[string][]*Hook

	// the biz modules
	modules map[string]*BizModule

	// middlware chain
	mwChain alice.Chain
}

//HookHandler the hook handler
type HookHandler func(*Fpm)

//Hook the hook handler
type Hook struct {
	f HookHandler
	p int
}

//NewHook create a new hook
func NewHook(f HookHandler, p int) *Hook {
	return &Hook{
		f: f,
		p: p,
	}
}

//Handler the bizHandler
type Handler func(*ctx.Ctx, *Fpm)

//New 初始化函数
//路由加载
//插件加载
//加载中间件
//执行init钩子函数
// BEFORE_INIT -> AFTER_INIT -> BEFORE_START -> BEFORE_SHUTDOWN(not sure) -> AFTER_SHUTDOWN(not sure)
func New() *Fpm {
	if defaultInstance != nil {
		return defaultInstance
	}
	//加载配置文件
	config.Init("")

	fpm := &Fpm{}
	fpm.v = version.Version
	fpm.buildAt = version.BuildAt

	fpm.mq = make(chan map[string]string, 1000)
	fpm.routers = mux.NewRouter()
	fpm.hooks = make(map[string][]*Hook, 0)
	fpm.modules = make(map[string]*BizModule, 0)

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
		c.JSON(map[string]interface{}{"Status": "UP", "StartAt": fpm.starttime, "version": fpm.v, "buildAt": fpm.buildAt})
	}).Methods("GET")

	fpm.Use(RecoverMiddleware)
	fpm.BindHandler("/api", api).Methods("POST")
	fpm.runHook("AFTER_INIT")
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
	return viper.InConfig(key)
}

//GetConfig get the config from the configfile
func (fpm *Fpm) GetConfig(key string) interface{} {
	return viper.Get(key)
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
	return handler(args)
}

//Use add some middleware
func (fpm *Fpm) Use(mw ...alice.Constructor) {
	fpm.mwChain = alice.New(mw...)
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
	return fpm.routers.Handle(url, fpm.mwChain.ThenFunc(f))
}

//Run 启动程序
func (fpm *Fpm) Run(addr string) {
	fpm.runHook("BEFORE_START")
	fpm.starttime = time.Now()
	log.Fatal(http.ListenAndServe(addr, fpm.routers))

}
