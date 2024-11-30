package script_generator

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// MeditationRequest represents the parameters for generating a meditation
type MeditationRequest struct {
	DurationMinutes int    `json:"duration_minutes"`
	Focus           string `json:"focus"`
	Style           string `json:"style,omitempty"`
}

// MeditationScript represents the generated meditation content
type MeditationScript struct {
	Content       string            `json:"content"`
	TimingMarkers map[string]string `json:"timing_markers"`
}

// ScriptGenerator handles the generation of meditation scripts
type ScriptGenerator struct {
	client *openai.Client
}

// NewScriptGenerator creates a new ScriptGenerator instance
func NewScriptGenerator(apiKey string) *ScriptGenerator {
	return &ScriptGenerator{
		client: openai.NewClient(apiKey),
	}
}

// GenerateScript creates a new meditation script based on the request
func (g *ScriptGenerator) GenerateScript(ctx context.Context, req MeditationRequest) (*MeditationScript, error) {
	prompt := g.buildPrompt(req)

	resp, err := g.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `You are an experienced meditation guide creating guided meditations. 
                    Format your response as JSON with two fields:
                    - "content": the meditation script with [PAUSE X] markers for pauses in seconds
                    - "timing_markers": a map of specific points in the meditation ("intro", "body", "closing")
                    Use a calm, soothing tone. Include clear breathing instructions and gentle guidance.
                    Start with a brief introduction, then guide the breathing, then the main practice, and end with a gentle closing.`,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate script: %w", err)
	}

	// Parse the response into our struct
	var script MeditationScript
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &script); err != nil {
		return nil, fmt.Errorf("failed to parse script: %w", err)
	}

	return &script, nil
}

// buildPrompt creates the prompt for the AI based on the request
func (g *ScriptGenerator) buildPrompt(req MeditationRequest) string {
	return fmt.Sprintf(`Create a %d-minute guided meditation focused on %s.
    Include appropriate pauses marked with [PAUSE X] where X is the number of seconds.
    The total of all meditation content and pauses should add up to approximately %d minutes.
    
    Key requirements:
    1. Start with a 15-second settling in period
    2. Include regular breathing guidance
    3. Focus specifically on %s
    4. End with a gentle return to awareness
    5. Use [PAUSE X] markers for silence periods
    6. Total length should be %d minutes
    7. Mark sections with specific timing markers for intro, body, and closing
    
    Style notes:
    - Use a gentle, calming tone
    - Include specific sensory guidance
    - Give clear but soft instructions
    - Use present-tense language
    - Avoid complex terminology`,
		req.DurationMinutes,
		req.Focus,
		req.DurationMinutes,
		req.Focus,
		req.DurationMinutes)
}
