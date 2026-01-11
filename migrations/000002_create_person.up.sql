-- Добавление таблицы для хранения ПД пользователя
CREATE TABLE IF NOT EXISTS person (
    id VARCHAR(36) PRIMARY KEY,
    active BOOLEAN DEFAULT false,
    last_name VARCHAR(100) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    email VARCHAR(50) NOT NULL UNIQUE,
    create_time_ms bigint NOT NULL,
    last_modify_time_ms bigint NOT NULL
);

COMMENT ON TABLE person IS 'Таблица для хранения персональных данных пользователя';

COMMENT ON COLUMN person.id IS 'Уникальный id пользователя (UUIDv7)';
COMMENT ON COLUMN person.active IS 'Если true, то пользователь активен (не заблокирован)';
COMMENT ON COLUMN person.last_name IS 'Фамилия пользователя';
COMMENT ON COLUMN person.first_name IS 'Имя пользователя';
COMMENT ON COLUMN person.middle_name IS 'Отчество пользователя';
COMMENT ON COLUMN person.email IS 'Уникальный email пользователя';
COMMENT ON COLUMN person.create_time_ms IS 'Unix timestamp создания пользователя в ms';
COMMENT ON COLUMN person.last_modify_time_ms IS 'Unix timestamp последнего изменения пользователя в ms';

-- Добавление индекса по дате создания, если он не существует.
CREATE INDEX IF NOT EXISTS idx_person_create_time_ms ON person (create_time_ms);