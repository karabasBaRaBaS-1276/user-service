-- Создание таблицы-справочника с информацией о приложениях, которые имеют право авторизоваться
CREATE TABLE IF NOT EXISTS gloss_auth_client (
    id VARCHAR(50) PRIMARY KEY,
    scope VARCHAR NOT NULL,
    page_path_login VARCHAR(250),
    active boolean NOT NULL,
    create_time_ms bigint NOT NULL,
    last_modify_time_ms bigint NOT NULL
);

COMMENT ON TABLE gloss_auth_client IS 'Справочник приложений-клиентов, которые имеют право авторизоваться';

COMMENT ON COLUMN gloss_auth_client.id IS 'Уникальный идентификатор клиента';
COMMENT ON COLUMN gloss_auth_client.scope IS 'Список допустимых прав клиента. Значения разделены запятой';
COMMENT ON COLUMN gloss_auth_client.page_path_login IS 'Путь, по которому доступна аутентификация пользователя. Заполняется, если для авторизации клиента требуется аутентификация пользователя.';
COMMENT ON COLUMN gloss_auth_client.active IS 'Если true, то запись активна';
COMMENT ON COLUMN gloss_auth_client.create_time_ms IS 'Unix timestamp создания записи в ms';
COMMENT ON COLUMN gloss_auth_client.last_modify_time_ms IS 'Unix timestamp последнего изменения записи в ms';