build-to-docker: main.go
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o khaos-monkey .
	docker build -t khaos-monkey .