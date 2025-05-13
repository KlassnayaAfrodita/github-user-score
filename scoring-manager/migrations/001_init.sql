CREATE TABLE IF NOT EXISTS scoring_status (
                                              application_id BIGSERIAL PRIMARY KEY,
                                              user_id INT NOT NULL,
                                              status INTEGER NOT NULL CHECK (status IN (0, 1, 2)),
                                              created_at TIMESTAMP DEFAULT now()
    );

CREATE TABLE IF NOT EXISTS scoring_result (
                                              application_id BIGINT PRIMARY KEY,
                                              user_id INT NOT NULL,
                                              score INT NOT NULL,
                                              FOREIGN KEY (application_id) REFERENCES scoring_status(application_id) ON DELETE CASCADE
    );
