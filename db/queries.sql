-- name: CreateVideo :exec
INSERT INTO video (
    hash_id, original_id, url, video_id, load_timestamp, path, title, duration,
    timestamp, filesize, thumbnail, channel_url, channel_id, user_id, channel, loaded_times
) VALUES (
             @hash_id, @original_id, @url, @video_id, @load_timestamp, @path, @title, @duration,
             @timestamp, @filesize, @thumbnail, @channel_url, @channel_id, @user_id, @channel, @loaded_times
         );


-- name: CreateVideoMin :exec
INSERT INTO video(hash_id, original_id, url, video_id, load_timestamp, user_id)
VALUES (@hash_id, @original_id, @url, @video_id, @load_timestamp, @user_id);



-- name: GetVideoByID :one
SELECT * FROM video WHERE hash_id = @hash_id;

-- name: UpdateVideo :exec
UPDATE video SET
                 title = @title,
                 url = @url,
                 duration = @duration,
                 path = @path,
                 filesize = @filesize,
                 channel = @channel,
                 loaded_times = @loaded_times
WHERE hash_id = @hash_id;

