-- Стартовое наполнение данных
INSERT INTO gloss_auth_client (
    id, scope, page_path_login, active, create_time_ms, last_modify_time_ms
)
VALUES (
    'myBank',
    'dbo',
    '/login',
    true,
    EXTRACT(epoch FROM CURRENT_TIMESTAMP) * 1000,
    EXTRACT(epoch FROM CURRENT_TIMESTAMP) * 1000
);