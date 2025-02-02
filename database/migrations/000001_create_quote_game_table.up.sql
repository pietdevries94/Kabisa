CREATE TABLE IF NOT EXISTS quote_game(
   id BLOB PRIMARY KEY,
   quote1_id INT NOT NULL,
   quote2_id INT NOT NULL,
   quote3_id INT NOT NULL,
   quote1_correct BOOLEAN NULL,
   quote2_correct BOOLEAN NULL,
   quote3_correct BOOLEAN NULL,
   created_at DATETIME NOT NULL,
   completed_at DATETIME NULL
);