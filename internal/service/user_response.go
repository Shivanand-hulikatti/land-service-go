package service

import (
	"encoding/json"
	"time"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/domain"
)

const (
	userTimestampLayout = "02-01-2006 15:04:05"
	userDOBSearchLayout = "2006-01-02"
	userDOBCreateLayout = "02/01/2006"
)

// normalizeUserResponseDates converts user service date strings to epoch millis (Java parseResponse parity).
func normalizeUserResponseDates(resp map[string]any, dobFormat string) {
	users, ok := resp["user"].([]any)
	if !ok {
		return
	}
	for _, u := range users {
		m, ok := u.(map[string]any)
		if !ok {
			continue
		}
		if s, ok := m["createdDate"].(string); ok && s != "" {
			m["createdDate"] = dateToLong(s, userTimestampLayout)
		}
		if s, ok := m["lastModifiedDate"].(string); ok && s != "" {
			m["lastModifiedDate"] = dateToLong(s, userTimestampLayout)
		}
		if s, ok := m["dob"].(string); ok && s != "" && dobFormat != "" {
			m["dob"] = dateToLong(s, dobFormat)
		}
		if s, ok := m["pwdExpiryDate"].(string); ok && s != "" {
			m["pwdExpiryDate"] = dateToLong(s, userTimestampLayout)
		}
	}
}

func dateToLong(value, layout string) *int64 {
	t, err := time.Parse(layout, value)
	if err != nil {
		return nil
	}
	ms := t.UnixMilli()
	return &ms
}

func decodeUserDetailResponse(resp map[string]any, dobFormat string) (*domain.UserDetailResponse, error) {
	normalizeUserResponseDates(resp, dobFormat)
	raw, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	var out domain.UserDetailResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
