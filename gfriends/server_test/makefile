packagename=server_test


.PHONY: build
build: clean
	go build -o ${packagename} main.go

.PHONY: run
run:
	go run -race main.go

.PHONY: clean
clean:
	rm -f ${packagename}
