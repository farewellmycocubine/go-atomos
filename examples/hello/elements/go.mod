module "elements"

go 1.16

require (
	github.com/hwangtou/go-atomos v0.1.3
	github.com/hwangtou/go-atomos/examples/hello/api v1.0.0
)

replace (
	github.com/hwangtou/go-atomos v0.1.3 => ../../../
	github.com/hwangtou/go-atomos/examples/hello/api v1.0.0 => ../api
)