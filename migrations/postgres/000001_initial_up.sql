DO $$
    DECLARE
        candidate_id UUID;
    BEGIN
        -- EXTENSIONS --
        EXECUTE 'CREATE EXTENSION IF NOT EXISTS pgcrypto';

        -- TABLES --
        EXECUTE '
        CREATE TABLE IF NOT EXISTS candidates (
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            full_name VARCHAR NOT NULL,
            email VARCHAR NOT NULL,
            phone BIGINT NOT NULL
        )
    ';

        EXECUTE '
        CREATE TABLE IF NOT EXISTS recruiters (
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            candidate_id UUID NOT NULL REFERENCES candidates (id),
            full_name VARCHAR NOT NULL,
            email VARCHAR NOT NULL,
            phone BIGINT NOT NULL
        )
    ';

        -- DATA --
        INSERT INTO candidates (full_name, email, phone)
        VALUES ('Nurdaulet', 'zzz@mail.ru', 3333)
        RETURNING id INTO candidate_id;

        INSERT INTO recruiters (candidate_id, full_name, email, phone)
        VALUES (candidate_id, 'Nariman', 'sss@mail.ru', 223);

        -- COMMIT не требуется внутри DO блока, так как он выполняется как единственная транзакция
    END
$$ LANGUAGE plpgsql;
