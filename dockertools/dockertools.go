package dockertools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/ronyv89/temple/structs"
)

type pullEvent struct {
	ID             string `json:"id"`
	Status         string `json:"status"`
	Error          string `json:"error,omitempty"`
	Progress       string `json:"progress,omitempty"`
	ProgressDetail struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progressDetail"`
}

func PullImage(imageName string) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	resp, err := cli.ImagePull(context.TODO(), "ruby", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	cursor := structs.Cursor{}
	layers := make([]string, 0)
	oldIndex := len(layers)

	var event *pullEvent
	decoder := json.NewDecoder(resp)
	cursor.Hide()

	for {
		if err := decoder.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}

			panic(err)
		}

		imageID := event.ID

		// Check if the line is one of the final two ones
		if strings.HasPrefix(event.Status, "Digest:") || strings.HasPrefix(event.Status, "Status:") {
			fmt.Printf("%s\n", event.Status)
			continue
		}

		// Check if ID has already passed once
		index := 0
		for i, v := range layers {
			if v == imageID {
				index = i + 1
				break
			}
		}

		// Move the cursor
		if index > 0 {
			diff := index - oldIndex

			if diff > 1 {
				down := diff - 1
				cursor.MoveDown(down)
			} else if diff < 1 {
				up := diff*(-1) + 1
				cursor.MoveUp(up)
			}

			oldIndex = index
		} else {
			layers = append(layers, event.ID)
			diff := len(layers) - oldIndex

			if diff > 1 {
				cursor.MoveDown(diff) // Return to the last row
			}

			oldIndex = len(layers)
		}

		cursor.ClearLine()

		if event.Status == "Pull complete" {
			fmt.Printf("%s: %s\n", event.ID, event.Status)
		} else {
			fmt.Printf("%s: %s %s\n", event.ID, event.Status, event.Progress)
		}

	}

	cursor.Show()

	if strings.Contains(event.Status, fmt.Sprintf("Downloaded newer image for %s", "rails")) {
		return
	}
}
