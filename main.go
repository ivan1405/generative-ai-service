package main

import (
	awsbedrock "gen-ai-service/app/aws-bedrock"
	chatgpt "gen-ai-service/app/chat-gpt"
	elevenlabs "gen-ai-service/app/eleven-labs"
	"gen-ai-service/app/service"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"
)

func main() {
	loadEnvVars()

	// Initialize chat-gpt
	chatGptHandler := chatgpt.NewChatGptHandler(os.Getenv("CHAT_GPT_KEY"))

	// Initialize aws-bedrock
	awsBedrockHandler := awsbedrock.NewAWSBedrockHandler()

	// Initialize eleven-labs
	elevenLabsHandler := elevenlabs.NewElevenLabsHandler(os.Getenv("ELEVEN_LABS_API_KEY"))

	genAIService := service.InitRouter()
	genAIService.ConfigureHandlers(
		chatGptHandler,
		awsBedrockHandler,
		elevenLabsHandler,
	)
	genAIService.Start()
}

func loadEnvVars() {
	if err := godotenv.Load(); err != nil {
		slog.Error("No .env file found")
	}
}
