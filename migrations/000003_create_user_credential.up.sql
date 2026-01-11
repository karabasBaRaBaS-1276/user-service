-- Создание таблицы учетных записей, если она не существует
CREATE TABLE IF NOT EXISTS user_credential (
    username VARCHAR(50) PRIMARY KEY,
    person_id VARCHAR(36) NOT NULL UNIQUE,
    salt VARCHAR(36) NOT NULL,
    alg VARCHAR(20) NOT NULL,
    formula VARCHAR(20) NOT NULL,
    pass_hash VARCHAR(512) NOT NULL,
    active boolean NOT NULL,
    create_time_ms bigint NOT NULL,
    last_modify_time_ms bigint NOT NULL
);

COMMENT ON TABLE user_credential IS 'Таблица для хранения учетных данных пользователя с целью их авторизации';

COMMENT ON COLUMN user_credential.username IS 'Уникальный идентификатор учетной записи';
COMMENT ON COLUMN user_credential.person_id IS 'Id пользователя из таблицы person';
COMMENT ON COLUMN user_credential.salt IS 'Соль для пароля (UUIDv4)';
COMMENT ON COLUMN user_credential.alg IS 'Алгоритм хеширования пароля';
COMMENT ON COLUMN user_credential.formula IS 'Id формулы для расчета строки хеширования';
COMMENT ON COLUMN user_credential.pass_hash IS 'Хеш пароля';
COMMENT ON COLUMN user_credential.active IS 'Если true, то учетная запись активна';
COMMENT ON COLUMN user_credential.create_time_ms IS 'Unix timestamp создания записи в ms';
COMMENT ON COLUMN user_credential.last_modify_time_ms IS 'Unix timestamp последнего изменения записи в ms';

-- Добавление индекса по дате создания, если он не существует.
CREATE INDEX IF NOT EXISTS idx_user_credential_create_time_ms ON user_credential (create_time_ms);