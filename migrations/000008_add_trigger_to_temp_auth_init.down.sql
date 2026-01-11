DROP TRIGGER IF EXISTS prevent_person_id_update_trigger
    ON user_service.temp_auth_init;

DROP FUNCTION IF EXISTS user_service.check_person_id();