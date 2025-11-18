package schema

import "time"

// ActionType represents the type of alert action
type ActionType string

const (
	ActionTypeNotificationPanel ActionType = "NotificationPanel"
	ActionTypeSendEmail         ActionType = "SendEmail"
	ActionTypeSendSms           ActionType = "SendSms"
	ActionTypeAutoscalingDbaas  ActionType = "AutoscalingDbaas"
)

// ExecutedAlertAction represents an executed alert action
type ExecutedAlertAction struct {
	ActionType   ActionType `json:"actionType,omitempty"`
	Success      bool       `json:"success,omitempty"`
	ErrorMessage string     `json:"errorMessage,omitempty"`
}

// AlertAction represents a possible alert action
type AlertAction struct {
	Key        string `json:"key,omitempty"`
	Disabled   bool   `json:"disabled,omitempty"`
	Executable bool   `json:"executable,omitempty"`
}

// AlertResponse represents an alert response
type AlertResponse struct {
	ID                   string                `json:"id,omitempty"`
	EventID              string                `json:"eventId,omitempty"`
	EventName            string                `json:"eventName,omitempty"`
	Username             string                `json:"username,omitempty"`
	ServiceCategory      string                `json:"serviceCategory,omitempty"`
	ServiceTypology      string                `json:"serviceTypology,omitempty"`
	ResourceID           string                `json:"resourceId,omitempty"`
	ServiceName          string                `json:"serviceName,omitempty"`
	ResourceTypology     string                `json:"resourceTypology,omitempty"`
	Metric               string                `json:"metric,omitempty"`
	LastReception        time.Time             `json:"lastReception,omitempty"`
	Rule                 string                `json:"rule,omitempty"`
	Theshold             int64                 `json:"theshold,omitempty"`
	Um                   string                `json:"um,omitempty"`
	Duration             string                `json:"duration,omitempty"`
	ThesholdExceedence   string                `json:"thesholdExceedence,omitempty"`
	Component            string                `json:"component,omitempty"`
	ClusterTypology      string                `json:"clusterTypology,omitempty"`
	Cluster              string                `json:"cluster,omitempty"`
	Clustername          string                `json:"clustername,omitempty"`
	NodePool             string                `json:"nodePool,omitempty"`
	Sms                  bool                  `json:"sms,omitempty"`
	Email                bool                  `json:"email,omitempty"`
	Panel                bool                  `json:"panel,omitempty"`
	Hidden               bool                  `json:"hidden,omitempty"`
	ExecutedAlertActions []ExecutedAlertAction `json:"executedAlertActions,omitempty"`
	Actions              []AlertAction         `json:"actions,omitempty"`
}

// AlertsListResponse represents a list of alerts
type AlertsListResponse struct {
	ListResponse
	Values []AlertResponse `json:"values"`
}
