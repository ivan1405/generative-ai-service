package awsbedrock

import (
	"gen-ai-service/app/service"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/bedrockruntime"
)

type AWSBedrockHandler struct {
	Client *bedrockruntime.BedrockRuntime
}

func (c *AWSBedrockHandler) Init(apiKey string) {
	mySession := session.Must(session.NewSession())
	c.Client = bedrockruntime.New(mySession, aws.NewConfig().WithRegion("us-west-2"))
}

func (c *AWSBedrockHandler) Type() string {
	return service.AmazonBedrockType
}

func (c *AWSBedrockHandler) ChatCompletion(message string) (string, error) {
	r, err := c.Client.InvokeModel(&bedrockruntime.InvokeModelInput{
		Body: []byte(message),
	})
	if err != nil {
		return "", err
	}
	return string(r.Body), nil
}
