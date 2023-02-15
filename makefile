# Copyright 2019 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.EXPORT_ALL_VARIABLES:

# settings
# TODO: make tls enable for net_dev private registry
REGISTRY?=ccr.ccs.tencentyun.com/edisonlai
TAG?=$(shell git describe --tags)

REPO_ROOT:=${CURDIR}
OUT_DIR=$(REPO_ROOT)/bin

# go
TARGETOS?=linux
TARGETARCH?=amd64
GOPROXY?=http://mirrors.aliyun.com/goproxy/,direct

# ldflags
VERSION_PKG=github.com/EdisonLai/ddns/pkg/utils
GIT_COMMIT=$(shell git rev-parse HEAD)
BUILD_DATE=$(shell date +%Y-%m-%dT%H:%M:%S%z)
ldflags="-s -w -X $(VERSION_PKG).Version=$(TAG) -X $(VERSION_PKG).GitVersion=${GIT_COMMIT} -X ${VERSION_PKG}.BuildTime=${BUILD_DATE}"

.PHONY: ddns
ddns:
	@echo + Building client binary
	rm -rf ${OUT_DIR} \
	GOSUMDB=off \
	go mod tidy
	GOARCH=${TARGETARCH} \
	GOOS=${TARGETOS} \
	CGO_ENABLED=0 \
	GO111MODULE=on \
	go build -v -o $(OUT_DIR)/ \
		    -ldflags $(ldflags) cmd/client/client.go
	go build -v -o $(OUT_DIR)/ \
		    -ldflags $(ldflags) cmd/server/server.go
	@echo + Built ddns binary to $(OUT_DIR)

.PHONY: image
image:
	docker build -t $(REGISTRY)/ddns:$(TAG) -f cmd/client/Dockerfile .
	docker build -t $(REGISTRY)/ddns-server:$(TAG) -f cmd/server/Dockerfile .
	@echo + Building image $(REGISTRY):$(TAG) successfully

.PHONY: push
push: push
	docker push $(REGISTRY)/ddns:$(TAG)
	docker push $(REGISTRY)/ddns-server:$(TAG)
