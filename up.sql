CREATE TABLE IF NOT EXISTS api_tokens (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name varchar(255) NOT NULL, -- glearning
    token text UNIQUE NOT NULL, -- some_secret_token
    active tinyint(1) NOT NULL DEFAULT '1'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO api_tokens (name, token, active) VALUES ('glearning', '6qPRWChjOmfziYo0dASFKS+vnkZGxHgg', 1);

drop table if exists sync_logs_batches;
drop table if exists sync_logs;

CREATE TABLE if not exists sync_logs (
  id CHAR(26) NOT NULL PRIMARY KEY,
  id_instansi CHAR(26) DEFAULT NULL,
  entity ENUM('mahasiswa', 'dosen', 'kelas', 'perkuliahan') NOT NULL,
  status ENUM('pending', 'processing', 'completed', 'failed') NOT NULL,
  total_data INT DEFAULT 0,
  start_time DATETIME NOT NULL,
  end_time DATETIME,
  error_message TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_entity_status (entity, status)
);

CREATE TABLE if not exists  sync_logs_batches (
   id CHAR(26) NOT NULL PRIMARY KEY,
   sync_log_id CHAR(26) NOT NULL,
   batch_sequence INT NOT NULL,
   total_data INT DEFAULT 0,
   success_count INT DEFAULT 0,
   failed_count INT DEFAULT 0,
   error_message TEXT,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   FOREIGN KEY (sync_log_id) REFERENCES sync_logs(id) ON DELETE CASCADE
);
