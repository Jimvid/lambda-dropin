package main

import (
	// "github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssecretsmanager"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type GolangAwsStackProps struct {
	awscdk.StackProps
}

func NewLambdaDropin(scope constructs.Construct, id string, props *GolangAwsStackProps) awscdk.Stack {
	var sprops awscdk.StackProps

	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)

	// Secret from secretsmanager
	secret := awssecretsmanager.Secret_FromSecretNameV2(stack, jsii.String("jwtTokenSecret"), jsii.String("lambdadropin/jwt-secret"))

	// DynamoDB table
	table := awsdynamodb.NewTable(stack, jsii.String("userTable"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("username"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		TableName: jsii.String("userTable"),
	})

	// Lambda functions
	userFunction := awslambda.NewFunction(stack, jsii.String("userFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("../cmd/user/function.zip"), nil),
		Handler: jsii.String("main"),
		Environment: &map[string]*string{
			"SECRET_ARN": secret.SecretArn(),
		},
	})

	// API Gateway
	api := awsapigateway.NewRestApi(stack, jsii.String("APIGateway"), &awsapigateway.RestApiProps{
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowHeaders: jsii.Strings("Content-Type", "Authorization"),
			AllowMethods: jsii.Strings("GET", "POST", "PUT", "DELETE"),
			AllowOrigins: jsii.Strings("*"),
		},
		CloudWatchRole: jsii.Bool(false),
		DeployOptions: &awsapigateway.StageOptions{
			LoggingLevel: awsapigateway.MethodLoggingLevel_INFO,
		},
	})

	// Grant function access
	table.GrantReadWriteData(userFunction)
	secret.GrantRead(userFunction, nil)

	// Integration
	integration := awsapigateway.NewLambdaIntegration(userFunction, nil)

	// Register resource
	registerResource := api.Root().AddResource(jsii.String("register"), nil)
	registerResource.AddMethod(jsii.String("POST"), integration, nil)

	// Login resource
	loginResource := api.Root().AddResource(jsii.String("login"), nil)
	loginResource.AddMethod(jsii.String("POST"), integration, nil)

	// protected
	protectedResource := api.Root().AddResource(jsii.String("protected"), nil)
	protectedResource.AddMethod(jsii.String("GET"), integration, nil)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewLambdaDropin(app, "LambdaDropin", &GolangAwsStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return nil
}
