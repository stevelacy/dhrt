VERSION=$(git describe --always --long)
NAME=dhrt

#GOOS=linux GOARCH=386 CGO_ENABLED=0
go build -i -v -ldflags="-X main.version=${VERSION}" -o ./build/${NAME} ./cmd/main.go
