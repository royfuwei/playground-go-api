package tools

import "time"

// String return string ref
func String(v string) *string { return &v }

// Bool return bool ref
func Bool(v bool) *bool { return &v }

// Int return int ref
func Int(v int) *int { return &v }

// Int64 return int64 ref
func Int64(v int64) *int64 { return &v }

func IntToInt32(v int) *int32 {
	v32 := int32(v)
	return &v32
}

// TimeParseByFormats 不確定隨錄欄位時間format，比對timeFormats，返回對應的time.Time
func TimeParseByFormats(timeString string) (time.Time, error) {
	var result time.Time
	var err error
	customTimeFormat := "2006-01-02 15:04:05MST"
	timeFormats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		customTimeFormat,
	}
	for _, format := range timeFormats {
		if format == customTimeFormat {
			zone, _ := time.Now().Zone()
			timeString = timeString + zone
		}
		result, err = time.Parse(format, timeString)
		if err == nil {
			break
		}
	}
	return result, err
}
