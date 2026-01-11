-- Создание временной таблицы с информацией о старте авторизации
CREATE TABLE IF NOT EXISTS temp_auth_init (
    id VARCHAR(50) PRIMARY KEY,
    person_id VARCHAR(36),
    response_type VARCHAR NOT NULL,
    client_id VARCHAR NOT NULL,
    redirect_uri VARCHAR NOT NULL,
    scope VARCHAR NOT NULL,
    state VARCHAR NOT NULL,
    code_challenge_method VARCHAR NOT NULL,
    code_challenge VARCHAR NOT NULL,
    create_time_ms bigint NOT NULL,
    last_modify_time_ms bigint NOT NULL
);

COMMENT ON TABLE temp_auth_init IS 'Временное хранение цепочек для авторизации';

COMMENT ON COLUMN temp_auth_init.id IS 'Уникальный код авторизации';
COMMENT ON COLUMN temp_auth_init.person_id IS 'Id пользователя из таблицы person';
COMMENT ON COLUMN temp_auth_init.response_type IS 'Тип ожидаемого ответа';
COMMENT ON COLUMN temp_auth_init.client_id IS 'Идентификатор клиентского приложения';
COMMENT ON COLUMN temp_auth_init.redirect_uri IS 'URI, куда будет направлен код авторизации после успешной аутентификации пользователя';
COMMENT ON COLUMN temp_auth_init.scope IS 'Область доступа (scope). Определяет, к каким ресурсам приложение запрашивает доступ';
COMMENT ON COLUMN temp_auth_init.state IS 'Параметр для предотвращения атаки CSRF';
COMMENT ON COLUMN temp_auth_init.code_challenge_method IS 'Способ шифрования, которым был зашифрован проверочный код (PKCE)';
COMMENT ON COLUMN temp_auth_init.code_challenge IS 'Зашифрованный проверочный код (PKCE)';
COMMENT ON COLUMN temp_auth_init.create_time_ms IS 'Unix timestamp создания записи в ms';
COMMENT ON COLUMN temp_auth_init.last_modify_time_ms IS 'Unix timestamp последнего изменения записи в ms';