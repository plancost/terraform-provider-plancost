package logging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type TfLogAdapter struct{}

func (w TfLogAdapter) Write(p []byte) (n int, err error) {
	var event map[string]interface{}
	if err := json.Unmarshal(p, &event); err != nil {
		// Fallback if not JSON
		tflog.Debug(context.Background(), string(p))
		return len(p), nil //nolint:nilerr
	}

	if event == nil {
		return len(p), nil
	}

	level := event["level"]
	msg := event["message"]

	// Remove fields handled by tflog or not needed
	delete(event, "level")
	delete(event, "message")
	delete(event, "time")

	switch level {
	case "trace":
		log.Println("[TRACE]", msg, event)
	case "debug":
		log.Println("[DEBUG]", msg, event)
	case "info":
		log.Println("[INFO]", msg, event)
	case "warn":
		log.Println("[WARN]", msg, event)
	case "error":
		log.Println("[ERROR]", msg, event)
	default:
		log.Println("[DEBUG]", msg, event)
	}
	return len(p), nil
}
