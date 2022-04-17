GOOS=linux GOARCH=mipsle go build -o movie-sync-server-linux-mipsle
GOOS=linux GOARCH=amd64 go build -o movie-sync-server-linux-amd64
GOOS=linux GOARCH=arm64 go build -o movie-sync-server-linux-arm64
GOOS=linux GOARCH=mips go build -o movie-sync-server-linux-mips
GOOS=windows GOARCH=amd64 go build -o movie-sync-server-win-amd64.exe