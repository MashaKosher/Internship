CREATE TABLE IF NOT EXISTS users (
  id        INTEGER  NOT NULL PRIMARY KEY,
  login     TEXT    NOT NULL,
  balance   DECIMAL(10, 2),
  win_rate  DECIMAL(3, 2)
);


CREATE TABLE IF NOT EXISTS daily_tasks (
  task_date        DATE PRIMARY KEY,
  referals_amount    INTEGER,
  wins_amount        INTEGER
);



CREATE TABLE IF NOT EXISTS seasons (
  id              BIGSERIAL               NOT NULL PRIMARY KEY,
  season_start    TIMESTAMP WITH TIME ZONE NOT NULL,
  season_end      TIMESTAMP WITH TIME ZONE NOT NULL,
  season_fund     INTEGER,
  season_status   TEXT
);

CREATE TABLE IF NOT EXISTS leaderboard (
  season_id BIGINT NOT NULL,
  user_id INTEGER NOT NULL,
  win INTEGER NOT NULL DEFAULT 0,
  PRIMARY KEY (season_id, user_id),
  FOREIGN KEY (season_id) REFERENCES seasons (id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);