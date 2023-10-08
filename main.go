package main

import (
	awsbedrock "gen-ai-service/app/aws-bedrock"
	chatgpt "gen-ai-service/app/chat-gpt"
	"gen-ai-service/app/service"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"
)

func main() {
	loadEnvVars()

	// Initialize chat-gpt
	chatGptHandler := &chatgpt.ChatGptHandler{}
	chatGptHandler.Init(os.Getenv("CHAT_GPT_KEY"))

	// Initialize aws-bedrock
	awsBedrockHandler := &awsbedrock.AWSBedrockHandler{}
	awsBedrockHandler.Init(os.Getenv("AWS_BEDROCK_KEY"))

	genAIService := service.InitRouter()
	genAIService.ConfigureHandlers(
		chatGptHandler,
		awsBedrockHandler,
	)
	genAIService.Start()
}

func loadEnvVars() {
	if err := godotenv.Load(); err != nil {
		slog.Error("No .env file found")
	}
}
