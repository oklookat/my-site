-- Table: public.articles

-- DROP TABLE public.articles;

DROP SEQUENCE IF EXISTS articles_id_seq CASCADE;
CREATE SEQUENCE articles_id_seq
    start 1
    increment 1
    NO MAXVALUE
    CACHE 1;
CREATE TABLE public.articles
(
    id bigint NOT NULL DEFAULT nextval('articles_id_seq'::regclass),
    user_id bigint NOT NULL,
    is_published boolean DEFAULT false,
    title character varying(124) COLLATE pg_catalog."default" DEFAULT 'Без названия'::character varying,
    content text COLLATE pg_catalog."default" NOT NULL,
    slug text COLLATE pg_catalog."default" NOT NULL,
    published_at timestamp with time zone,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    CONSTRAINT articles_pkey PRIMARY KEY (id),
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