//Package core the core fpm
package core

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//Fpm the core type defination
type Fpm struct {
	// the routers, include the api, health, something else
	routers *mux.Router

	// the message queue, for pub and sub
	mq chan map[string]string

	// the lifecycle hooks for
	hooks map[string][]*Hook
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
type Handler func(*Ctx)

//New 初始化函数
//路由加载
//插件加载
//加载中间件
//执行init钩子函数
// BEFORE_INIT -> AFTER_INIT -> BEFORE_START -> BEFORE_SHUTDOWN(not sure) -> AFTER_SHUTDOWN(not sure)
func (fpm *Fpm) New() {
	fpm.mq = make(chan map[string]string, 1000)
	fpm.routers = mux.NewRouter()
	fpm.hooks = make(map[string][]*Hook, 0)

	fpm.loadPlugin()

}

//Init run the init
func (fpm *Fpm) Init() {
	fpm.runHook("BEFORE_INIT")

	fpm.runHook("AFTER_INIT")
}

//Get get some key/val from the context
func (fpm *Fpm) Get(key string) {

}

//Set set a key/val item into the context
func (fpm *Fpm) Set(key string, value interface{}) {

}

//loadPlugin load the plugins
func (fpm *Fpm) loadPlugin() {

}

//GetConfig get the config from the configfile
func (fpm *Fpm) GetConfig(key string) {

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

//BindHandler 绑定接口路由
func (fpm *Fpm) BindHandler(url string, handler Handler) *mux.Route {
	return fpm.routers.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		handler(WrapCtx(fpm, w, r))
	})
}

//Run 启动程序
func (fpm *Fpm) Run(addr string) {
	fpm.runHook("BEFORE_START")

	log.Fatal(http.ListenAndServe(addr, fpm.routers))

}
