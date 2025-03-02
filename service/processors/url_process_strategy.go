package processors

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"log/slog"
	"strings"
	"tg-dispatcher/domain"
	"tg-dispatcher/domain/models"
	"tg-dispatcher/lib"
	"tg-dispatcher/lib/logger/sl"
	"tg-dispatcher/storage"

	"github.com/Vsevololod/tg-api-contracts-lib/gen/go/messages"
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

	tracer := otel.Tracer("tg-dispatcher")
	ctx, span := tracer.Start(update.Context, "ProcessMessage")
	defer span.End()

	videoId := lib.GetVideoIdFromUrl(update.Message.Text)
	span.SetAttributes(attribute.String("video-id", videoId))
	video, err := s.videoProvider.GetVideoById(ctx, videoId)
	if err != nil {
		if errors.Is(err, storage.ErrVideoNotFound) {
			err := s.videoSaver.SaveVideoMin(ctx,
				update.UUID,
				update.UpdateID,
				update.Message.Text,
				videoId,
				update.Message.From.ID,
			)
			if err != nil {
				span.RecordError(err)
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
				Context: ctx,
			}
		} else {
			s.log.Error("Cannot get video", sl.Err(err), sl.Req(update))
		}
	}
	message := messagesv1.TgSendMessage{
		Text:   video.Title,
		UserId: uint64(video.UserID),
		Type:   messagesv1.MessageType_IMAGE,
		Params: map[string]string{
			messagesv1.MessageParams_FILE_URL.String():  video.Path,
			messagesv1.MessageParams_PHOTO_URL.String(): video.Thumbnail,
		},
	}

	s.outputMessageChannel <- domain.MessageReq{
		UUID:        update.UUID,
		Destination: domain.VideoMessageSand,
		Message:     &message,
		Context:     ctx,
	}

	return true
}
func (s UrlProcessStrategy) CanProcess(update domain.Update) bool {
	return strings.Contains(update.Message.Text, "https")
}
