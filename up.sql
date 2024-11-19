CREATE TABLE IF NOT EXISTS api_tokens (
    id bigint(20) unsigned NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name varchar(255) NOT NULL, -- glearning
    token text UNIQUE NOT NULL, -- some_secret_token
    active tinyint(1) NOT NULL DEFAULT '1'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO api_tokens (name, token, active) VALUES ('glearning', '6qPRWChjOmfziYo0dASFKS+vnkZGxHgg', 1);