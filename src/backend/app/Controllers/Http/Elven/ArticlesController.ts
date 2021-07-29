import Env from '@ioc:Adonis/Core/Env'
import {HttpContextContract} from '@ioc:Adonis/Core/HttpContext'
import Article from "App/Models/Elven/Article"
import ArticleValidator from "App/Common/Elven/_VALIDATORS/ArticleValidator"
import ErrorConstructors from "App/Common/Elven/_TOOLS/ErrorConstructors"


const pageSize = Env.get('PAGINATION_SIZE') // default: 16

export default class ArticlesController {

  // GET url/
  // params:
  // show = published, drafts, all
  // by = created, updated, published
  // start = newest (DESC), oldest (ASC)
  // preview = true (content < 480 symbols), false (gives you full articles)
  public async index(ctx: HttpContextContract) {
    const isAdmin = ctx['user'] && ctx['user'].role === 'admin'
    // VALIDATION START //
    let show = ctx.request.input('show', 'published')
    show = show.toLowerCase()
    if (show !== 'published' && show !== 'drafts') {
      return ctx.response.status(400).send(await ErrorConstructors.publicError('show должен быть published, drafts или all'))
    } else {
      if ((show === 'drafts' || show === 'all') && !isAdmin) {
        return ctx.response.status(403).send(await ErrorConstructors.publicError('Доступ запрещен.'))
      }
    }
    let by = ctx.request.input('by', 'published')
    by = by.toLowerCase()
    if (by !== 'created' && by !== 'published' && by !== 'updated') {
      return ctx.response.status(400).send(await ErrorConstructors.publicError('by должен быть created, published или updated'))
    } else {
      if (by === 'created' || by === 'updated') {
        if (!isAdmin) {
          return ctx.response.status(403).send(await ErrorConstructors.publicError('Доступ запрещен.'))
        } else {
          if (by === 'created') {
            by = 'created_at'
          } else if (by === 'updated') {
            by = 'updated_at'
          }
        }
      } else if (by === 'published') {
        by = 'published_at'
      }
    }
    let start = ctx.request.input('start', 'newest')
    start = start.toLowerCase()
    if (start !== 'newest' && start !== 'oldest') {
      return ctx.response.status(400).send(await ErrorConstructors.publicError('start должен быть newest или oldest'))
    } else {
      if (start === 'newest') {
        start = 'DESC'
      } else if (start === 'oldest') {
        start = 'ASC'
      }
    }
    let preview = ctx.request.input('preview', 'false')
    if(preview !== 'false' && preview !== 'true'){
      return ctx.response.status(400).send(await ErrorConstructors.publicError('preview должно быть true или false'))
    } else if(preview === 'true'){
      preview = true
    } else if(preview === 'false'){
      preview = false
    }
    let page = ctx.request.input('page', 1)
    if(page < 1){
      return ctx.response.status(400).send(await ErrorConstructors.publicError('page не может быть меньше нуля'))
    }
    // VALIDATION END //
    let articles
    if (show === 'published') {
      articles = await Article.query().where('is_published', 'true').orderBy(by, start).paginate(page, pageSize)
    } else if (show === 'drafts') {
      articles = await Article.query().where('is_published', 'false').orderBy(by, start).paginate(page, pageSize)
    } else if (show === 'all') {
      articles = await Article.query().orderBy(by, start).paginate(page, pageSize)
    }
    if (articles) {
      let preview = ctx.request.input('preview', 'false')
      preview.toLowerCase()
      if (preview) {
        for (let i = 0; articles.length > i; i++) {
          let content = articles[i].content
          content = JSON.parse(content)
          content = content.blocks[0].data.text
          if (content.length > 408) {
            content = content.slice(0, 408) + '...'
            articles[i].content = content
          } else {
            articles[i].content = content
          }
        }
      }
      return ctx.response.status(200).send(articles)
    } else {
      return ctx.response.status(500).send(await ErrorConstructors.publicError('При получении записей произошла ошибка.'))
    }
  }

  // GET url/:id
  public async show(ctx: HttpContextContract) {
    const article = await Article.find(ctx.params.id)
    if (article) {
      return ctx.response.status(200).send(article)
    } else {
      return ctx.response.notFound(await ErrorConstructors.publicError('Запись не найдена.'))
    }
  }

  // POST url/
  public async store(ctx: HttpContextContract) {
    let article = new Article()
    try {
      Object.assign(article, await ArticleValidator.whenCreate(ctx.request))
    } catch (error) {
      return ctx.response.badRequest(error)
    }
    const user = ctx['user']
    try {
      await user.related('articles').save(article)
      return ctx.response.status(200).send(article)
    } catch (error) {
      return ctx.response.internalServerError(await ErrorConstructors.publicError('Не удалось создать запись.'))
    }
  }

  // PUT OR PATCH url/:id
  public async update(ctx: HttpContextContract) {
    let article = await Article.find(ctx.params.id)
    if (!article) {
      return ctx.response.notFound(await ErrorConstructors.publicError('Запись не найдена.'))
    }
    try {
      Object.assign(article, await ArticleValidator.whenUpdate(ctx.request, article))
    } catch (error) {
      return ctx.response.badRequest(error)
    }
    const user = ctx['user']
    try {
      await user.related('articles').save(article)
      return ctx.response.status(200).send(article)
    } catch (error) {
      return ctx.response.internalServerError(await ErrorConstructors.publicError('Не удалось сохранить запись.'))
    }
  }

  // DELETE url/:id
  public async destroy(ctx: HttpContextContract) {
    const article = await Article.find(ctx.params.id)
    if (!article) {
      return ctx.response.notFound(await ErrorConstructors.publicError('Запись не найдена.'))
    }
    try {
      await article.delete()
      return ctx.response.status(200).send(await ErrorConstructors.publicError('Запись удалена.'))
    } catch (error) {
      return ctx.response.internalServerError(await ErrorConstructors.publicError('Ошибка при удалении записи.'))
    }
  }

}
