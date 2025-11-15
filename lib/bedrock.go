package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"log"
	"os"
)

const defaultRegion = "us-east-1"

var brc *bedrockruntime.Client

func init() {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = defaultRegion
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	brc = bedrockruntime.NewFromConfig(cfg)
}

const claudePromptFormat = "\n\nHuman: %s\n\nAssistant:"

func SendBedrock(ctx context.Context, msg string) (string, error) {
	msg = fmt.Sprintf(claudePromptFormat, msg)

	payload := Request{Prompt: msg, MaxTokensToSample: 2048}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	output, err := brc.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		Body:        payloadBytes,
		ModelId:     aws.String("anthropic.claude-v2:1"),
		ContentType: aws.String("application/json"),
	})

	if err != nil {
		return "", err
	}

	var resp Response

	err = json.Unmarshal(output.Body, &resp)

	if err != nil {
		return "", err
	}

	return resp.Completion, nil
}

//request/response model

type Request struct {
	Prompt            string   `json:"prompt"`
	MaxTokensToSample int      `json:"max_tokens_to_sample"`
	Temperature       float64  `json:"temperature,omitempty"`
	TopP              float64  `json:"top_p,omitempty"`
	TopK              int      `json:"top_k,omitempty"`
	StopSequences     []string `json:"stop_sequences,omitempty"`
}

type Response struct {
	Completion string `json:"completion"`
}
