package events

import (
    "encoding/json"
    "time"

    "github.com/google/uuid"
)

type Event struct {
    ID        string    `json:"id"`
    Type      string    `json:"type"`
    Source    string    `json:"source"`
    Version   string    `json:"version"`
    Timestamp time.Time `json:"timestamp"`
    Data      []byte    `json:"data"`
}

func NewEvent(eventType, source, version string, data interface{}) (*Event, error) {
    dataBytes, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }

    return &Event{
        ID:        uuid.New().String(),
        Type:      eventType,
        Source:    source,
        Version:   version,
        Timestamp: time.Now().UTC(),
        Data:      dataBytes,
    }, nil
}

const (
    UserRegisteredEvent    = "user.registered"
    UserTierUpgradedEvent  = "user.tier.upgraded"
    UserQuotaUpdatedEvent  = "user.quota.updated"
)

type UserRegisteredData struct {
    UserID    string `json:"user_id"`
    Email     string `json:"email"`
    FullName  string `json:"full_name"`
    Tier      string `json:"tier"`
    CreatedAt string `json:"created_at"`
}

type UserTierUpgradedData struct {
    UserID    string `json:"user_id"`
    OldTier   string `json:"old_tier"`
    NewTier   string `json:"new_tier"`
    UpgradedAt string `json:"upgraded_at"`
}

type QuotaData struct {
    Used  int `json:"used"`
    Limit int `json:"limit"`
}

type QuotaInfoData struct {
    AIDescription QuotaData `json:"ai_description"`
    AIVideo       QuotaData `json:"ai_video"`
    AutoPosting   QuotaData `json:"auto_posting"`
}

type UserQuotaUpdatedData struct {
    UserID    string       `json:"user_id"`
    Quotas    QuotaInfoData `json:"quotas"`
    UpdatedAt string       `json:"updated_at"`
}
