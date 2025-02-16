package processors

import (
	"log/slog"
	"strings"
	"tg-dispatcher/domain"
)

type PlaylistProcessStrategy struct {
	log *slog.Logger
}

func (s PlaylistProcessStrategy) GetName() string {
	return "today_playlist"
}

func (s PlaylistProcessStrategy) GetDescription() string {
	return "Сгенерировать плейлист для сегодняшних подкастов"
}
func (s PlaylistProcessStrategy) Process(update domain.Update) bool {
	s.log.Info("Process playlist", update)
	return true
}
func (s PlaylistProcessStrategy) CanProcess(update domain.Update) bool {
	return strings.Contains(update.Message.Text, s.GetName())
}
