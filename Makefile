network=solid
database_connect=postgres://lowcoder:@postgres:5432/solid?sslmode=disable
# command template for docker info
is_network=$(shell docker network ls -f 'name=solid' -q)
docker_container_list = $(shell docker ps -aq)
docker_image_list = $(shell docker images -q)
migration_version = $(shell make migrate-version 2>&1 1>/dev/null | cut -d ' ' -f1)
# command template for init extensions database
install_pgcrypto_extension = $(shell psql -U postgres -d solid -h localhost -c "CREATE EXTENSION pgcrypto CASCADE")

init-run: build run

run:
	docker run -p 5432:5432 -d --rm --name postgres --network $(network) postgres
	# TODO: in bad case golang service will run before postgres service
	docker run -p 8080:8080 -d -v $(PWD):/go/src/github.com/ilya-korotya/solid --rm --name solid --network $(network) solid

build:
	if [ -z $(is_network) ] ; then docker network create -d bridge $(network) ; fi
	docker build -t postgres ./database/postgres -f ./database/postgres/Dockerfile
	docker build -t solid . -f Dockerfile

# for testing full circle creating containers
clear:
	# with '|| true' we ignore error
	docker stop $(docker_container_list) || true
	docker rm $(docker_container_list) || true
	docker rmi $(docker_image_list) || true
	if [ ! -z $(is_network) ] ;	then docker network rm $(network) ; fi

rebuild:
	docker stop solid || true
	docker run -p 8080:8080 -d -v $(PWD):/go/src/github.com/ilya-korotya/solid --rm --name solid --network $(network) solid

# work with migration for postgres
migration: 
	# install extensions for custom database. FIX THIS SHIT PLZ  (╯°□°）╯︵ ┻━┻
	# it would be nice if they were installed at the start of the database
	# we will can init_database make as 'psql -U ... -d ... --file /path/to/init-db'
	docker exec -it -d postgres bash $(install_pgcrypto_extension) || true
	docker run -v $(PWD)/migrations/postgres:/migrations --network $(network) migrate/migrate -path=/migrations/ -database $(database_connect) up

migrate-create: 
	docker run --user=`id -u` -v "$(PWD)/migrations/postgres":/migrations --network $(network) migrate/migrate create -dir=/migrations/ -ext=.sql $(name)

migrate-down:
	docker run  -v "$(PWD)/migrations/postgres":/migrations --network $(network) migrate/migrate -path=/migrations/ -database $(database_connect) down $(count)

migrate-version: 
	docker run -v "$(PWD)/migrations/postgres":/migrations --network $(network) migrate/migrate -path=/migrations/ -database $(database_connect) version

migrate-fix:
	docker run  -v "$(PWD)/migrations/postgres":/migrations --network $(network) migrate/migrate -path=/migrations/ -database $(database_connect) force $(migration_version)

# TODO: add command(s) for run unit tests and create coverage
