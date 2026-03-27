-- 外部キー制約を有効化（PostgreSQL はデフォルトで ON）
SET client_encoding = 'UTF8';

-- ユーザーテーブル
CREATE TABLE app_user (
    id          SERIAL PRIMARY KEY,
    password    TEXT NOT NULL
);

-- 部屋テーブル
CREATE TABLE room (
    id          SERIAL PRIMARY KEY,
    score       NUMERIC(10,2) NOT NULL,   -- 小数も扱える
    user_id     INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES app_user(id)
);

-- 画像テーブル
CREATE TABLE img (
    id          SERIAL PRIMARY KEY,
    score       NUMERIC(10,2) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    room_id     INTEGER NOT NULL,
    FOREIGN KEY (room_id) REFERENCES room(id)
);

-- 画像ごとの物体数テーブル
CREATE TABLE detail (
    img_id      INTEGER NOT NULL,
    obj_id      INTEGER NOT NULL,
    obj_cnt     INTEGER NOT NULL,
    PRIMARY KEY (img_id, obj_id),
    FOREIGN KEY (img_id) REFERENCES img(id)
);
