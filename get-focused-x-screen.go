package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/vcraescu/go-xrandr"
)

type swayOutput struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Rect struct {
		X      int `json:"x"`
		Y      int `json:"y"`
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"rect"`
	Focus              []int   `json:"focus"`
	Border             string  `json:"border"`
	CurrentBorderWidth int     `json:"current_border_width"`
	Layout             string  `json:"layout"`
	Orientation        string  `json:"orientation"`
	Percent            float64 `json:"percent"`
	WindowRect         struct {
		X      int `json:"x"`
		Y      int `json:"y"`
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"window_rect"`
	DecoRect struct {
		X      int `json:"x"`
		Y      int `json:"y"`
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"deco_rect"`
	Geometry struct {
		X      int `json:"x"`
		Y      int `json:"y"`
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"geometry"`
	Window             interface{}   `json:"window"`
	Urgent             bool          `json:"urgent"`
	Marks              []interface{} `json:"marks"`
	FullscreenMode     int           `json:"fullscreen_mode"`
	Nodes              []interface{} `json:"nodes"`
	FloatingNodes      []interface{} `json:"floating_nodes"`
	Sticky             bool          `json:"sticky"`
	Type               string        `json:"type"`
	Active             bool          `json:"active"`
	Dpms               bool          `json:"dpms"`
	Primary            bool          `json:"primary"`
	Make               string        `json:"make"`
	Model              string        `json:"model"`
	Serial             string        `json:"serial"`
	Scale              float64       `json:"scale"`
	ScaleFilter        string        `json:"scale_filter"`
	Transform          string        `json:"transform"`
	AdaptiveSyncStatus string        `json:"adaptive_sync_status"`
	CurrentWorkspace   string        `json:"current_workspace"`
	Modes              []struct {
		Width   int `json:"width"`
		Height  int `json:"height"`
		Refresh int `json:"refresh"`
	} `json:"modes"`
	CurrentMode struct {
		Width   int `json:"width"`
		Height  int `json:"height"`
		Refresh int `json:"refresh"`
	} `json:"current_mode"`
	MaxRenderTime   int    `json:"max_render_time"`
	Focused         bool   `json:"focused"`
	SubpixelHinting string `json:"subpixel_hinting"`
}

type xrandrOutput struct {
	Name   string
	X      int
	Y      int
	Width  int
	Height int
}

func swayFocusedOutput() (*swayOutput, error) {
	cmd := exec.Command("swaymsg", "-t", "get_outputs")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to get focused sway output: %w", err)
	}

	var outputs []swayOutput
	if err := json.NewDecoder(&stdout).Decode(&outputs); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	for _, output := range outputs {
		if !output.Active {
			continue
		}
		if output.Focused {
			return &output, nil
		}
	}

	return nil, nil
}

func main() {
	swayOutput, err := swayFocusedOutput()
	if err != nil {
		log.Fatalf("failed to get sway outputs: %v", err)
	}
	log.Printf("found focused output: %s '%s %s %s'", swayOutput.Name, swayOutput.Make, swayOutput.Model, swayOutput.Serial)

	screens, err := xrandr.GetScreens()
	if err != nil {
		log.Fatalf("failed to get xrandr outputs: %v", err)
	}

	for _, screen := range screens {
		for _, monitor := range screen.Monitors {
			if monitor.Position.X == swayOutput.Rect.X && monitor.Position.Y == swayOutput.Rect.Y {
				log.Printf("matched X output %s", monitor.ID)
				fmt.Println(monitor.ID)
				return
			}
		}
	}

	log.Fatalf("No matching monitor found")
}
