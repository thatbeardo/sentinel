check_install:
		go get -u github.com/swaggo/swag/cmd/swag
		go get -u github.com/swaggo/gin-swagger
		go get -u github.com/swaggo/files

swagger: check_install
		swag init -g application.go