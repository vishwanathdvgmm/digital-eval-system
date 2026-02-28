// Package logger provides a custom colorful log formatter for the Go node.
// It produces well-aligned, colored terminal output with clear component
// labels, method tags, and spacing.
package logger

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// ANSI color codes
const (
	reset = "\033[0m"
	bold  = "\033[1m"
	dim   = "\033[2m"

	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"

	bgRed   = "\033[41m"
	bgGreen = "\033[42m"
	bgBlue  = "\033[44m"
	bgCyan  = "\033[46m"
)

// PrettyFormatter is a logrus Formatter producing colourful, well-spaced logs.
type PrettyFormatter struct{}

func (f *PrettyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buf bytes.Buffer

	// ── Timestamp ────────────────────────────────────────
	ts := entry.Time.Format("15:04:05")
	buf.WriteString(fmt.Sprintf("%s%s%s ", dim, ts, reset))

	// ── Level badge ──────────────────────────────────────
	levelBadge := formatLevel(entry.Level)
	buf.WriteString(levelBadge)
	buf.WriteString("  ")

	// ── Component tag (if msg starts with [Component]) ───
	msg := entry.Message
	if strings.HasPrefix(msg, "[") {
		end := strings.Index(msg, "]")
		if end > 0 {
			component := msg[1:end]
			msg = strings.TrimSpace(msg[end+1:])
			buf.WriteString(fmt.Sprintf("%s%s %-18s %s ", bold, cyan, component, reset))
			if msg != "" {
				buf.WriteString(dim)
				buf.WriteString("│ ")
				buf.WriteString(reset)
				buf.WriteString(msg)
			}
			appendFields(&buf, entry.Data)
			buf.WriteByte('\n')
			return buf.Bytes(), nil
		}
	}

	// ── HTTP request log (method + path + duration) ──────
	method, hasMethod := entry.Data["method"]
	path, hasPath := entry.Data["path"]
	if hasMethod && hasPath {
		methodStr := fmt.Sprintf("%v", method)
		pathStr := fmt.Sprintf("%v", path)
		buf.WriteString(formatMethod(methodStr))
		buf.WriteString(fmt.Sprintf("  %-42s", pathStr))
		if dur, ok := entry.Data["duration"]; ok {
			buf.WriteString(fmt.Sprintf("%s%v%s", dim, dur, reset))
		}
		// trace id dimmed at end
		if trace, ok := entry.Data["trace_id"]; ok {
			buf.WriteString(fmt.Sprintf("  %s%v%s", dim, trace, reset))
		}
		buf.WriteByte('\n')
		return buf.Bytes(), nil
	}

	// ── Generic message ──────────────────────────────────
	buf.WriteString(msg)
	appendFields(&buf, entry.Data)
	buf.WriteByte('\n')
	return buf.Bytes(), nil
}

func formatLevel(level logrus.Level) string {
	switch level {
	case logrus.TraceLevel, logrus.DebugLevel:
		return fmt.Sprintf("%s%s DBG %s", dim, white, reset)
	case logrus.InfoLevel:
		return fmt.Sprintf("%s%s INF %s", bold, green, reset)
	case logrus.WarnLevel:
		return fmt.Sprintf("%s%s WRN %s", bold, yellow, reset)
	case logrus.ErrorLevel:
		return fmt.Sprintf("%s%s ERR %s", bold, red, reset)
	case logrus.FatalLevel, logrus.PanicLevel:
		return fmt.Sprintf("%s%s FTL %s", bold+bgRed, white, reset)
	default:
		return fmt.Sprintf("%s%s ??? %s", dim, white, reset)
	}
}

func formatMethod(method string) string {
	padded := fmt.Sprintf("%-6s", method)
	switch strings.ToUpper(method) {
	case "GET":
		return fmt.Sprintf("%s%s%s", green, padded, reset)
	case "POST":
		return fmt.Sprintf("%s%s%s", blue, padded, reset)
	case "PUT", "PATCH":
		return fmt.Sprintf("%s%s%s", yellow, padded, reset)
	case "DELETE":
		return fmt.Sprintf("%s%s%s", red, padded, reset)
	default:
		return fmt.Sprintf("%s%s%s", white, padded, reset)
	}
}

func appendFields(buf *bytes.Buffer, data logrus.Fields) {
	// skip fields we already handled
	skip := map[string]bool{
		"method": true, "path": true, "duration": true,
		"trace_id": true, "trace": true,
	}
	first := true
	for k, v := range data {
		if skip[k] {
			continue
		}
		if first {
			buf.WriteString("  ")
			first = false
		}
		buf.WriteString(fmt.Sprintf(" %s%s%s=%v", cyan, k, reset, v))
	}
	// always append trace if present (for non-HTTP logs)
	if trace, ok := data["trace"]; ok {
		buf.WriteString(fmt.Sprintf("  %s%v%s", dim, trace, reset))
	}
}
