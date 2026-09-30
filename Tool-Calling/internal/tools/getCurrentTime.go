package tools

import (
	"fmt"
	"time"

	"google.golang.org/genai"
)

type CurrentTimeArgs struct {
	Timezone string `json:"timezone"`
}

type CurrentTimeResult struct {
	Timezone    string `json:"timezone"`
	CurrentTime string `json:"time"`
}

func GetCurrentTime(args CurrentTimeArgs) (CurrentTimeResult, error) {
	if args.Timezone == "" {
		return CurrentTimeResult{}, fmt.Errorf("Time zone not found")
	}

	location, err := time.LoadLocation(args.Timezone)
	if err != nil {
		return CurrentTimeResult{}, err
	}

	timeZone := CurrentTimeResult{
		Timezone:    args.Timezone,
		CurrentTime: time.Now().In(location).Format(time.RFC3339),
	}

	return timeZone, nil
}

func CurrentTimeDeclaration() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        "get_current_time",
		Description: "Get the current date and time in an IANA timezone." + "Use when the other user asks for the current time" + "If the location is unclear ask the user to clarify",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"timezone": {
					Type: genai.TypeString,
					Description: "IANA timezone, such as " +
						"Asia/Kolkata or America/New_York",
				},
			},
			Required: []string{"timezone"},
		},
	}
}
