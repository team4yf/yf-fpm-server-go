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

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"

	"github.com/gorilla/mux"
	"github.com/team4yf/fpm-go-pkg/log"
	"github.com/team4yf/fpm-go-pkg/utils"
	"github.com/team4yf/yf-fpm-server-go/config"
	"github.com/team4yf/yf-fpm-server-go/ctx"
	"github.com/team4yf/yf-fpm-server-go/errno"
	"github.com/team4yf/yf-fpm-server-go/middleware"
	"github.com/team4yf/yf-fpm-server-go/pkg/cache"
	"github.com/team4yf/yf-fpm-server-go/pkg/db"
	"github.com/team4yf/yf-fpm-server-go/version"
)

var (
	registerPlugins map[string]*Plugin

	errNoMethod = errors.New("NO_METHOD_DEFINED")

	defaultInstance *Fpm

	tempMapData map[string]interface{}
)

func init() {
	tempMapData = make(map[string]interface{})
	registerPlugins = make(map[string]*Plugin)
}

//Register register a nonamed plugin
func Register(handler func(*Fpm)) {
	RegisterByPlugin(&Plugin{
		Name:    utils.GenShortID(),
		Handler: handler,
	})
}

//RegisterByPlugin register some plugin
func RegisterByPlugin(event *Plugin) {
	registerPlugins[event.Name] = event
}

//Fpm the core type definition
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

	// the distribute locker
	distributeLocker cache.SyncLocker

	//appInfo
	appInfo *AppInfo

	//database interface
	database map[string]func() db.Database

	//health check data
	HealthCheckData *HealthCheckRsp
}

//HookHandler the hook handler
type HookHandler func(*Fpm)

//Plugin the plugin
type Plugin struct {
	Handler   func(*Fpm)
	Name      string
	Deps      []string
	Installed bool
	V         string
}

//FilterHandler the hook handler
type FilterHandler func(app *Fpm, biz string, args *BizParam) (bool, error)

//MessageHandler message handler
type MessageHandler func(topic string, data interface{})

//AppInfo the basic config of the app
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

//HealthCheckRsp health check rsp
type HealthCheckRsp struct {
	Status    string     `json:"status"`
	Hostname  string     `json:"hostname"`
	StartAt   *time.Time `json:"startAt"`
	Version   string     `json:"version"`
	BuildAt   string     `json:"buildAt"`
	GitHashID string     `json:"gitHashId"`
	GitBranch string     `json:"gitBranch"`
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

//New create with default configFile
func New() *Fpm {
	return NewWithConfig("")
}

//NewWithConfig create with specified configFile
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
	fpm.database = make(map[string]func() db.Database, 0)
	fpm.appInfo = &AppInfo{
		Name:    "yf-fpm-server-go",
		Mode:    "debug",
		Domain:  "localhost",
		Version: version.Version,
		Addr:    ":9090",
	}
	fpm.HealthCheckData = &HealthCheckRsp{
		Status:    "UP",
		StartAt:   &fpm.starttime,
		Hostname:  utils.GetHostname(),
		Version:   fpm.v,
		BuildAt:   fpm.buildAt,
		GitBranch: "master",
		GitHashID: "-",
	}

	if err := viper.Unmarshal(&(fpm.appInfo)); err != nil {
		panic(err)
	}

	fpm.loadPlugin()
	defaultInstance = fpm

	//pass the jwt key
	utils.InitJWTUtil(fpm.GetConfig("jwt.secret").(string))

	return fpm
}

//Default get pointer of default instance
func Default() *Fpm {
	return defaultInstance
}

//Init init the server instance
// - load routers
// - load plugins
// - load middlewares
// - execute init hooks
// - run hooks BEFORE_INIT -> AFTER_INIT -> BEFORE_START -> BEFORE_SHUTDOWN(not sure) -> AFTER_SHUTDOWN(not sure)
func (fpm *Fpm) Init() {
	fpm.runHook("BEFORE_INIT")

	fpm.BindHandler("/health", func(c *ctx.Ctx, _ *Fpm) {
		c.JSON(fpm.HealthCheckData)
	}).Methods("GET")

	fpm.BindHandler("/ping", func(c *ctx.Ctx, _ *Fpm) {
		c.JSON(map[string]interface{}{"status": "UP", "timestamp": time.Now().Unix()})
	}).Methods("GET")

	fpm.Use(middleware.Recover)

	if fpm.HasConfig("auth") {
		basicAuthConfig := middleware.BasicAuthConfig{
			Enable: false,
		}
		if err := fpm.FetchConfig("auth", &basicAuthConfig); err != nil {
			panic(err)
		}
		fpm.Use(middleware.BasicAuth(&basicAuthConfig))
	}

	if fpm.HasConfig("serverAuth") {
		serverAuthConfig := middleware.ServerAuthConfig{
			Enable: false,
		}
		if err := fpm.FetchConfig("serverAuth", &serverAuthConfig); err != nil {
			panic(err)
		}
		fpm.Use(middleware.ServerAuth(&serverAuthConfig))
	}

	fpm.BindHandler("/api", api).Methods("POST")

	fpm.BindHandler("/biz/{module}/{method}", biz).Methods("POST", "GET")
	fpm.BindHandler("/webhook/{upstream}/{event}/{method}", webhook).Methods("POST")
	fpm.SetStatic("/static/", "./static")
	initOauth2(fpm)

	registerPrometheus(fpm)
	attachProfiler(fpm.routers)

	fpm.runHook("AFTER_INIT")
}

//SetStatic set static
// prefix should starts and ends with slash like: /static/
// dir can be a relative or absolte filepath
func (fpm *Fpm) SetStatic(prefix, dir string) {
	fpm.routers.PathPrefix(prefix).Handler(http.StripPrefix(prefix, http.FileServer(http.Dir(dir))))
}

//GetAppInfo get the basic info of the app from the config
func (fpm *Fpm) GetAppInfo() *AppInfo {
	return fpm.appInfo
}
func initOauth2(fpm *Fpm) {

	//grant_type=client_credentials&client_id=000000&client_secret=999999&scope=read
	//curl "localhost:9090/oauth/token?grant_type=client_credentials&client_id=123123123&client_secret=123123123&scope=admin"
	//we get  {"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEyMzEyMzEyMyJ9.5KEyBxDH2NGVKoiA0J6IPB4QPlvZi9zPH9SSKTWF2h8","expires_in":7200,"scope":"admin","token_type":"Bearer"}

	fpm.BindHandler("/oauth/token", func(c *ctx.Ctx, fpm *Fpm) {
		querys := c.Querys()
		gt := querys["grant_type"]
		if gt != "client_credentials" {
			c.BizError(errno.OAuthOnlySupportClientErr)
			return
		}
		id := querys["client_id"]
		secret := querys["client_secret"]
		//TODO need check here
		if id != "123123123" || secret != "123123123" {
			c.BizError(errno.OAuthClientAuthErr)
			return
		}
		scope := querys["scope"]
		exp := 720000
		tokenStr, _ := utils.GenerateToken(&jwt.MapClaims{
			"id":  id,
			"exp": exp,
		})

		c.JSON(map[string]interface{}{
			"access_token": tokenStr,
			"expires_in":   exp,
			"scope":        scope,
			"token_type":   "Bearer",
		})
	}).Methods("GET")
}

func biz(c *ctx.Ctx, fpm *Fpm) {
	method := c.Param("method")
	module := c.Param("module")
	method = module + "." + method
	param := BizParam{}
	if c.GetRequest().Method == "POST" {
		c.ParseBody(&param)
	} else {
		//Get
		querys := c.Querys()
		utils.Interface2Struct(querys, &param)
	}

	data, err := fpm.Execute(method, &param)
	if err != nil {
		c.BizError(errno.Wrap(err))
		return
	}
	c.JSON(ResponseOK(data))
}

func api(c *ctx.Ctx, fpm *Fpm) {
	var data APIReq
	if err := c.ParseBody(&data); err != nil {
		c.BizError(errno.Wrap(err))
		return
	}
	method := data.Method
	if data.Raw != nil {
		p := BizParam{}
		// Raw 可能是string类型，也可能是 object 类型
		switch data.Raw.(type) {
		case string:
			if err := utils.StringToStruct(data.Raw.(string), &p); err != nil {
				c.BizError(errno.Wrap(err))
				return
			}
		default:
			if err := utils.Interface2Struct(data.Raw.(string), &p); err != nil {
				c.BizError(errno.Wrap(err))
				return
			}
		}
		data.Param = &p
	}

	result, err := fpm.Execute(method, data.Param)
	if err != nil {
		c.BizError(errno.Wrap(err))
		return
	}
	c.JSON(ResponseOK(result))
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

	c.JSON(ResponseOK(1))
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

//SetDistributeLocker set the instance of Distribute Locker
func (fpm *Fpm) SetDistributeLocker(l cache.SyncLocker) {
	fpm.distributeLocker = l
}

//GetDistributeLocker get the instance of Distribute Locker
func (fpm *Fpm) GetDistributeLocker() (cache.SyncLocker, bool) {
	if fpm.distributeLocker == nil {
		return nil, false
	}
	return fpm.distributeLocker, true
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

//GetDatabase get some key/val from the context
func (fpm *Fpm) GetDatabase(name string) (db.Database, bool) {
	provider, ok := fpm.database[name]
	if !ok {
		return nil, false
	}
	return provider(), true
}

//SetDatabase set a db interface
func (fpm *Fpm) SetDatabase(name string, provider func() db.Database) {
	fpm.database[name] = provider
}

//loadPlugin load the plugins
func (fpm *Fpm) loadPlugin() {
	//the plugin should contains dependence of the other
	//we should make them run with sequence

	for {
		done := true
		for name, event := range registerPlugins {
			done = true
			if event.Installed {
				continue
			}
			done = false
			if len(event.Deps) < 1 {
				// no dep, run now
				event.Installed = true
				event.Handler(fpm)
				continue
			}
			depDone := true
			for _, d := range event.Deps {
				p, ok := registerPlugins[d]
				if !ok {
					panic(fmt.Sprintf("plugin: %s required: %s, but now installed.", name, d))
				}
				if !p.Installed {
					// dep not installed yet
					depDone = false
				}
			}
			if depDone {
				// all dependency installed
				event.Installed = true
				event.Handler(fpm)
			}
		}
		if done {
			break
		}
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

//InstalledPlugins get all installed plugins
func (fpm *Fpm) InstalledPlugins() []string {
	names := make([]string, 0)
	for m := range registerPlugins {
		names = append(names, m)
	}
	return names
}

//IsInstalledPlugin check if the plugin installed
func (fpm *Fpm) IsInstalledPlugin(name string) bool {
	_, ok := registerPlugins[name]
	return ok
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

//runHook run hooks
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

//runFilter run filter functions
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

//AddFilter add a filter function
func (fpm *Fpm) AddFilter(filterName, event string, handler FilterHandler, priority int) {
	filters, exists := fpm.filters["_"+filterName+"_"+event]
	if !exists {
		filters = make([]*Filter, 0)
	}
	filters = append(filters, NewFilter(handler, priority))
	fpm.filters["_"+filterName+"_"+event] = filters
}

//AddHook add a hook function
func (fpm *Fpm) AddHook(hookName string, handler HookHandler, priority int) {
	hooks, exists := fpm.hooks[hookName]
	if !exists {
		hooks = make([]*Hook, 0)
	}
	hooks = append(hooks, NewHook(handler, priority))
	fpm.hooks[hookName] = hooks
}

//Execute execute biz function
func (fpm *Fpm) Execute(biz string, args *BizParam) (data interface{}, err error) {
	defer func() {
		if err != nil {
			incBizExecuteVec(biz, "fail")
		} else {
			incBizExecuteVec(biz, "success")
		}
	}()
	ok := false
	if ok, err = fpm.runFilter("_"+biz+"_before", biz, args); !ok {
		return
	}
	bizPath := strings.Split(biz, ".")
	moduleName := bizPath[0]
	module, exists := fpm.modules[moduleName]
	if !exists {
		err = errNoMethod
		return
	}
	bizName := strings.Join(bizPath[1:], ".")
	handler, exists := (*module)[bizName]
	if !exists {
		err = errNoMethod
		return
	}
	data, err = handler(args)
	if err != nil {
		return
	}
	(*args).__result__ = data
	if ok, err = fpm.runFilter("_"+biz+"_after", biz, args); !ok {
		log.Errorf("run _%s_after error: %v", biz, err)
	}
	if data == nil {
		data = 1
	}
	return
}

//Use add some middleware
func (fpm *Fpm) Use(mws ...func(next http.Handler) http.Handler) {
	for _, mw := range mws {
		fpm.routers.Use(mw)
	}
}

//AddBizModule add biz module
func (fpm *Fpm) AddBizModule(name string, module *BizModule) {
	fpm.modules[name] = module
}

//BindHandler bind router handler
func (fpm *Fpm) BindHandler(url string, handler Handler) *mux.Route {
	f := func(w http.ResponseWriter, r *http.Request) {
		handler(ctx.WrapCtx(w, r), fpm)
	}
	return fpm.routers.HandleFunc(url, f)

}

//Run startup server
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
