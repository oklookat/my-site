import BaseSchema from '@ioc:Adonis/Lucid/Schema'

export default class Tokens extends BaseSchema {
  protected tableName = 'tokens'

  public async up() {
    this.schema.raw(`
      CREATE TABLE public.tokens
      (
        id         uuid PRIMARY KEY                                   default gen_random_uuid() NOT NULL,
        user_id    uuid                                                                         NOT NULL,
        token      character varying(4024) COLLATE pg_catalog."default" NOT NULL,
        last_ip    character varying(32) COLLATE pg_catalog."default" DEFAULT 'unknown':: character varying,
        last_agent varchar COLLATE pg_catalog."default"               DEFAULT 'unknown'::varchar,
        auth_ip    character varying(32) COLLATE pg_catalog."default" DEFAULT 'unknown':: character varying,
        auth_agent varchar COLLATE pg_catalog."default"               DEFAULT 'unknown'::varchar,
        created_at timestamp with time zone                           default current_timestamp,
        updated_at timestamp with time zone,
        CONSTRAINT fk_users_tokens FOREIGN KEY (user_id)
          REFERENCES public.users (id) MATCH SIMPLE
          ON UPDATE NO ACTION
          ON DELETE CASCADE
      ) TABLESPACE pg_default;
      ALTER TABLE public.tokens
        OWNER to postgres;
    `)

    this.schema.raw(`
    CREATE TRIGGER token_before_insert_or_update BEFORE INSERT OR UPDATE ON tokens FOR EACH ROW EXECUTE PROCEDURE before_insert_or_update();
    CREATE TRIGGER token_before_insert BEFORE INSERT ON tokens FOR EACH ROW EXECUTE PROCEDURE before_insert();
    CREATE TRIGGER token_before_update BEFORE UPDATE ON tokens FOR EACH ROW EXECUTE PROCEDURE before_update();
    `)
  }

  public async down() {
    this.schema.dropTable(this.tableName)
  }
}
