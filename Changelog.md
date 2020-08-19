## Changelog

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