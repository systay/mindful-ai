package script_generator

import (
	"context"
	"fmt"
	"github.com/systay/mindful-ai/go/utils"
	"os"
	"testing"
)
import "github.com/stretchr/testify/require"

func TestName(t *testing.T) {
	err := utils.LoadEnv()
	require.NoError(t, err)
	apiKey := os.Getenv("OPENAI_API_KEY")
	generator := NewScriptGenerator(apiKey)

	req := MeditationRequest{
		Technique:     FocusedAttention,
		SessionLength: 12,
		GuidanceLevel: "brief",
		FocusObject:   "breath and pain in neck",
		Goal:          "calm and focus",
	}
	script, err := generator.GenerateScript(context.Background(), req)

	require.NoError(t, err)
	fmt.Println(script.Content)
	fmt.Println(script.TimingMarkers)
}
