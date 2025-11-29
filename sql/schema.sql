-- Prices table
CREATE TABLE IF NOT EXISTS prices (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    region TEXT NOT NULL,
    price REAL NOT NULL,
    unit TEXT,
    source TEXT,
    recorded_at TEXT NOT NULL,
    created_at TEXT DEFAULT (datetime('now'))
);

-- Weather history table
CREATE TABLE IF NOT EXISTS weather_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    region TEXT NOT NULL,
    temp_c REAL,
    humidity INTEGER,
    rain_mm REAL,
    fetched_at TEXT NOT NULL,
    created_at TEXT DEFAULT (datetime('now'))
);