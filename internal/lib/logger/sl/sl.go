package sl

import "log/slog"

func Error(err error) slog.Attr {
	return slog.Attr{"error", slog.StringValue(err.Error())}

}
