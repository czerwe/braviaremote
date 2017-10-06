dependencys:
	go get github.com/Sirupsen/logrus
	go get github.com/gorilla/mux
	go get github.com/jessevdk/go-flags
	go get github.com/czerwe/gobravia
	
armbuild:
	env GOOS=linux CGO_ENABLED=0 GOARCH=arm go build braviactl.go

build:
	env GOOS=linux CGO_ENABLED=0 go build braviactl.go

buildimage:
	docker build -t abuild -f Dockerfile_build .

buildindocker:
	make buildimage
	docker run -v  $(shell pwd):/go/src/alexaservice --workdir /go/src/alexaservice abuild make build

docker:
	docker build -t braviaremote .