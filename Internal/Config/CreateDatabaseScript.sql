-- Users table (core user information)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Stats definition table (flexible stat types)
CREATE TABLE stat_types (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,  -- 'strength', 'perception', etc.
    display_name TEXT NOT NULL, -- 'Strength', 'Perception'
    min_value INTEGER DEFAULT 1,
    max_value INTEGER DEFAULT 10,
    default_value INTEGER DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- User stats (dynamic attributes)
CREATE TABLE user_stats (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stat_type_id INTEGER NOT NULL REFERENCES stat_types(id) ON DELETE CASCADE,
    value INTEGER NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_id, stat_type_id)
);

-- User results (win/lose/draw records)
CREATE TABLE user_results (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    total_wins INTEGER NOT NULL DEFAULT 0,
    total_draws INTEGER NOT NULL DEFAULT 0,
    total_lose INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_id)
);

-- Indexes for performance
CREATE INDEX idx_user_stats_user ON user_stats(user_id);
CREATE INDEX idx_user_stats_type ON user_stats(stat_type_id);