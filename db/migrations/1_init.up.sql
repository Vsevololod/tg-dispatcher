create table if not exists "user"
(
    id         bigint       not null
        primary key,
    first_name varchar(255) not null,
    username   varchar(255) not null
);

create table if not exists video
(
    hash_id        varchar(255) not null
        primary key,
    original_id    bigint       not null,
    url            text         not null,
    video_id       text         not null,
    load_timestamp bigint       not null,
    path           text         not null,
    title          text         not null,
    duration       bigint       not null,
    timestamp      bigint       not null,
    filesize       bigint,
    thumbnail      varchar(255) not null,
    channel_url    varchar(255) not null,
    channel_id     varchar(255) not null,
    user_id        bigint       not null
        references "user",
    channel        varchar(255) default 'none'::character varying,
    loaded_times   bigint       default 0
);

create unique index video_video_id
    on video (video_id);

create index video_user_id
    on video (user_id);