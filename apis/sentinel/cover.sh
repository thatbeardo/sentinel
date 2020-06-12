go test ./... -v -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
echo .
echo Refresh your browser page for coverage.html