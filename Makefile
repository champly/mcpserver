build:
	GOOS=linux GOARCH=amd64 go build -o mcpserver main.go
	docker build -t registry.cn-hangzhou.aliyuncs.com/champly/mcpserver:v1 .
	docker push registry.cn-hangzhou.aliyuncs.com/champly/mcpserver:v1