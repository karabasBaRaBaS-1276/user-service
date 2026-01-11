-- Создание таблицы с публичными ключами для проверки подписи JWT токена
CREATE TABLE IF NOT EXISTS public_key (
    kid VARCHAR(36) PRIMARY KEY,
    json_public_jwk VARCHAR NOT NULL,
    expiration_time_ms BIGINT NOT NULL,
    create_time_ms BIGINT NOT NULL,
    last_modify_time_ms BIGINT NOT NULL
);

COMMENT ON TABLE public_key IS 'Публичные ключи для проверки подписи JWT токена';

COMMENT ON COLUMN public_key.kid IS 'Уникальный id ключа';
COMMENT ON COLUMN public_key.json_public_jwk IS 'JSON представление публичного JWK';
COMMENT ON COLUMN public_key.expiration_time_ms IS 'Время действия ключа в Unix timestamp ms';
COMMENT ON COLUMN public_key.create_time_ms IS 'Unix timestamp создания записи в ms';
COMMENT ON COLUMN public_key.last_modify_time_ms IS 'Unix timestamp последнего изменения записи в ms';