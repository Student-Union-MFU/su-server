CREATE TABLE steps (
    id            SERIAL PRIMARY KEY,
    user_id       INT REFERENCES users(id) ON DELETE CASCADE,
    step_count    INT NOT NULL,
    recorded_date DATE NOT NULL,
    synced_at     TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, recorded_date)
);

