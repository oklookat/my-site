---------------------- SERVICE FUNCTIONS START ----------------------

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
CREATE DOMAIN ulid AS varchar(26) DEFAULT generate_ulid()
CONSTRAINT ulid_length_check CHECK (char_length(value) = 26)
CONSTRAINT ulid_upper_bound CHECK (value <= '7ZZZZZZZZZZZZZZZZZZZZZZZZZ');

-- custom types
DROP DOMAIN IF EXISTS created CASCADE;
CREATE DOMAIN created AS timestamp with time zone DEFAULT current_timestamp NOT NULL;
---------------------- SERVICE FUNCTIONS END ----------------------


---------------------- BASIC FUNCTIONS. ACTUAL FOR ALL MODELS. ----------------------
-------- before insert or update --------
CREATE OR REPLACE FUNCTION before_insert_or_update() RETURNS TRIGGER AS
$$
BEGIN
    -- set updated date
    NEW.updated_at
        = now();
    RETURN NEW;
END;
$$
language 'plpgsql';


-------- before insert --------
CREATE
    OR REPLACE FUNCTION before_insert()
    RETURNS TRIGGER AS
$$
BEGIN
    -- set created date
    NEW.created_at
        = now();
    RETURN NEW;
END;
$$
language 'plpgsql';


-------- before update --------
CREATE OR REPLACE FUNCTION before_update()
RETURNS TRIGGER AS
$$
BEGIN
    IF
        -- prevent id changing
        NEW.id != OLD.id THEN
        NEW.id = OLD.id;
    END IF;
    IF
        -- prevent created date changing
        NEW.created_at != OLD.created_at THEN
        NEW.created_at = OLD.created_at;
    END IF;
    RETURN NEW;
END;
$$
language 'plpgsql';
---------------------- BASIC FUNCTIONS END ----------------------


---------------------- TABLES START ----------------------
---- USERS ----
CREATE EXTENSION IF NOT EXISTS pgcrypto;
ALTER TABLE IF EXISTS users DROP CONSTRAINT IF EXISTS user_length CASCADE;
DROP TYPE IF EXISTS user_roles CASCADE;
CREATE TYPE user_roles AS ENUM ('user', 'admin');
CREATE TABLE users
(
    id         ulid PRIMARY KEY,
    role       user_roles DEFAULT 'user'::user_roles NOT NULL,
    username   varchar(24) UNIQUE NOT NULL CHECK (length(username) > 0),
    password   varchar(256) NOT NULL CHECK (length(password) > 8),
    reg_ip     varchar(64) DEFAULT NULL,
    reg_agent  varchar(324) DEFAULT NULL,
    created_at created,
    updated_at timestamp with time zone,
    CONSTRAINT users_length CHECK (length(username) >= 4 AND (length(password) >= 8))
) TABLESPACE pg_default;
ALTER TABLE users OWNER to postgres;
CREATE TRIGGER user_before_insert_or_update BEFORE INSERT OR UPDATE ON users FOR EACH ROW
EXECUTE PROCEDURE before_insert_or_update();
CREATE TRIGGER user_before_insert BEFORE INSERT ON users FOR EACH ROW
EXECUTE PROCEDURE before_insert();
CREATE TRIGGER user_before_update BEFORE UPDATE ON users FOR EACH ROW
EXECUTE PROCEDURE before_update();


---- TOKENS ----
CREATE TABLE tokens
(
    id ulid PRIMARY KEY,
    user_id ulid NOT NULL
    REFERENCES users(id)
    ON UPDATE CASCADE
    ON DELETE CASCADE,
    token varchar(2048) NOT NULL UNIQUE,
    last_ip varchar(64) DEFAULT NULL,
    last_agent varchar(324) DEFAULT NULL,
    auth_ip varchar(64) DEFAULT NULL,
    auth_agent varchar(324) DEFAULT NULL,
    created_at created,
    updated_at timestamp with time zone
) TABLESPACE pg_default;
ALTER TABLE tokens OWNER to postgres;
CREATE TRIGGER token_before_insert_or_update BEFORE INSERT OR UPDATE ON tokens FOR EACH ROW
EXECUTE PROCEDURE before_insert_or_update();
CREATE TRIGGER token_before_insert BEFORE INSERT ON tokens FOR EACH ROW
EXECUTE PROCEDURE before_insert();
CREATE TRIGGER token_before_update BEFORE UPDATE ON tokens FOR EACH ROW
EXECUTE PROCEDURE before_update();


---- FILES ----
CREATE TABLE files
(
    id ulid PRIMARY KEY,
    user_id ulid NOT NULL
    REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    hash varchar(32) NOT NULL UNIQUE CHECK (length(hash) = 32),
    path varchar(512) NOT NULL UNIQUE,
    name varchar(104) NOT NULL UNIQUE,
    original_name varchar(264) DEFAULT NULL,
    extension varchar(64) DEFAULT NULL,
    size bigint NOT NULL,
    created_at created,
    updated_at timestamp with time zone
) TABLESPACE pg_default;
ALTER TABLE files OWNER to postgres;
CREATE TRIGGER file_before_insert_or_update BEFORE INSERT OR UPDATE ON files FOR EACH ROW
EXECUTE PROCEDURE before_insert_or_update();
CREATE TRIGGER file_before_insert BEFORE INSERT ON files FOR EACH ROW
EXECUTE PROCEDURE before_insert();
CREATE TRIGGER file_before_update BEFORE UPDATE ON files FOR EACH ROW
EXECUTE PROCEDURE before_update();


---- ARTICLES ----
CREATE TABLE articles
(
    id ulid  PRIMARY KEY,
    user_id ulid NOT NULL
    REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    cover_id ulid DEFAULT NULL
    REFERENCES files(id) ON UPDATE CASCADE ON DELETE SET NULL,
    is_published boolean DEFAULT false,
    title varchar(124) DEFAULT 'Untitled'::varchar,
    content varchar(816000) DEFAULT ''::varchar NOT NULL,
    published_at timestamp with time zone,
    created_at created,
    updated_at timestamp with time zone
) TABLESPACE pg_default;
ALTER TABLE articles OWNER to postgres;
CREATE TRIGGER article_before_insert_or_update BEFORE INSERT OR UPDATE ON articles FOR EACH ROW
EXECUTE PROCEDURE before_insert_or_update();
CREATE TRIGGER article_before_insert BEFORE INSERT ON articles FOR EACH ROW
EXECUTE PROCEDURE before_insert();
CREATE TRIGGER article_before_update BEFORE UPDATE ON articles FOR EACH ROW
EXECUTE PROCEDURE before_update();