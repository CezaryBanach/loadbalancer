all: backend loadbalancer

backend: backend/main.go
	@echo "building backend ..."
	# cp ./backend/response_template.txt ./bin/response_template.txt
	go build -o ./bin/backend ./backend/main.go

loadbalancer: loadbalancer/main.go
	@echo "building loadbalancer ..."
	go build -o ./bin/loadbalancer ./loadbalancer/main.go

.PHONY: backend loadbalancer clean


clean:
	@echo "Cleaning up binaries..."
	rm -f  ./bin/*