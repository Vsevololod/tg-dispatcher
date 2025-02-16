package processors

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"tg-dispatcher/domain"
	"tg-dispatcher/domain/models"
	"tg-dispatcher/lib"
	"tg-dispatcher/lib/logger/sl"
	"tg-dispatcher/storage"
)

type UrlProcessStrategy struct {
	outputMessageChannel chan domain.MessageReq
	videoProvider        VideoProvider
	videoSaver           VideoSaver
	log                  *slog.Logger
}

type VideoSaver interface {
	SaveVideoMin(ctx context.Context, hashId string, originalId int64, url string, videoId string, userId int64) error
}

type VideoProvider interface {
	GetVideoById(ctx context.Context, videoId string) (models.Video, error)
}

func (s UrlProcessStrategy) GetName() string {
	return "process_url"
}

func (s UrlProcessStrategy) GetDescription() string {
	return "Загрузка видео по урлу"
}

func (s UrlProcessStrategy) Process(update domain.Update) bool {
	s.log.Info("Process Url:", update)
	videoId := lib.GetVideoIdFromUrl(update.Message.Text)
	video, err := s.videoProvider.GetVideoById(context.Background(), videoId)
	if err != nil {
		if errors.Is(err, storage.ErrVideoNotFound) {
			err := s.videoSaver.SaveVideoMin(context.Background(),
				update.UUID,
				update.UpdateID,
				update.Message.Text,
				videoId,
				update.Message.From.ID,
			)
			if err != nil {
				s.log.Error("Cannot save video", sl.Err(err), sl.Req(update))
				return false
			}
			s.outputMessageChannel <- domain.MessageReq{
				UUID:        update.UUID,
				Destination: domain.VideoDownload,
				Message: domain.VideoDownloadReq{
					Id:  update.UUID,
					Url: update.Message.Text,
				},
			}
		} else {
			s.outputMessageChannel <- domain.MessageReq{
				UUID:        update.UUID,
				Destination: domain.VideoMessageSand,
				Message: domain.MessageSendReq{
					UserId: update.Message.From.ID,
					HashId: video.HashID,
				},
			}
		}
	}
	return true
}
func (s UrlProcessStrategy) CanProcess(update domain.Update) bool {
	return strings.Contains(update.Message.Text, "https")
}
