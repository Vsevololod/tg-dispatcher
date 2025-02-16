package sl

import (
	"log/slog"
	"tg-dispatcher/domain"
)

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func Req(req domain.Update) slog.Attr {
	return slog.Attr{
		Key:   "req",
		Value: slog.AnyValue(req),
	}
}
