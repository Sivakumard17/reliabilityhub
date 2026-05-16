package incidents

import "time"

// AlertManagerPayload is the webhook body sent by AlertManager
// Follows the AlertManager webhook v4 format
type AlertManagerPayload struct {
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	TruncatedAlerts   int               `json:"truncatedAlerts"`
	Status            string            `json:"status"` // "firing" or "resolved"
	Receiver          string            `json:"receiver"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Alerts            []Alert           `json:"alerts"`
}

// Alert is a single alert within the payload
type Alert struct {
	Status       string            `json:"status"` // "firing" or "resolved"
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

// severityFromLabels maps AlertManager severity labels to our Severity type
func SeverityFromLabels(labels map[string]string) Severity {
	switch labels["severity"] {
	case "critical":
		return SeverityCritical
	case "high":
		return SeverityHigh
	case "warning":
		return SeverityMedium
	case "info":
		return SeverityInfo
	default:
		return SeverityMedium
	}
}

// titleFromAlert generates a human-readable incident title
func TitleFromAlert(alert Alert) string {
	if summary, ok := alert.Annotations["summary"]; ok {
		return summary
	}
	if alertname, ok := alert.Labels["alertname"]; ok {
		return alertname
	}
	return "Unknown Alert"
}

// serviceFromAlert extracts service name from alert labels
func ServiceFromAlert(alert Alert) string {
	for _, key := range []string{"service", "job", "app", "deployment"} {
		if v, ok := alert.Labels[key]; ok {
			return v
		}
	}
	return "unknown"
}
