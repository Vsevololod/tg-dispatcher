package processors

import "tg-dispatcher/domain"

type MessageProcessStrategy interface {
	GetName() string
	GetDescription() string
	Process(update domain.Update) bool
	CanProcess(update domain.Update) bool
}
