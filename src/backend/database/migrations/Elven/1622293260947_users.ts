import BaseSchema from '@ioc:Adonis/Lucid/Schema'

export default class Users extends BaseSchema {
  protected tableName = 'users'

  public async up() {
    this.schema.raw(`
      CREATE EXTENSION IF NOT EXISTS pgcrypto;
      ALTER TABLE IF EXISTS users DROP CONSTRAINT IF EXISTS user_length CASCADE;
      DROP TYPE IF EXISTS user_roles CASCADE;
      CREATE TYPE user_roles AS ENUM ('user', 'admin');
      `)
    this.schema.raw(`
      CREATE TABLE public.users
      (
        id         uuid PRIMARY KEY                                    default gen_random_uuid() NOT NULL,
        role       user_roles                                          DEFAULT 'user'::user_roles NOT NULL,
        username   character varying(24) COLLATE pg_catalog."default"                            NOT NULL,
        password   character varying(104) COLLATE pg_catalog."default"                            NOT NULL,
        reg_ip     character varying(54) COLLATE pg_catalog."default"  DEFAULT 'unknown':: character varying,
        reg_agent  varchar COLLATE pg_catalog."default" DEFAULT 'unknown'::varchar,
        created_at timestamp with time zone                            default current_timestamp NOT NULL,
        updated_at timestamp with time zone,
        CONSTRAINT users_username_key UNIQUE (username),
        CONSTRAINT users_length CHECK (length(username) >= 4 AND (length(password) >= 8))
      ) TABLESPACE pg_default;
      ALTER TABLE public.users
        OWNER to postgres;
    `)

    // basic functions. Actual for all models.
    this.schema.raw(`
    CREATE OR REPLACE FUNCTION before_insert_or_update()
    RETURNS TRIGGER AS $$
    BEGIN
      NEW.updated_at = now();
      RETURN NEW;
    END;
    $$ language 'plpgsql';


    CREATE OR REPLACE FUNCTION before_insert()
    RETURNS TRIGGER AS $$
    BEGIN
      NEW.created_at = now();
    RETURN NEW;
    END;
    $$ language 'plpgsql';


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

    CREATE TRIGGER user_before_insert_or_update BEFORE INSERT OR UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE before_insert_or_update();
    CREATE TRIGGER user_before_insert BEFORE INSERT ON users FOR EACH ROW EXECUTE PROCEDURE before_insert();
    CREATE TRIGGER user_before_update BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE before_update();
    `)
  }

  public async down() {
    this.schema.dropTable(this.tableName)
  }
}
