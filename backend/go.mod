module github.com/halimath/fate-core-remote-table/backend

go 1.22

toolchain go1.22.1

replace github.com/halimath/httputils v0.0.0 => ../../../go/httputils

require (
	github.com/google/uuid v1.6.0
	github.com/halimath/expect v0.5.1
	github.com/halimath/httputils v0.0.0
	github.com/halimath/jose v0.0.0-20210820062418-4ca508234dee
	github.com/halimath/kvlog v0.11.1
	github.com/oapi-codegen/runtime v1.1.1
	github.com/sethvargo/go-envconfig v1.0.1
)

require github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
