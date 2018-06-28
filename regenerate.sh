#!/bin/bash


rm -rf modelsv2 clientv2

docker run --rm -v $(pwd):/go/src/github.com/fnproject/fn_go -v $(pwd)/../fn/docs/swagger_v2.yml:/go/src/github.com/fnproject/fn/swagger.yml -w /go/src/github.com/fnproject/fn_go quay.io/goswagger/swagger:0.14.0 generate client -f /go/src/github.com/fnproject/fn/swagger.yml -A fn -m modelsv2 -c clientv2

docker run --rm -v $(pwd):/go/src/github.com/fnproject/fn_go -v $(pwd)/../fn/docs/swagger.yml:/go/src/github.com/fnproject/fn/swagger.yml -w /go/src/github.com/fnproject/fn_go quay.io/goswagger/swagger:0.14.0 generate client -f /go/src/github.com/fnproject/fn/swagger.yml -A fn
