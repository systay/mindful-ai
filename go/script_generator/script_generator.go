package script_generator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// TechniqueType represents the allowed meditation techniques as an integer
type TechniqueType int

const (
	BodyScan           TechniqueType = iota // 0
	FocusedAttention                        // 1
	LovingKindness                          // 2
	MindfulnessEmotion                      // 3
	GratitudePractice                       // 4
)

// techniqueNames maps TechniqueType values to their string equivalents
var techniqueNames = map[TechniqueType]string{
	BodyScan:           "body_scan",
	FocusedAttention:   "focused_attention",
	LovingKindness:     "loving_kindness",
	MindfulnessEmotion: "mindfulness_emotion",
	GratitudePractice:  "gratitude_practice",
}

// techniqueValues maps string names to their TechniqueType equivalents
var techniqueValues = map[string]TechniqueType{
	"body_scan":           BodyScan,
	"focused_attention":   FocusedAttention,
	"loving_kindness":     LovingKindness,
	"mindfulness_emotion": MindfulnessEmotion,
	"gratitude_practice":  GratitudePractice,
}

// MarshalJSON implements custom JSON marshalling for TechniqueType
func (t TechniqueType) MarshalJSON() ([]byte, error) {
	name, exists := techniqueNames[t]
	if !exists {
		return nil, errors.New("invalid TechniqueType")
	}
	return json.Marshal(name)
}

// UnmarshalJSON implements custom JSON unmarshalling for TechniqueType
func (t *TechniqueType) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}
	val, exists := techniqueValues[name]
	if !exists {
		return errors.New("invalid technique name")
	}
	*t = val
	return nil
}

type MeditationRequest struct {
	Technique         TechniqueType `json:"technique"`                    // e.g., "body_scan", "focused_attention"
	SessionLength     int           `json:"session_length"`               // in minutes
	GuidanceLevel     string        `json:"guidance_level"`               // e.g., "detailed", "brief"
	FocusObject       string        `json:"focus_object,omitempty"`       // e.g., "breath", "mantra", for techniques like focused attention
	CompassionTargets []string      `json:"compassion_targets,omitempty"` // Targets for loving-kindness meditation, e.g., ["self", "family", "all_beings"]
	EmotionLabels     []string      `json:"emotion_labels,omitempty"`     // Labels for mindfulness of emotions, e.g., ["joy", "anger"]
	GratitudeScope    string        `json:"gratitude_scope,omitempty"`    // e.g., "self", "others", "world"
	AmbientSound      string        `json:"ambient_sound,omitempty"`      // e.g., "nature", "silence", "white_noise"
	VoicePreference   string        `json:"voice_preference,omitempty"`   // e.g., "calm", "whisper"
	Goal              string        `json:"goal,omitempty"`               // e.g., "relaxation", "focus"
}

// MeditationScript represents the generated meditation content
type MeditationScript struct {
	Content       string            `json:"content"`
	TimingMarkers map[string]string `json:"timing_markers"`
}

func (s MeditationScript) ToString() string {
	return fmt.Sprintf("Content: %s\nTimingMarkers: %v", s.Content, s.TimingMarkers)
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
	prompt, err := BuildMeditationPrompt(req)
	if err != nil {
		return nil, err
	}

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
                    - "timing_markers": a map of specific points in the meditation ("intro", "body", "closing")`,
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
	content := resp.Choices[0].Message.Content
	content = strings.Replace(content, "\n", "\\n", 0)
	if err := json.Unmarshal([]byte(content), &script); err != nil {
		return nil, fmt.Errorf("failed to parse script: %w\ncontent: %s", err, content)
	}

	return &script, nil
}

func BuildMeditationPrompt(request MeditationRequest) (string, error) {
	switch request.Technique {
	case BodyScan:
		return BuildBodyScanPrompt(request), nil
	case FocusedAttention:
		return BuildFocusedAttentionPrompt(request), nil
	case LovingKindness:
		return BuildLovingKindnessPrompt(request), nil
	case MindfulnessEmotion:
		return BuildMindfulnessEmotionPrompt(request), nil
	case GratitudePractice:
		return BuildGratitudePracticePrompt(request), nil
	default:
		return "", fmt.Errorf("unsupported technique")
	}
}

// Functions to build prompts for each meditation technique

func BuildBodyScanPrompt(request MeditationRequest) string {
	prompt := fmt.Sprintf(
		`You are a meditation teacher. Create a %d-minute body scan meditation script with %s guidance. The goal is %s. Use a %s voice tone. Include references to %s ambient sounds. Follow traditional mindfulness practices for body scan meditation.`,
		request.SessionLength,
		request.GuidanceLevel,
		request.Goal,
		request.VoicePreference,
		request.AmbientSound,
	)
	return prompt
}

func BuildFocusedAttentionPrompt(request MeditationRequest) string {
	focusObject := request.FocusObject
	if focusObject == "" {
		focusObject = "the breath"
	}

	if request.GuidanceLevel == "" {
		request.GuidanceLevel = "brief"
	}
	if request.VoicePreference == "" {
		request.VoicePreference = "calm"
	}

	prompt := fmt.Sprintf(
		`You are a meditation teacher. Create a %d-minute focused attention meditation script where the focus is on %s. Provide %s guidance. The goal is %s. Incorporate traditional practices of focused attention meditation.`,
		request.SessionLength,
		focusObject,
		request.GuidanceLevel,
		request.Goal,
	)
	return prompt
}

func BuildLovingKindnessPrompt(request MeditationRequest) string {
	targets := "self and others"
	if len(request.CompassionTargets) > 0 {
		targets = strings.Join(request.CompassionTargets, ", ")
	}
	prompt := fmt.Sprintf(
		`You are a meditation teacher. Create a %d-minute loving-kindness (metta) meditation script focusing on cultivating compassion towards %s. Provide %s guidance. Use a %s voice tone and include references to %s ambient sounds. Follow traditional loving-kindness meditation practices.`,
		request.SessionLength,
		targets,
		request.GuidanceLevel,
		request.VoicePreference,
		request.AmbientSound,
	)
	return prompt
}

func BuildMindfulnessEmotionPrompt(request MeditationRequest) string {
	emotions := "various emotions"
	if len(request.EmotionLabels) > 0 {
		emotions = strings.Join(request.EmotionLabels, ", ")
	}
	prompt := fmt.Sprintf(
		`You are a meditation teacher. Create a %d-minute mindfulness of emotions meditation script, guiding the listener to observe and acknowledge emotions such as %s. Provide %s guidance. The goal is %s. Use a %s voice tone and include references to %s ambient sounds. Incorporate traditional mindfulness practices.`,
		request.SessionLength,
		emotions,
		request.GuidanceLevel,
		request.Goal,
		request.VoicePreference,
		request.AmbientSound,
	)
	return prompt
}

func BuildGratitudePracticePrompt(request MeditationRequest) string {
	scope := request.GratitudeScope
	if scope == "" {
		scope = "self, others, and the world"
	}
	prompt := fmt.Sprintf(
		`You are a meditation teacher. Create a %d-minute gratitude meditation script focusing on cultivating gratitude towards %s. Provide %s guidance. The goal is %s. Use a %s voice tone and include references to %s ambient sounds. Incorporate traditional gratitude meditation practices.`,
		request.SessionLength,
		scope,
		request.GuidanceLevel,
		request.Goal,
		request.VoicePreference,
		request.AmbientSound,
	)
	return prompt
}
