build:
	@cd cmd/user && GOOS=linux GOARCH=amd64 go build -o bootstrap
	@cd cmd/user && zip function.zip bootstrap

diff:
	@cd infrastructure && cdk diff

deploy:
	@cd infrastructure && cdk deploy
