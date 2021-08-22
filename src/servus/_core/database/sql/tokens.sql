-- Table: public.tokens

-- DROP TABLE public.tokens;

DROP SEQUENCE IF EXISTS tokens_id_seq CASCADE;
CREATE SEQUENCE tokens_id_seq
    start 1
    increment 1
    NO MAXVALUE
    CACHE 1;
CREATE TABLE public.tokens
(
    id bigint NOT NULL DEFAULT nextval('tokens_id_seq'::regclass),
    user_id bigint NOT NULL,
    token text COLLATE pg_catalog."default" NOT NULL,
    last_ip character varying(32) COLLATE pg_catalog."default" DEFAULT 'unknown'::character varying,
    last_agent character varying(128) COLLATE pg_catalog."default" DEFAULT 'unknown'::character varying,
    auth_ip character varying(32) COLLATE pg_catalog."default" DEFAULT 'unknown'::character varying,
    auth_agent character varying(128) COLLATE pg_catalog."default" DEFAULT 'unknown'::character varying,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    CONSTRAINT tokens_pkey PRIMARY KEY (id),
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