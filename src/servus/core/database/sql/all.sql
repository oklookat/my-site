-- Table: public.users

-- DROP TABLE public.users;

CREATE EXTENSION IF NOT EXISTS pgcrypto;
DROP TYPE IF EXISTS user_roles CASCADE;
CREATE TYPE user_roles AS ENUM ('user', 'admin');
CREATE TABLE public.users
(
    id uuid PRIMARY KEY default gen_random_uuid(),
    role user_roles DEFAULT 'user'::user_roles NOT NULL,
    username character varying(24) COLLATE pg_catalog."default" NOT NULL,
    password character varying(64) COLLATE pg_catalog."default" NOT NULL,
    reg_ip character varying(54) COLLATE pg_catalog."default" DEFAULT 'unknown'::character varying,
    reg_agent character varying(128) COLLATE pg_catalog."default" DEFAULT 'unknown'::character varying,
    created_at timestamp with time zone default current_timestamp NOT NULL,
    updated_at timestamp with time zone,
    CONSTRAINT users_username_key UNIQUE (username),
	CONSTRAINT username_min_length check (length(username) >= 4),
	CONSTRAINT password_min_length check (length(password) >= 8)
)

TABLESPACE pg_default;
ALTER TABLE public.users
    OWNER to postgres;
	
	
-------- BASIC FUNCTIONS START
CREATE OR REPLACE FUNCTION before_insert_or_update()   
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';
CREATE TRIGGER user_before_insert_or_update BEFORE INSERT OR UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE before_insert_or_update();



CREATE OR REPLACE FUNCTION before_insert()   
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';
CREATE TRIGGER user_before_insert BEFORE INSERT ON users FOR EACH ROW EXECUTE PROCEDURE before_insert();


CREATE OR REPLACE FUNCTION before_update()   
RETURNS TRIGGER AS $$
BEGIN
	IF NEW.id != OLD.id THEN
		NEW.id = OLD.id;
	END IF;	
	IF NEW.created_at != OLD.created_at THEN
		NEW.created_at = OLD.created_at;
	END IF;
    RETURN NEW;   
END;
$$ language 'plpgsql';
CREATE TRIGGER user_before_update BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE before_update();
-------- BASIC FUNCTIONS END



-- Table: public.tokens
-- DROP TABLE public.tokens;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE public.tokens
(
    id uuid PRIMARY KEY default gen_random_uuid(),
    user_id uuid NOT NULL,
    token text COLLATE pg_catalog."default" NOT NULL,
    last_ip character varying(32) COLLATE pg_catalog."default" DEFAULT 'unknown'::character varying,
    last_agent character varying(128) COLLATE pg_catalog."default" DEFAULT 'unknown'::character varying,
    auth_ip character varying(32) COLLATE pg_catalog."default" DEFAULT 'unknown'::character varying,
    auth_agent character varying(128) COLLATE pg_catalog."default" DEFAULT 'unknown'::character varying,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    CONSTRAINT fk_users_tokens FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE public.tokens
    OWNER to postgres;
	

-------- BASIC FUNCTIONS START
CREATE OR REPLACE FUNCTION before_insert_or_update()   
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';
CREATE TRIGGER token_before_insert_or_update BEFORE INSERT OR UPDATE ON tokens FOR EACH ROW EXECUTE PROCEDURE before_insert_or_update();



CREATE OR REPLACE FUNCTION before_insert()   
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';
CREATE TRIGGER token_before_insert BEFORE INSERT ON tokens FOR EACH ROW EXECUTE PROCEDURE before_insert();


CREATE OR REPLACE FUNCTION before_update()   
RETURNS TRIGGER AS $$
BEGIN
	IF NEW.id != OLD.id THEN
		NEW.id = OLD.id;
	END IF;	
	IF NEW.created_at != OLD.created_at THEN
		NEW.created_at = OLD.created_at;
	END IF;
    RETURN NEW;   
END;
$$ language 'plpgsql';
CREATE TRIGGER token_before_update BEFORE UPDATE ON tokens FOR EACH ROW EXECUTE PROCEDURE before_update();
-------- BASIC FUNCTIONS END


-- Table: public.files
-- DROP TABLE public.files;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE public.files
(
    id uuid PRIMARY KEY default gen_random_uuid(),
    user_id uuid NOT NULL,
    hash text COLLATE pg_catalog."default" NOT NULL,
    path text COLLATE pg_catalog."default" NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    original_name text COLLATE pg_catalog."default" NOT NULL DEFAULT 'unknown'::text,
    extension text COLLATE pg_catalog."default" DEFAULT 'unknown'::text,
    size text COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    CONSTRAINT fk_users_files FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE public.files
    OWNER to postgres;
	



-------- BASIC FUNCTIONS START
CREATE OR REPLACE FUNCTION before_insert_or_update()   
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';
CREATE TRIGGER file_before_insert_or_update BEFORE INSERT OR UPDATE ON files FOR EACH ROW EXECUTE PROCEDURE before_insert_or_update();



CREATE OR REPLACE FUNCTION before_insert()   
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';
CREATE TRIGGER file_before_insert BEFORE INSERT ON files FOR EACH ROW EXECUTE PROCEDURE before_insert();


CREATE OR REPLACE FUNCTION before_update()   
RETURNS TRIGGER AS $$
BEGIN
	IF NEW.id != OLD.id THEN
		NEW.id = OLD.id;
	END IF;	
	IF NEW.created_at != OLD.created_at THEN
		NEW.created_at = OLD.created_at;
	END IF;
    RETURN NEW;   
END;
$$ language 'plpgsql';
CREATE TRIGGER file_before_update BEFORE UPDATE ON files FOR EACH ROW EXECUTE PROCEDURE before_update();
-------- BASIC FUNCTIONS END


-- Table: public.articles
-- DROP TABLE public.articles;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE public.articles
(
    id uuid PRIMARY KEY default gen_random_uuid(),
    user_id uuid NOT NULL,
    is_published boolean DEFAULT false,
    title character varying(124) COLLATE pg_catalog."default" DEFAULT 'Без названия'::character varying,
    content text COLLATE pg_catalog."default" NOT NULL,
    slug text COLLATE pg_catalog."default" NOT NULL,
    published_at timestamp with time zone,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    CONSTRAINT articles_slug_key UNIQUE (slug),
    CONSTRAINT fk_users_articles FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
)

TABLESPACE pg_default;

ALTER TABLE public.articles
    OWNER to postgres;
	

-------- BASIC FUNCTIONS START
CREATE OR REPLACE FUNCTION before_insert_or_update()   
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';
CREATE TRIGGER article_before_insert_or_update BEFORE INSERT OR UPDATE ON articles FOR EACH ROW EXECUTE PROCEDURE before_insert_or_update();



CREATE OR REPLACE FUNCTION before_insert()   
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';
CREATE TRIGGER article_before_insert BEFORE INSERT ON articles FOR EACH ROW EXECUTE PROCEDURE before_insert();


CREATE OR REPLACE FUNCTION before_update()   
RETURNS TRIGGER AS $$
BEGIN
	IF NEW.id != OLD.id THEN
		NEW.id = OLD.id;
	END IF;	
	IF NEW.created_at != OLD.created_at THEN
		NEW.created_at = OLD.created_at;
	END IF;
    RETURN NEW;   
END;
$$ language 'plpgsql';
CREATE TRIGGER article_before_update BEFORE UPDATE ON articles FOR EACH ROW EXECUTE PROCEDURE before_update();
-------- BASIC FUNCTIONS END