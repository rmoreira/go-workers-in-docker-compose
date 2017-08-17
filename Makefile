default: build up

build:
	govendor init
	govendor add +external
	docker build -t goworkers:build -f Dockerfile.build .
	docker run --rm -v $(PWD):/tmp -t goworkers:build
	docker build -t goworkers:worker .

getdeps:
	go get -u github.com/kardianos/govendor

up:
	docker-compose up

clean:
	docker-compose rm -fs
