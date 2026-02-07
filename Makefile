dev:
	go run github.com/githubnemo/CompileDaemon -build="go build -o main ." -command="./main"

install-tools:
	go mod download
