import BaseSchema from '@ioc:Adonis/Lucid/Schema'

export default class Articles extends BaseSchema {
  protected tableName = 'articles'

  public async up() {
    this.schema.raw(`
      CREATE TABLE public.articles
      (
        id           uuid PRIMARY KEY                                    default gen_random_uuid() NOT NULL,
        user_id      uuid                                                                          NOT NULL,
        is_published boolean                                             DEFAULT false,
        title        character varying(124) COLLATE pg_catalog."default" DEFAULT 'Untitled':: character varying,
        content      text COLLATE pg_catalog."default"                                             NOT NULL,
        slug         character varying(256) COLLATE pg_catalog."default",
        published_at timestamp with time zone,
        created_at   timestamp with time zone                            default current_timestamp,
        updated_at   timestamp with time zone,
        CONSTRAINT articles_slug_key UNIQUE (slug),
        CONSTRAINT fk_users_articles FOREIGN KEY (user_id)
          REFERENCES public.users (id) MATCH SIMPLE
          ON UPDATE NO ACTION
          ON DELETE CASCADE
      ) TABLESPACE pg_default;

      ALTER TABLE public.articles
        OWNER to postgres;
    `)

    this.schema.raw(`
    CREATE TRIGGER article_before_insert_or_update BEFORE INSERT OR UPDATE ON articles FOR EACH ROW EXECUTE PROCEDURE before_insert_or_update();
    CREATE TRIGGER article_before_insert BEFORE INSERT ON articles FOR EACH ROW EXECUTE PROCEDURE before_insert();
    CREATE TRIGGER article_before_update BEFORE UPDATE ON articles FOR EACH ROW EXECUTE PROCEDURE before_update();
    `)
  }

  public async down() {
    this.schema.dropTable(this.tableName)
  }
}
