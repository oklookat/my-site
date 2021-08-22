-- Table: public.users

-- DROP TABLE public.users;

DROP SEQUENCE IF EXISTS users_id_seq CASCADE;
CREATE SEQUENCE users_id_seq
    start 1
    increment 1
    NO MAXVALUE
    CACHE 1;
DROP TYPE IF EXISTS user_roles CASCADE;
CREATE TYPE user_roles AS ENUM ('user', 'admin');
CREATE TABLE public.users
(
    id bigint DEFAULT nextval('users_id_seq'::regclass) NOT NULL,
    role user_roles DEFAULT 'user'::user_roles NOT NULL,
    username character varying(24) COLLATE pg_catalog."default" NOT NULL,
    password character varying(64) COLLATE pg_catalog."default" NOT NULL,
    reg_ip character varying(32) COLLATE pg_catalog."default" DEFAULT 'unknown'::character varying,
    reg_agent character varying(128) COLLATE pg_catalog."default" DEFAULT 'unknown'::character varying,
    created_at timestamp with time zone default current_timestamp NOT NULL,
    updated_at timestamp with time zone,
    CONSTRAINT users_pkey PRIMARY KEY (id),
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