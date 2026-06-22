CREATE TABLE events (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    content     TEXT NOT NULL,
    location    VARCHAR(255),
    date        VARCHAR(50) NOT NULL,
    time        VARCHAR(50) NOT NULL,
    link        VARCHAR(255),
    created_at  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE event_images (
    id          SERIAL PRIMARY KEY,
    event_id    INT REFERENCES events(id) ON DELETE CASCADE,
    url         VARCHAR(255) NOT NULL,
    position    INT DEFAULT 0
);

CREATE TABLE lost_and_found (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    content     TEXT NOT NULL,
    date        VARCHAR(50) NOT NULL,
    time        VARCHAR(50) NOT NULL,
    link        VARCHAR(255),
    created_at  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE lost_and_found_images (
    id          SERIAL PRIMARY KEY,
    item_id     INT REFERENCES lost_and_found(id) ON DELETE CASCADE,
    url         VARCHAR(255) NOT NULL,
    position    INT DEFAULT 0
);
