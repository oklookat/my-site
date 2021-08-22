-- Table: public.files

-- DROP TABLE public.files;

DROP SEQUENCE IF EXISTS files_id_seq CASCADE;
CREATE SEQUENCE files_id_seq
    start 1
    increment 1
    NO MAXVALUE
    CACHE 1;
	

CREATE TABLE public.files
(
    id bigint NOT NULL DEFAULT nextval('files_id_seq'::regclass),
    user_id bigint NOT NULL,
    hash text COLLATE pg_catalog."default" NOT NULL,
    path text COLLATE pg_catalog."default" NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    original_name text COLLATE pg_catalog."default" NOT NULL DEFAULT 'unknown'::text,
    extension text COLLATE pg_catalog."default" DEFAULT 'unknown'::text,
    size text COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    CONSTRAINT files_pkey PRIMARY KEY (id),
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