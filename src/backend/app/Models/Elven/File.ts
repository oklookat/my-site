import { DateTime } from 'luxon'
import { BaseModel, column } from '@ioc:Adonis/Lucid/Orm'

export default class File extends BaseModel {
  @column({ isPrimary: true })
  public id: string

  @column()
  public user_id: string

  @column()
  public hash: string

  @column()
  public path: string

  @column()
  public original_name: string | undefined

  @column()
  public name: string

  @column()
  public extension: string | undefined

  @column()
  public size: string

  @column.dateTime()
  public created_at: DateTime

  @column.dateTime()
  public updated_at: DateTime
}
