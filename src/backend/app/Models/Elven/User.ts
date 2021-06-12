import { DateTime } from 'luxon'
import {BaseModel, beforeSave, column, hasMany, HasMany} from '@ioc:Adonis/Lucid/Orm'
import Token from "App/Models/Elven/Token";
import Article from "App/Models/Elven/Article";
import File from "App/Models/Elven/File";
import Hooks from "App/Common/Elven/_MODEL_HOOKS/Hooks";

export default class User extends BaseModel {
  @column({ isPrimary: true })
  public id: number

  @hasMany(() => Token, {
    foreignKey: 'user_id',
  })
  public tokens: HasMany<typeof Token>

  @hasMany(() => Article, {
    foreignKey: 'author_id',
  })
  public articles: HasMany<typeof Article>

  @hasMany(() => File, {
    foreignKey: 'user_id',
  })
  public files: HasMany<typeof File>

  @column()
  public role: string

  @column()
  public username: string

  @column()
  public password: string

  @column()
  public reg_ip: string

  @column()
  public reg_agent: string

  @column.dateTime({ autoCreate: true })
  public created_at: DateTime

  @column.dateTime({ autoCreate: true, autoUpdate: true })
  public updated_at: DateTime

  @beforeSave()
  public static async hashPassword (user: User) {
    await Hooks.userValidate(user)
    user.password = await Hooks.hashPassword(user)
  }

}
