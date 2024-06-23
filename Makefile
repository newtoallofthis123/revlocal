build:
	@cd server && go build -o ../bin/server
	@cd client && go build -o ../bin/client
	@cd proxy && go build -o ../bin/proxy

server: build
	@./bin/server

client: build
	@./bin/client

proxy: build
	@./bin/proxy

clean:
	@rm -rf bin
	@mkdir bin

.PHONY: build server client proxy clean