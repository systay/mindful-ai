package tts

import (
	"github.com/systay/mindful-ai/go/utils"
	"os"
	"testing"
)
import "github.com/stretchr/testify/require"

func TestName(t *testing.T) {
	err := utils.LoadEnv()
	require.NoError(t, err)
	apiKey := os.Getenv("ELEVENLABS_API_KEY")

	script := `Welcome to this focused attention meditation, where we will concentrate on our breath and any discomfort you may be feeling in your neck. [PAUSE 5] Begin by finding a comfortable seated position. Sit upright, but try not to strain any muscles. Rest your hands in your lap or on your knees. [PAUSE 5] Close your eyes. [PAUSE 5]
'intro': 'As we begin our journey, let's start by bringing our attention to our breath. [PAUSE 5] Notice how it feels as you inhale and exhale. Don't try to change your breathing, just observe it. [PAUSE 5]
'body': 'Now, gently shift your focus to your neck. [PAUSE 5] Notice any tension or discomfort that you feel. [PAUSE 5] Imagine that each breath you take is flowing into your neck. [PAUSE 5] As you exhale, imagine that the tension is leaving your body with your breath. [PAUSE 5] Continue this practice, focusing on your breath and releasing tension in your neck. [PAUSE 5] If your mind starts to wander, gently bring your attention back to your breath and your neck. [PAUSE 5]
'closing': 'Now, slowly start to bring your focus back to your surroundings. [PAUSE 5] Wiggle your fingers and your toes. [PAUSE 5] When you're ready, slowly open your eyes. [PAUSE 5] Thank you for joining me in this focused attention meditation. [PAUSE 5]`

	err = TextToSpeech(script, "output_audio.mp3", apiKey)
	require.NoError(t, err)
}
