package postgres

import (
	"database/sql"
	"encoding/json"
	"strings"
)

func columnIndexMap(cols []string) map[string]int {
	m := make(map[string]int, len(cols))
	for i, c := range cols {
		m[strings.ToLower(c)] = i
	}
	return m
}

func scanRowValues(rows *sql.Rows) ([]interface{}, []string, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	values := make([]interface{}, len(cols))
	ptrs := make([]interface{}, len(cols))
	for i := range values {
		ptrs[i] = &values[i]
	}
	if err := rows.Scan(ptrs...); err != nil {
		return nil, nil, err
	}
	return values, cols, nil
}

func valString(row []interface{}, idx map[string]int, names ...string) string {
	for _, name := range names {
		i, ok := idx[strings.ToLower(name)]
		if !ok || row[i] == nil {
			continue
		}
		switch v := row[i].(type) {
		case string:
			return v
		case []byte:
			return string(v)
		}
	}
	return ""
}

func valInt64(row []interface{}, idx map[string]int, names ...string) (int64, bool) {
	for _, name := range names {
		i, ok := idx[strings.ToLower(name)]
		if !ok || row[i] == nil {
			continue
		}
		switch v := row[i].(type) {
		case int64:
			return v, true
		case int32:
			return int64(v), true
		case int:
			return int64(v), true
		case float64:
			return int64(v), true
		}
	}
	return 0, false
}

func valFloat64(row []interface{}, idx map[string]int, names ...string) (*float64, bool) {
	for _, name := range names {
		i, ok := idx[strings.ToLower(name)]
		if !ok || row[i] == nil {
			continue
		}
		switch v := row[i].(type) {
		case float64:
			f := v
			return &f, true
		case float32:
			f := float64(v)
			return &f, true
		case int64:
			f := float64(v)
			return &f, true
		}
	}
	return nil, false
}

func valBool(row []interface{}, idx map[string]int, names ...string) (*bool, bool) {
	for _, name := range names {
		i, ok := idx[strings.ToLower(name)]
		if !ok || row[i] == nil {
			continue
		}
		switch v := row[i].(type) {
		case bool:
			return &v, true
		}
	}
	return nil, false
}

func valJSONRaw(row []interface{}, idx map[string]int, names ...string) json.RawMessage {
	s := valString(row, idx, names...)
	if s == "" || s == "{}" || s == "null" {
		return nil
	}
	return json.RawMessage(s)
}

func int64Ptr(v int64, ok bool) *int64 {
	if !ok {
		return nil
	}
	return &v
}
