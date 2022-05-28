PROJECT=${p}

all:
	@echo 'Start Compose Demo:'
	@echo '  docker-compose up -d'
	@echo ''
	@echo 'Build Image:'
	@echo '  docker build . -t init-golang:0.0.0'

run:
	@echo "docker-compose up -d"

docker:
	@echo "docker build . -t init-golang:0.0.0"
	@echo "docker tag init-golang:0.0.0 docker-hub.kittymeow.cc:5000/init-golang:0.0.0"

linux:
	CGO_ENABLED=0 GOOS=linux GOATCH=amd64 bash build.sh ${PROJECT}

macos:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 bash build.sh ${PROJECT}

clbuild:
	rm -rf bin/*

climg:
	@docker images | grep '<none>'
	@docker images | grep '<none>' | awk '{print $$3}' | xargs -n 1 docker rmi
