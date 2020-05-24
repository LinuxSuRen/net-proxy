build:
	go build -o bin/net-proxy
	rm -rf net-proxy && ln -s bin/net-proxy net-proxy

image:
	docker build -t surenpi/net-proxy .

run:
	./net-proxy

run-image:
	docker run -p 8089:8089 surenpi/net-proxy
