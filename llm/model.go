package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"google.golang.org/api/option"
)

const (
	ProviderOpenAI = "openai"
	ProviderGemini = "gemini"
	ProviderDoubao = "doubao"

	// Model constants
	geminiModel = "gemini-pro"
)

const (
	DefaultApiKey   = "NmUzZTQzOGMtYTM4MC00ZWQ1LWI1OTctZTAxY2I4MmJjNGRm"
	DefaultEndpoint = "ZXAtMjAyNTAxMTAyMDI1MDMtZmRrZ3E="
)

const llmPrompt = `Generate a concise and informative Git commit message based on the following code diff.

The commit message should follow these rules:
1. Follow the Conventional Commits format: <type>(<scope>): <description>
2. The body should be one paragraph
3. The body should explain WHAT and WHY (not HOW)
4. Each line should be less than 72 characters
5. There should be a line break between the title and the body

Here's the diff:`

func GenerateGeminiCommitMessage(diff, apiKey string) (string, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", fmt.Errorf("creating Gemini client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(geminiModel)
	prompt := fmt.Sprintf("%s\n%s", llmPrompt, diff)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("generating commit message: %w", err)
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		generatedMessage := resp.Candidates[0].Content.Parts[0].(genai.Text)
		return strings.TrimSpace(string(generatedMessage)), nil
	}

	return "", fmt.Errorf("no commit message generated by Gemini")
}

func generateOpenAICommitMessage(diff, apiKey string) (string, error) {
	// TODO: Implement OpenAI integration
	return "", fmt.Errorf("OpenAI integration not implemented yet")
}

func GenerateDoubaoCommitMessage(diff, apiKey string, endpointId string) (string, error) {
	client := arkruntime.NewClientWithApiKey(
		apiKey,
	)

	ctx := context.Background()

	prompt := fmt.Sprintf("%s\n%s", llmPrompt, diff)

	req := model.ChatCompletionRequest{
		Model: endpointId,
		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleSystem,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String("你是豆包，是由字节跳动开发的 AI 人工智能助手"),
				},
			},
			{
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String(prompt),
				},
			},
		},
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	return *resp.Choices[0].Message.Content.StringValue, nil

}
