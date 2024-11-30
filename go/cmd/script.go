package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/systay/mindful-ai/go/script_generator"
	"github.com/systay/mindful-ai/go/utils"
	"os"
)

func scriptCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "script",
		Example: "mindful script",
		Short:   "Create a meditation script",
		Args:    cobra.NoArgs,
		Run: func(_ *cobra.Command, args []string) {
			req := script_generator.MeditationRequest{
				DurationMinutes: 10,
				Focus:           "breath",
				Style:           "slow",
			}
			ctx := context.Background()
			err := utils.LoadEnv()
			fail(err, "loading ENV")
			apiKey := os.Getenv("OPENAI_API_KEY")
			generator := script_generator.NewScriptGenerator(apiKey)
			script, err := generator.GenerateScript(ctx, req)
			fail(err, "generating script")
			fmt.Println(script.ToString())
		},
	}
}
