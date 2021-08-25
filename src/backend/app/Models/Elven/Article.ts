import { DateTime } from 'luxon'
import {BaseModel, beforeSave, column} from '@ioc:Adonis/Lucid/Orm'
import Hooks from "App/Common/Elven/_MODEL_HOOKS/Hooks";

export default class Article extends BaseModel {
  @column({ isPrimary: true })
  public id: string

  @column()
  public user_id: string

  @column()
  public is_published: boolean

  @column()
  public title: string

  @column()
  public content: string

  @column()
  public slug: string | undefined

  @column()
  public published_at: Date

  @column.dateTime()
  public created_at: DateTime

  @column.dateTime()
  public updated_at: DateTime

  @beforeSave()
  public static async hooksBefore (article: Article) {
    article.slug = Hooks.autoSlug(article)
    article.published_at = Hooks.setPublishedDate(article)
  }

  // @afterSave()
  // public static async hooksAfter (article: Article) {
  //   // тут можно добавить проверку на is_published
  //   // и если он будет true, то отправлять уведомление в Telegram о новой записи
  // }

}
