package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"tg-dispatcher/storage/postgresql/gen"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"tg-dispatcher/domain/models" // Импорт моделей
	"tg-dispatcher/storage"
)

type Storage struct {
	queries *gen.Queries
	db      *pgxpool.Pool
}

// New создает новое подключение к PostgreSQL и инициализирует sqlc
func New(dsn string) (*Storage, error) {
	const op = "storage.postgresql.New"

	// Подключаемся к PostgreSQL через pgxpool
	conn, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Создаем объект Queries из sqlc
	queries := gen.New(conn)

	return &Storage{db: conn, queries: queries}, nil
}

// SaveVideoMin сохраняет видео в БД
func (s *Storage) SaveVideoMin(ctx context.Context, hashId string, originalId int64, url string, videoId string, userId int64) error {
	const op = "storage.postgresql.SaveVideo"

	// Вставляем данные в таблицу video с помощью sqlc
	err := s.queries.CreateVideoMin(ctx, gen.CreateVideoMinParams{
		HashID:        hashId,
		OriginalID:    originalId,
		Url:           url,
		VideoID:       videoId,
		LoadTimestamp: getUnixTimestamp(),
		UserID:        userId,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// sqlc не поддерживает LastInsertId, так как PostgreSQL использует RETURNING
	return nil
}

// GetVideoById получает видео из БД по video_id
func (s *Storage) GetVideoById(ctx context.Context, videoId string) (models.Video, error) {
	const op = "storage.postgresql.Video"

	// Получаем видео из БД с помощью sqlc
	videoDB, err := s.queries.GetVideoByID(ctx, videoId)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return models.Video{}, fmt.Errorf("%s: %w", op, storage.ErrVideoNotFound)
		}
		return models.Video{}, fmt.Errorf("%s: %w", op, err)
	}

	// Маппим структуру sqlc в `models.Video`
	video := ConvertSQLCVideoToModel(videoDB)
	return video, nil
}

// ConvertSQLCVideoToModel конвертирует структуру `db.Video` (от sqlc) в `models.Video`
func ConvertSQLCVideoToModel(v gen.Video) models.Video {
	return models.Video{
		HashID:        v.HashID,
		OriginalID:    v.OriginalID,
		URL:           v.Url,
		VideoID:       v.VideoID,
		LoadTimestamp: v.LoadTimestamp,
		Path:          v.Path,
		Title:         v.Title,
		Duration:      v.Duration,
		Timestamp:     v.Timestamp,
		Filesize:      getInt64orZero(v.Filesize),
		Thumbnail:     v.Thumbnail,
		ChannelURL:    v.ChannelUrl,
		ChannelID:     v.ChannelID,
		UserID:        v.UserID,
		Channel:       getStrOrEmpty(v.Channel),
		LoadedTimes:   getInt64orZero(v.LoadedTimes),
	}
}

// getUnixTimestamp возвращает текущее время в формате Unix Timestamp
func getUnixTimestamp() int64 {
	return time.Now().Unix()
}

func getStrOrEmpty(s pgtype.Text) string {
	if s.Valid {
		return s.String
	} else {
		return ""
	}
}

// parseInt64 конвертирует строку в int64 (если число)
func getInt64orZero(value pgtype.Int8) int64 {
	if value.Valid {
		return value.Int64
	} else {
		return 0
	}
}
