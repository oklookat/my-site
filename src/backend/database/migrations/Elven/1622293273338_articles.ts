import BaseSchema from '@ioc:Adonis/Lucid/Schema'

export default class Articles extends BaseSchema {
  protected tableName = 'articles'

  public async up () {
    this.schema.createTable(this.tableName, (table) => {
      table.increments('id')
      table
        .integer('author_id')
        .unsigned()
        .references('users.id')
        .onDelete('CASCADE')

      table.boolean('is_published').notNullable().defaultTo(false)
      table.text('thumbnail').nullable()
      table.string('title', 124).notNullable().defaultTo('Без названия')
      table.text('content').notNullable()
      table.text('slug').notNullable().unique()
      table.dateTime('published_at').nullable()

      table.timestamp('created_at', { useTz: true })
      table.timestamp('updated_at', { useTz: true })
    })
  }

  public async down () {
    this.schema.dropTable(this.tableName)
  }
}
