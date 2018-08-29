all:
	./build.sh
test:
	go test ./pkg/ -v
clean:
	rm -rf build/*
setup:
	dep ensure
