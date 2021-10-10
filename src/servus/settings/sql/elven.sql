---------------------- SERVICE FUNCTIONS START
-- pgulid is based on OK Log's Go implementation of the ULID spec
--
-- https://github.com/oklog/ulid
-- https://github.com/ulid/spec
--
-- Copyright 2016 The Oklog Authors
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
-- http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE FUNCTION generate_ulid()
    RETURNS TEXT
AS
$$
DECLARE
-- Crockford's Base32
    encoding  BYTEA = '0123456789ABCDEFGHJKMNPQRSTVWXYZ';
    timestamp BYTEA = E'\\000\\000\\000\\000\\000\\000';
    output    TEXT  = '';
    unix_time BIGINT;
    ulid      BYTEA;
BEGIN
    -- 6 timestamp bytes
    unix_time = (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT;
    timestamp = SET_BYTE(timestamp, 0, (unix_time >> 40)::BIT(8)::INTEGER);
    timestamp = SET_BYTE(timestamp, 1, (unix_time >> 32)::BIT(8)::INTEGER);
    timestamp = SET_BYTE(timestamp, 2, (unix_time >> 24)::BIT(8)::INTEGER);
    timestamp = SET_BYTE(timestamp, 3, (unix_time >> 16)::BIT(8)::INTEGER);
    timestamp = SET_BYTE(timestamp, 4, (unix_time >> 8)::BIT(8)::INTEGER);
    timestamp = SET_BYTE(timestamp, 5, unix_time::BIT(8)::INTEGER);

    -- 10 entropy bytes
    ulid = timestamp || gen_random_bytes(10);

    -- Encode the timestamp
    output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 0) & 224) >> 5));
    output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 0) & 31)));
    output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 1) & 248) >> 3));
    output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 1) & 7) << 2) | ((GET_BYTE(ulid, 2) & 192) >> 6)));
    output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 2) & 62) >> 1));
    output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 2) & 1) << 4) | ((GET_BYTE(ulid, 3) & 240) >> 4)));
    output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 3) & 15) << 1) | ((GET_BYTE(ulid, 4) & 128) >> 7)));
    output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 4) & 124) >> 2));
    output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 4) & 3) << 3) | ((GET_BYTE(ulid, 5) & 224) >> 5)));
    output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 5) & 31)));

    -- Encode the entropy
    output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 6) & 248) >> 3));
    output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 6) & 7) << 2) | ((GET_BYTE(ulid, 7) & 192) >> 6)));
    output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 7) & 62) >> 1));
    output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 7) & 1) << 4) | ((GET_BYTE(ulid, 8) & 240) >> 4)));
    output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 8) & 15) << 1) | ((GET_BYTE(ulid, 9) & 128) >> 7)));
    output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 9) & 124) >> 2));
    output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 9) & 3) << 3) | ((GET_BYTE(ulid, 10) & 224) >> 5)));
    output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 10) & 31)));
    output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 11) & 248) >> 3));
    output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 11) & 7) << 2) | ((GET_BYTE(ulid, 12) & 192) >> 6)));
    output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 12) & 62) >> 1));
    output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 12) & 1) << 4) | ((GET_BYTE(ulid, 13) & 240) >> 4)));
    output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 13) & 15) << 1) | ((GET_BYTE(ulid, 14) & 128) >> 7)));
    output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 14) & 124) >> 2));
    output = output || CHR(GET_BYTE(encoding, ((GET_BYTE(ulid, 14) & 3) << 3) | ((GET_BYTE(ulid, 15) & 224) >> 5)));
    output = output || CHR(GET_BYTE(encoding, (GET_BYTE(ulid, 15) & 31)));

    RETURN output;
END
$$
    LANGUAGE plpgsql
    VOLATILE;

-- create ulid type from https://github.com/geckoboard/pgulid/issues/4
DROP DOMAIN IF EXISTS uild CASCADE;
create domain ulid as varchar(26)
    default generate_ulid()
    not null
    constraint ulid_length_check check (char_length(value) = 26)
    constraint ulid_upper_bound check (value <= '7ZZZZZZZZZZZZZZZZZZZZZZZZZ');
---------------------- SERVICE FUNCTIONS END


---------------------- TABLES START
-- users
CREATE
    EXTENSION IF NOT EXISTS pgcrypto;
ALTER TABLE IF EXISTS users
    DROP CONSTRAINT IF EXISTS user_length CASCADE;
DROP TYPE IF EXISTS user_roles CASCADE;
CREATE TYPE user_roles AS ENUM ('user', 'admin');

CREATE TABLE public.users
(
    id         ulid PRIMARY KEY,
    role       user_roles                                DEFAULT 'user'::user_roles NOT NULL,
    username   varchar(24) COLLATE pg_catalog."default"                             NOT NULL,
    password   varchar(256) COLLATE pg_catalog."default"                            NOT NULL,
    reg_ip     varchar(64) COLLATE pg_catalog."default"  DEFAULT NULL,
    reg_agent  varchar(324) COLLATE pg_catalog."default" DEFAULT NULL,
    created_at timestamp with time zone                  DEFAULT current_timestamp  NOT NULL,
    updated_at timestamp with time zone,
    CONSTRAINT users_username_key UNIQUE (username),
    CONSTRAINT users_length CHECK (length(username) >= 4 AND (length(password) >= 8))
) TABLESPACE pg_default;
ALTER TABLE public.users
    OWNER to postgres;

-- basic functions. Actual for all models.
CREATE
    OR REPLACE FUNCTION before_insert_or_update()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at
        = now();
    RETURN NEW;
END;
$$
    language 'plpgsql';


CREATE
    OR REPLACE FUNCTION before_insert()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.created_at
        = now();
    RETURN NEW;
END;
$$
    language 'plpgsql';


CREATE
    OR REPLACE FUNCTION before_update()
    RETURNS TRIGGER AS
$$
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
        UPDATE
    ON users
    FOR EACH ROW
EXECUTE PROCEDURE before_insert_or_update();
CREATE TRIGGER user_before_insert
    BEFORE INSERT
    ON users
    FOR EACH ROW
EXECUTE PROCEDURE before_insert();
CREATE TRIGGER user_before_update
    BEFORE UPDATE
    ON users
    FOR EACH ROW
EXECUTE PROCEDURE before_update();


-- tokens
CREATE TABLE public.tokens
(
    id         ulid PRIMARY KEY,
    user_id    ulid                                       NOT NULL,
    token      varchar(2048) COLLATE pg_catalog."default" NOT NULL,
    last_ip    varchar(64) COLLATE pg_catalog."default"  DEFAULT NULL,
    last_agent varchar(324) COLLATE pg_catalog."default" DEFAULT NULL,
    auth_ip    varchar(64) COLLATE pg_catalog."default"  DEFAULT NULL,
    auth_agent varchar(324) COLLATE pg_catalog."default" DEFAULT NULL,
    created_at timestamp with time zone                  DEFAULT current_timestamp,
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
        UPDATE
    ON tokens
    FOR EACH ROW
EXECUTE PROCEDURE before_insert_or_update();
CREATE TRIGGER token_before_insert
    BEFORE INSERT
    ON tokens
    FOR EACH ROW
EXECUTE PROCEDURE before_insert();
CREATE TRIGGER token_before_update
    BEFORE UPDATE
    ON tokens
    FOR EACH ROW
EXECUTE PROCEDURE before_update();


-- articles
CREATE TABLE public.articles
(
    id           ulid PRIMARY KEY,
    user_id      ulid  NOT NULL,
    is_published boolean                                   DEFAULT false,
    title        varchar(124) COLLATE pg_catalog."default" DEFAULT 'Untitled':: varchar,
    content      jsonb NOT NULL,
    slug         varchar(256) COLLATE pg_catalog."default",
    published_at timestamp with time zone,
    created_at   timestamp with time zone                  DEFAULT current_timestamp,
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
        UPDATE
    ON articles
    FOR EACH ROW
EXECUTE PROCEDURE before_insert_or_update();
CREATE TRIGGER article_before_insert
    BEFORE INSERT
    ON articles
    FOR EACH ROW
EXECUTE PROCEDURE before_insert();
CREATE TRIGGER article_before_update
    BEFORE UPDATE
    ON articles
    FOR EACH ROW
EXECUTE PROCEDURE before_update();


-- files
CREATE TABLE public.files
(
    id            ulid PRIMARY KEY,
    user_id       ulid                                      NOT NULL,
    hash          varchar(256) COLLATE pg_catalog."default"  NOT NULL,
    path          varchar(512) COLLATE pg_catalog."default" NOT NULL,
    name          varchar(512) COLLATE pg_catalog."default" NOT NULL,
    original_name varchar(512) COLLATE pg_catalog."default" DEFAULT NULL,
    extension     varchar(64) COLLATE pg_catalog."default"  DEFAULT NULL,
    size          bigint                                    NOT NULL,
    created_at    timestamp with time zone                  DEFAULT current_timestamp,
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
        UPDATE
    ON files
    FOR EACH ROW
EXECUTE PROCEDURE before_insert_or_update();
CREATE TRIGGER file_before_insert
    BEFORE INSERT
    ON files
    FOR EACH ROW
EXECUTE PROCEDURE before_insert();
CREATE TRIGGER file_before_update
    BEFORE UPDATE
    ON files
    FOR EACH ROW
EXECUTE PROCEDURE before_update();
---------------------- TABLES END