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
    room_name   TEXT NOT NULL,
    score       NUMERIC(10,2) NOT NULL,   -- 小数も扱える
    user_id     INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES app_user(id)
);

-- 画像テーブル
CREATE TABLE img (
    id          SERIAL PRIMARY KEY,
    score       NUMERIC(10,2) NOT NULL,
    note        TEXT,                     -- 物体検出結果などのメモ
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

-- ============================
-- app_user（ユーザー）
-- ============================
INSERT INTO app_user (password) VALUES
('pass_user1'),
('pass_user2'),
('pass_user3');

-- ============================
-- room（部屋）
-- ============================
INSERT INTO room (room_name, score, user_id) VALUES
('リビング', 85.50, 1),
('キッチン', 72.30, 1),
('寝室', 90.00, 2),
('書斎', 65.75, 3);

-- ============================
-- detail（画像ごとの物体数）
-- ============================
INSERT INTO detail (img_id, obj_id, obj_cnt) VALUES
(1, 1, 3),
(1, 2, 1),
(2, 1, 5),
(2, 3, 2),
(3, 2, 4),
(4, 1, 1),
(4, 4, 6),
(5, 3, 2),
(5, 5, 1);
