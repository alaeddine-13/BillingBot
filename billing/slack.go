package main

import (
	"fmt"

	"github.com/nlopes/slack"
)


func reportMonthlyThresholdExceeded(costResult *CostResult, threshold float32) {

	// topServicesText := build
	attachments := []slack.Attachment{
		{
			Color: "#ff0000",
			Title: fmt.Sprintf("This month's cost: $%f", costResult.Total),
			Text:  fmt.Sprintf(
					"The threshold of $%f was exceeded by resources in the project.\nCosts per day:\n%+v",
					threshold,
					costResult,
				),
		},
	}

	msg := &slack.WebhookMessage{
		Text:        fmt.Sprintf("AWS monthly cost exceeded threshold (set to *$%f*)", threshold),
		Attachments: attachments,
	}


	slack.PostWebhook(getWebhookURL(), msg)
	return
}
