create table video
(
    hash_id        varchar(255)                               not null
        primary key,
    original_id    bigint                                     not null,
    url            text                                       not null,
    video_id       text                                       not null,
    load_timestamp bigint                                     not null,
    path           text         default ''::text              not null,
    title          text         default ''::text              not null,
    duration       bigint       default 0                     not null,
    timestamp      bigint       default 0                     not null,
    filesize       bigint       default 0,
    thumbnail      varchar(255) default ''::character varying not null,
    channel_url    varchar(255) default ''::character varying not null,
    channel_id     varchar(255) default ''::character varying not null,
    user_id        bigint                                     not null
        references "user",
    channel        varchar(255) default 'none'::character varying,
    loaded_times   bigint       default 0
);