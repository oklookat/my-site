import BaseSchema from '@ioc:Adonis/Lucid/Schema'

export default class Files extends BaseSchema {
  protected tableName = 'files'

  public async up() {
    this.schema.raw(`
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
    `)

    this.schema.raw(`
    CREATE TRIGGER file_before_insert_or_update BEFORE INSERT OR UPDATE ON files FOR EACH ROW EXECUTE PROCEDURE before_insert_or_update();
    CREATE TRIGGER file_before_insert BEFORE INSERT ON files FOR EACH ROW EXECUTE PROCEDURE before_insert();
    CREATE TRIGGER file_before_update BEFORE UPDATE ON files FOR EACH ROW EXECUTE PROCEDURE before_update();
    `)
  }

  public async down() {
    this.schema.dropTable(this.tableName)
  }
}
