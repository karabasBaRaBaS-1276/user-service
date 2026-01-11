-- Видится критически важным в таблице temp_auth_init не допускать изменений person_id, если он установлен.

-- Добавим функцию проверки значения для триггера
CREATE OR REPLACE FUNCTION user_service.check_person_id() RETURNS TRIGGER
AS $$
    BEGIN
        IF OLD.person_id IS DISTINCT FROM NEW.person_id AND OLD.person_id IS NOT NULL THEN
            RAISE EXCEPTION 'Cannot update person_id';
        END IF;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION user_service.check_person_id() IS 'Триггерная функция для проверки обновления поля person_id';

DO $$
    BEGIN
        -- Проверяем существование триггера именно для конкретной таблицы
        IF NOT EXISTS (
            SELECT 1 
            FROM pg_trigger 
            WHERE tgrelid = 'user_service.temp_auth_init'::regclass 
            AND tgname = 'prevent_person_id_update_trigger'
        ) THEN
            CREATE TRIGGER prevent_person_id_update_trigger
            BEFORE UPDATE ON user_service.temp_auth_init
            FOR EACH ROW
            EXECUTE FUNCTION user_service.check_person_id();

            -- Комментарий вешаем ТОЛЬКО если триггер создан или существует
            EXECUTE 'COMMENT ON TRIGGER prevent_person_id_update_trigger ON user_service.temp_auth_init IS ''Триггер для запрета обновления поля person_id''';
        END IF;
    END 
$$;