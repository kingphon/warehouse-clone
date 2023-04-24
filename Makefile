#!bin/bash

export DOMAIN_WAREHOUSE_ADMIN=localhost:3000
export DOMAIN_WAREHOUSE_APP=localhost:3001
export ENV=develop
export ZOOKEEPER_URI=127.0.0.1:2181
export ZOOKEEPER_PREFIX_EXTERNAL=/selly
export ZOOKEEPER_PREFIX_WAREHOUSE_COMMON=/selly_warehouse/common
export ZOOKEEPER_PREFIX_WAREHOUSE_ADMIN=/selly_warehouse/admin
export ZOOKEEPER_PREFIX_WAREHOUSE_APP=/selly_warehouse/app


# make update-submodules branch=develop
update-submodules:
	git submodule update --init --recursive && \
	git submodule foreach git checkout $(branch) && \
	git submodule foreach git pull origin $(branch)

run-admin:
	go run cmd/admin/main.go


run-app:
	go run cmd/app/main.go

swagger-admin:
	swag init -d ./ -g cmd/admin/main.go \
    --exclude ./pkg/app \
    -o ./docs/admin --pd

swagger-app:
	swag init -d ./ -g cmd/app/main.go \
	--exclude ./pkg/admin \
	-o ./docs/app --pd
# delete submodules folder in git cache
# git rm --cached submodules