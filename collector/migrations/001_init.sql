CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username TEXT NOT NULL UNIQUE
);

CREATE TABLE user_stats (
                            user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                            repos INTEGER NOT NULL,
                            stars INTEGER NOT NULL,
                            forks INTEGER NOT NULL,
                            commits INTEGER NOT NULL,
                            updated_at TIMESTAMP DEFAULT now(),
                            PRIMARY KEY (user_id)
);
