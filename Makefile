all:
	./build.sh
test:
	go test ./pkg/ -v
