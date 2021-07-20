import { DateTime } from 'luxon'
import { BaseModel, column } from '@ioc:Adonis/Lucid/Orm'

export default class Token extends BaseModel {
  @column({ isPrimary: true })
  public id: number

  @column()
  public user_id: number

  @column()
  public last_ip: string

  @column()
  public last_agent: string | undefined

  @column()
  public auth_ip: string

  @column()
  public auth_agent: string | undefined

  @column()
  public token: string

  @column.dateTime({ autoCreate: true })
  public created_at: DateTime

  @column.dateTime({ autoCreate: true, autoUpdate: true })
  public updated_at: DateTime
}
