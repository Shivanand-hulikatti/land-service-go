package service

import (
	"encoding/json"
	"fmt"
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
		normalizeUserMapDates(m, dobFormat)
	}
}

func normalizeUserMapDates(m map[string]any, dobFormat string) {
	replaceDateField(m, "createdDate", userTimestampLayout)
	replaceDateField(m, "lastModifiedDate", userTimestampLayout)
	replaceStringField(m, "createdBy")
	replaceStringField(m, "lastModifiedBy")
	if dobFormat != "" {
		replaceDateField(m, "dob", dobFormat)
	}
	replaceDateField(m, "pwdExpiryDate", userTimestampLayout)
}

func replaceDateField(m map[string]any, field, layout string) {
	s, ok := m[field].(string)
	if !ok || s == "" {
		return
	}
	m[field] = dateToLong(s, layout)
}

func dateToLong(value, layout string) *int64 {
	t, err := time.Parse(layout, value)
	if err != nil {
		return nil
	}
	ms := t.UnixMilli()
	return &ms
}

func replaceStringField(m map[string]any, field string) {
	switch v := m[field].(type) {
	case float64:
		m[field] = fmt.Sprintf("%.0f", v)
	case int:
		m[field] = fmt.Sprintf("%d", v)
	case int64:
		m[field] = fmt.Sprintf("%d", v)
	}
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
