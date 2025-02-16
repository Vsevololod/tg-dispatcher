package service

import (
	"log/slog"
	"tg-dispatcher/domain"
	"tg-dispatcher/service/processors"
)

// MessageProcessService — сервис обработки сообщений
type MessageProcessService struct {
	inputMessageChannel  chan domain.Update
	outputMessageChannel chan domain.MessageReq
	allProcessors        []*processors.MessageProcessContext
	log                  *slog.Logger
}

// NewMessageProcessService создает новый сервис и принимает канал для сообщений
func NewMessageProcessService(inputMessageChannel chan domain.Update, outputMessageChannel chan domain.MessageReq, allProcessors []*processors.MessageProcessContext, logger *slog.Logger) *MessageProcessService {
	return &MessageProcessService{
		inputMessageChannel:  inputMessageChannel,
		outputMessageChannel: outputMessageChannel,
		allProcessors:        allProcessors,
		log:                  logger}
}

// StartProcessing запускает обработку сообщений в отдельной горутине
func (s *MessageProcessService) StartProcessing(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			for msg := range s.inputMessageChannel {
				s.ProcessMessage(workerID, msg)
			}
		}(i)
	}
}

// ProcessMessage выполняет обработку сообщения
func (s *MessageProcessService) ProcessMessage(workerID int, msg domain.Update) {
	s.log.Info("Worker %d: Обрабатываю сообщение с ID %d от %s: %s",
		workerID, msg.UpdateID, msg.Message.From.FirstName, msg.Message.Text)

	for _, processor := range s.allProcessors {
		if processor.CanProcess(msg) {
			processor.Process(msg)
		}
	}
}
