import BaseSchema from '@ioc:Adonis/Lucid/Schema'

export default class Files extends BaseSchema {
  protected tableName = 'files'

  public async up () {
    this.schema.createTable(this.tableName, (table) => {
      table.increments('id')

      table.text('hash').notNullable()
      table.text('path').notNullable()
      table.text('original_name').notNullable().defaultTo('unknown')
      table.text('name').notNullable()
      table.string('extension').notNullable().defaultTo('unknown')
      table.string('size').notNullable()

      table
        .integer('user_id')
        .unsigned()
        .references('users.id')
        .onDelete('CASCADE')

      table.timestamp('created_at', { useTz: true })
      table.timestamp('updated_at', { useTz: true })
    })
  }

  public async down () {
    this.schema.dropTable(this.tableName)
  }
}
