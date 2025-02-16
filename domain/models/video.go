package models

// Video представляет структуру для хранения информации о видео.
type Video struct {
	HashID        string `db:"hash_id"`        // Первичный ключ
	OriginalID    int64  `db:"original_id"`    // Оригинальный ID видео
	URL           string `db:"url"`            // URL видео
	VideoID       string `db:"video_id"`       // ID видео
	LoadTimestamp int64  `db:"load_timestamp"` // Временная метка загрузки
	Path          string `db:"path"`           // Путь к файлу видео
	Title         string `db:"title"`          // Заголовок видео
	Duration      int64  `db:"duration"`       // Длительность видео
	Timestamp     int64  `db:"timestamp"`      // Временная метка (возможно, публикации)
	Filesize      int64  `db:"filesize"`       // Размер файла (может быть NULL)
	Thumbnail     string `db:"thumbnail"`      // Ссылка на миниатюру
	ChannelURL    string `db:"channel_url"`    // URL канала
	ChannelID     string `db:"channel_id"`     // ID канала
	UserID        int64  `db:"user_id"`        // ID пользователя (ссылка на "user")
	Channel       string `db:"channel"`        // Название канала (по умолчанию 'none')
	LoadedTimes   int64  `db:"loaded_times"`   // Количество загрузок (по умолчанию 0)
}
