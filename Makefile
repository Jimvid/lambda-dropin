build:
	@cd cmd/user && GOOS=linux GOARCH=amd64 go build -o bootstrap
	@cd cmd/user && zip function.zip bootstrap

diff:
	@cd infrastructure && cdk diff
	
synth:
	@cd infrastructure && cdk synth

deploy:
	@cd infrastructure && cdk deploy
