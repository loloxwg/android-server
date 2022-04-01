SRCDIRS_EXCLUDE = proto log logs deploy ops ops_initctl
SRCDIRS_ALL = $(sort $(subst /,,$(dir $(wildcard */*.go))))
SRCDIRS = $(filter-out $(SRCDIRS_EXCLUDE), $(SRCDIRS_ALL))

PKGDIRS_EXCLUDE=$(GOROOT)/pkg
PKGDIRS_ALL = $(addsuffix /pkg, $(subst :, ,$(GOPATH)))
PKGDIRS = $(filter-out $(PKGDIRS_EXCLUDE), $(PKGDIRS_ALL))

EXEC_PREFIX = android_server
BUILD_VERSION = v0.1.0
BUILD_TIME = $(shell date "+%F %T")
BUILD_NAME = android_server_$(shell date "+%Y%m%d%H")
EXEC_NAME = $(EXEC_PREFIX)
REGISTRY = registry.sensetime.com

all: build_main 
	@for subdir in $(SRCDIRS);do \
		cd $$subdir; go install; cd ..; \
	done 

build_main:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags \
	" \
		-X 'main.BuildVersion=${BUILD_VERSION}'     \
		-X 'main.BuildTime=${BUILD_TIME}'     \
		-X 'main.BuildName=${BUILD_NAME}'     \
	" \
	-o $(EXEC_NAME) cmd/main.go

build_local:
	go build -ldflags \
	" \
		-X 'main.BuildVersion=${BUILD_VERSION}'     \
		-X 'main.BuildTime=${BUILD_TIME}'     \
		-X 'main.BuildName=${BUILD_NAME}'     \
	" \
	-o $(EXEC_NAME) cmd/main.go

show:
	@echo "==================src====================="
	@echo SRCDIRS_ALL: $(SRCDIRS_ALL)
	@echo SRCDIRS_EXCLUDE: $(SRCDIRS_EXCLUDE)
	@echo SRCDIRS: $(SRCDIRS)
	@echo "==================pkg====================="
	@echo PKGDIRS_EXCLUDE: $(PKGDIRS_EXCLUDE)
	@echo PKGDIRS_ALL: $(PKGDIRS_ALL)
	@echo PKGDIRS: $(PKGDIRS)
	@echo "================clean====================="
	@for subdir in $(PKGDIRS); do \
		cd $$subdir;\
		module_name=`echo $(CURDIR)|awk -F"/" '{print $$(NF)}'`;\
		result=`find . |grep $$module_name |head -n1|awk -F"." '{print $$2}'`; \
		if [ -n "$$result" ];then \
			echo clean_dirs:$$subdir$$result; \
		fi \
	done

image:
	docker build --network=host -t $(REGISTRY)/storage/$(EXEC_NAME):$(BUILD_VERSION) -f Dockerfile .

push: image
	docker push $(REGISTRY)/storage/$(EXEC_NAME):$(BUILD_VERSION)

test:
	go test ./...
	
lint:
	@golangci-lint run --deadline=5m

clean:
	rm ${EXEC_NAME}

