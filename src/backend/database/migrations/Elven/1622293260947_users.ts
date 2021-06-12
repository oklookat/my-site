import BaseSchema from '@ioc:Adonis/Lucid/Schema'

export default class Users extends BaseSchema {
  protected tableName = 'users'

  public async up () {
    this.schema.createTable(this.tableName, (table) => {
      table.increments('id').primary()

      table.string('role').notNullable().defaultTo('user')
      table.string('username').notNullable().unique()
      table.text('password').notNullable()
      table.text('reg_ip').nullable().defaultTo('unknown')
      table.text('reg_agent').nullable().defaultTo('unknown')

      table.timestamp('created_at', { useTz: true })
      table.timestamp('updated_at', { useTz: true })
    })
  }

  public async down () {
    this.schema.dropTable(this.tableName)
  }
}
