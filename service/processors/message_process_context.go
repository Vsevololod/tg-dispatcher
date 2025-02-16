package processors

import (
	"log/slog"
	"tg-dispatcher/domain"
)

type MessageProcessContext struct {
	strategy MessageProcessStrategy
}

func create(strategy MessageProcessStrategy) *MessageProcessContext {
	return &MessageProcessContext{strategy: strategy}
}

func (c MessageProcessContext) Process(update domain.Update) bool {
	return c.strategy.Process(update)
}

func (c MessageProcessContext) CanProcess(update domain.Update) bool {
	return c.strategy.CanProcess(update)
}

func CreateAllProcessors(outputMessageChannel chan domain.MessageReq, videoProvider VideoProvider, videoSaver VideoSaver, logger *slog.Logger) []*MessageProcessContext {
	return []*MessageProcessContext{
		create(&UrlProcessStrategy{outputMessageChannel, videoProvider, videoSaver, logger}),
		create(&PlaylistProcessStrategy{logger}),
	}
}
