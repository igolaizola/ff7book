package ff7book

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Config struct {
	Debug  bool
	Input  string
	Output string
	Key    string
	Model  string
}

const prePrompt = `You are an award-winning novelist and professional book formatter.
Your task is to convert the following video game script into a full, polished novel,
ready for eBook and print conversion.
You must:
    1. Integrate the provided story summary and context to enrich world-building and narrative consistency.
    2. Produce valid Markdown (or HTML) with front-matter metadata.
    3. Divide the story into numbered chapters with titles.
    4. Include a Table of Contents after the front-matter.
    5. Use prose-novel style: show-don't-tell, varied sentence lengths, POV shifts only at chapter breaks.
    6. Embed scene breaks as horizontal rules ("---") where needed.
    7. No additional commentary—only the novel content.

Below is the synopsis and full script:


`

// Generate runs the ff7book process.
func Generate(ctx context.Context, cfg *Config) error {
	log.Println("running")
	defer log.Println("finished")

	// Read input script
	input, err := os.ReadFile(cfg.Input)
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}

	// Prepare the GenAI client
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.Key))
	if err != nil {
		return fmt.Errorf("couldn't create genai client: %w", err)
	}
	model := client.GenerativeModel(cfg.Model)

	// Open (or create) output file for streaming writes
	f, err := os.OpenFile(cfg.Output,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0o644,
	)
	if err != nil {
		return fmt.Errorf("open output file: %w", err)
	}
	defer f.Close()

	// Build prompt and start stream
	prompt := prePrompt + string(input)
	stream := model.GenerateContentStream(ctx, genai.Text(prompt))
	if stream == nil {
		return fmt.Errorf("failed to start stream")
	}

	// Stream loop: write each chunk to both stdout AND file immediately
	for {
		resp, err := stream.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			log.Printf("stream error: %v", err)
			break
		}
		for _, part := range resp.Candidates[0].Content.Parts {
			txt, ok := part.(genai.Text)
			if !ok {
				log.Printf("unexpected part type: %T", part)
				continue
			}
			chunk := string(txt)

			// 1) real-time stdout feedback
			fmt.Print(chunk)

			// 2) append to file
			if _, err := f.WriteString(chunk); err != nil {
				log.Printf("file write error: %v", err)
			}
		}
	}

	return nil
}
