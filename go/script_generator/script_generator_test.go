package script_generator

import (
	"context"
	"fmt"
	"os"
	"testing"
)
import "github.com/stretchr/testify/require"

func TestName(t *testing.T) {
	err := LoadTestEnv()
	require.NoError(t, err)
	apiKey := os.Getenv("OPENAI_API_KEY")
	generator := NewScriptGenerator(apiKey)
	script, err := generator.GenerateScript(context.Background(), MeditationRequest{
		DurationMinutes: 5,
		Focus:           "first on breathing, and then shift to focus on hearing sounds",
		Style:           "slow with long pauses. should be mostly quiet, with a little guidence here and there",
	})

	require.NoError(t, err)
	fmt.Println(script.Content)
	fmt.Println(script.TimingMarkers)
}
