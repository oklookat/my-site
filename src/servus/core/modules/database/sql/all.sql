-- users
CREATE
EXTENSION IF NOT EXISTS pgcrypto;
ALTER TABLE IF EXISTS users DROP CONSTRAINT IF EXISTS user_length CASCADE;
DROP TYPE IF EXISTS user_roles CASCADE;
CREATE TYPE user_roles AS ENUM ('user', 'admin');

CREATE TABLE public.users
(
    id         uuid PRIMARY KEY                                   default gen_random_uuid() NOT NULL,
    role       user_roles                                         DEFAULT 'user'::user_roles NOT NULL,
    username   character varying(24) COLLATE pg_catalog."default"                           NOT NULL,
    password   character varying(104) COLLATE pg_catalog."default"                          NOT NULL,
    reg_ip     character varying(54) COLLATE pg_catalog."default" DEFAULT 'unknown':: character varying,
    reg_agent  varchar COLLATE pg_catalog."default"               DEFAULT 'unknown':: varchar,
    created_at timestamp with time zone                           default current_timestamp NOT NULL,
    updated_at timestamp with time zone,
    CONSTRAINT users_username_key UNIQUE (username),
    CONSTRAINT users_length CHECK (length(username) >= 4 AND (length(password) >= 8))
) TABLESPACE pg_default;
ALTER TABLE public.users
    OWNER to postgres;


-- basic functions. Actual for all models.
CREATE
OR REPLACE FUNCTION before_insert_or_update()
    RETURNS TRIGGER AS $$
BEGIN
      NEW.updated_at
= now();
RETURN NEW;
END;
    $$
language 'plpgsql';


    CREATE
OR REPLACE FUNCTION before_insert()
    RETURNS TRIGGER AS $$
BEGIN
      NEW.created_at
= now();
RETURN NEW;
END;
    $$
language 'plpgsql';


    CREATE
OR REPLACE FUNCTION before_update()
    RETURNS TRIGGER AS $$
BEGIN
      IF
NEW.id != OLD.id THEN
        NEW.id = OLD.id;
END IF;
      IF
NEW.created_at != OLD.created_at THEN
        NEW.created_at = OLD.created_at;
END IF;
RETURN NEW;
END;
    $$
language 'plpgsql';

CREATE TRIGGER user_before_insert_or_update
    BEFORE INSERT OR
UPDATE ON users FOR EACH ROW
    EXECUTE PROCEDURE before_insert_or_update();
CREATE TRIGGER user_before_insert
    BEFORE INSERT
    ON users
    FOR EACH ROW EXECUTE PROCEDURE before_insert();
CREATE TRIGGER user_before_update
    BEFORE UPDATE
    ON users
    FOR EACH ROW EXECUTE PROCEDURE before_update();


-- tokens
CREATE TABLE public.tokens
(
    id         uuid PRIMARY KEY                                   default gen_random_uuid() NOT NULL,
    user_id    uuid                                                                         NOT NULL,
    token      character varying(4024) COLLATE pg_catalog."default"                         NOT NULL,
    last_ip    character varying(32) COLLATE pg_catalog."default" DEFAULT 'unknown':: character varying,
    last_agent varchar COLLATE pg_catalog."default"               DEFAULT 'unknown':: varchar,
    auth_ip    character varying(32) COLLATE pg_catalog."default" DEFAULT 'unknown':: character varying,
    auth_agent varchar COLLATE pg_catalog."default"               DEFAULT 'unknown':: varchar,
    created_at timestamp with time zone                           default current_timestamp,
    updated_at timestamp with time zone,
    CONSTRAINT fk_users_tokens FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
) TABLESPACE pg_default;
ALTER TABLE public.tokens
    OWNER to postgres;


CREATE TRIGGER token_before_insert_or_update
    BEFORE INSERT OR
UPDATE ON tokens FOR EACH ROW
    EXECUTE PROCEDURE before_insert_or_update();
CREATE TRIGGER token_before_insert
    BEFORE INSERT
    ON tokens
    FOR EACH ROW EXECUTE PROCEDURE before_insert();
CREATE TRIGGER token_before_update
    BEFORE UPDATE
    ON tokens
    FOR EACH ROW EXECUTE PROCEDURE before_update();


-- articles
CREATE TABLE public.articles
(
    id           uuid PRIMARY KEY                                    default gen_random_uuid() NOT NULL,
    user_id      uuid                                                                          NOT NULL,
    is_published boolean                                             DEFAULT false,
    title        character varying(124) COLLATE pg_catalog."default" DEFAULT 'Untitled':: character varying,
    content      text COLLATE pg_catalog."default"                                             NOT NULL,
    slug         character varying(256) COLLATE pg_catalog."default",
    published_at timestamp with time zone,
    created_at   timestamp with time zone                            default current_timestamp,
    updated_at   timestamp with time zone,
    CONSTRAINT articles_slug_key UNIQUE (slug),
    CONSTRAINT fk_users_articles FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
) TABLESPACE pg_default;

ALTER TABLE public.articles
    OWNER to postgres;


CREATE TRIGGER article_before_insert_or_update
    BEFORE INSERT OR
UPDATE ON articles FOR EACH ROW
    EXECUTE PROCEDURE before_insert_or_update();
CREATE TRIGGER article_before_insert
    BEFORE INSERT
    ON articles
    FOR EACH ROW EXECUTE PROCEDURE before_insert();
CREATE TRIGGER article_before_update
    BEFORE UPDATE
    ON articles
    FOR EACH ROW EXECUTE PROCEDURE before_update();


-- files
CREATE TABLE public.files
(
    id            uuid PRIMARY KEY                           default gen_random_uuid() NOT NULL,
    user_id       uuid                              NOT NULL,
    hash          text COLLATE pg_catalog."default" NOT NULL,
    path          text COLLATE pg_catalog."default" NOT NULL,
    name          text COLLATE pg_catalog."default" NOT NULL,
    original_name text COLLATE pg_catalog."default" NOT NULL DEFAULT 'unknown'::text,
    extension     text COLLATE pg_catalog."default" NOT NULL DEFAULT 'unknown'::text,
    size          text COLLATE pg_catalog."default" NOT NULL,
    created_at    timestamp with time zone                   default current_timestamp,
    updated_at    timestamp with time zone,
    CONSTRAINT fk_users_files FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
) TABLESPACE pg_default;
ALTER TABLE public.files
    OWNER to postgres;

CREATE TRIGGER file_before_insert_or_update
    BEFORE INSERT OR
UPDATE ON files FOR EACH ROW
    EXECUTE PROCEDURE before_insert_or_update();
CREATE TRIGGER file_before_insert
    BEFORE INSERT
    ON files
    FOR EACH ROW EXECUTE PROCEDURE before_insert();
CREATE TRIGGER file_before_update
    BEFORE UPDATE
    ON files
    FOR EACH ROW EXECUTE PROCEDURE before_update();