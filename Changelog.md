## Changelog

#### Otc 15

Feature:
- add `grpc` for biz api methods.
  - should add `fpm-go-plugin-grpc`

#### Sep 17

Feature:
- add `/metrics` for prometheus.
- add `/debug/pprof/` for debug

#### Sep 12

Feature:
- add `Basic Auth` for `/biz` to protect the api

#### Sep 5

Feature:
- add `GetToken`, `GetRemoteIP`, `GetRequestID` functions for `ctx.Ctx`
- change the Register Argument, it should contains `Name` & `Deps` 
#### Sep 3 v0.2.0

Feature:
- support run with no config file. use all default value.
- add jwt api, `/oauth/token`
- add `/biz/{module}/{method}` url for execute biz directly


#### August 31 v0.1.14

Feature:
- change the db interface

#### August 25 v0.1.9~v0.1.12

Feature:
- add cache interface
- offer a redis plugin `fpm-go-plugin-redis` todo `fpm-go-plugin-leveldb`
- defined db interface

#### August 22 v0.1.8

Feature:
- add default config for logger & addr, the log & addr node for config is not nessery
    - the default log output is `STDOUT`, the default addr is `:9090`
- add `GetHeader(string) string`  for fpm.Ctx
- change the `webhook` url to `/webhook/:upstream/:event/:data`
- add `ping` api

#### August 20 v0.1.5

Feature:
- add webhook router
- server static file folder

#### August 19 v0.1.3

Feature:
- add filter for biz execute
- before filter should return true, the biz executor should not run if not.
- after filter can run anyway, log error if fail, but the filterchain will crash, they will not run.

- the `__result__` will insert the param in the after filter.

#### August 17

Feature:
- use fpm.Logger to use the ref of the log

#### August 13

Feature:
- config reader with viper.
- add startup time and version for the core instance.
- panicHolder middleware for the core.