package util

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/crypto/blake2b"
)

// HashBase64 generate a composite unique index value.
func HashBase64(prefix string, ss ...string) string {
	var b bytes.Buffer
	b.WriteString(prefix)
	for _, s := range ss {
		b.WriteString(s)
	}
	h := blake2b.Sum256(b.Bytes())
	return base64.RawURLEncoding.EncodeToString(h[:])
}

// UnixMS returns a Unix time, the number of milliseconds elapsed since January 1, 1970 UTC.
func UnixMS(ts ...time.Time) int64 {
	var t time.Time
	if len(ts) == 0 {
		t = time.Now()
	} else {
		t = ts[0]
	}
	t = t.UTC().Truncate(time.Millisecond)
	return t.Unix()*1000 + int64(t.Nanosecond()/1e6)
}

// StringsHas ...
func StringsHas(ss []string, filter func(s string) bool) bool {
	for _, s := range ss {
		if filter(s) {
			return true
		}
	}
	return false
}

// StringsFilter ...
func StringsFilter(ss []string, filter func(s string) bool) []string {
	data := make([]string, 0, len(ss))
	for _, s := range ss {
		if filter(s) {
			data = append(data, s)
		}
	}
	return data
}

// FormatUID ...
func FormatUID(uid string) string {
	return fmt.Sprintf("<%s>", uid)
}

// FormatUIDs ...
func FormatUIDs(uids []string) []string {
	data := make([]string, len(uids))
	for i, uid := range uids {
		data[i] = fmt.Sprintf("<%s>", uid)
	}
	return data
}

// FormatStr ...
func FormatStr(s string) string {
	return strconv.Quote(s)
}

// FormatStrs ...
func FormatStrs(ss []string) []string {
	data := make([]string, len(ss))
	for i, s := range ss {
		data[i] = strconv.Quote(s)
	}
	return data
}
