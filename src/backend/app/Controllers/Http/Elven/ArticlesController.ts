import {HttpContextContract} from '@ioc:Adonis/Core/HttpContext'
import Article from "App/Models/Elven/Article";
import ElvenValidators from "App/Common/Elven/_VALIDATORS/ElvenValidators"
import ElvenTools from "App/Common/Elven/_TOOLS/ElvenTools"

const pageSize = 16

export default class ArticlesController {

  // GET url/
  // params:
  // show = published, drafts, all
  // by = created, published, title
  // start = up (DESC), down (ASC)
  public async index(ctx: HttpContextContract) {
    const isAdmin = ctx['user'] && ctx['user'].role === 'admin'
    // VALIDATION START //
    let show = ctx.request.input('show', 'published')
    show = show.toLowerCase()
    if (show !== 'published' && show !== 'drafts') {
      return ctx.response.status(400).send(await ElvenTools.publicErrorConstructor('show должен быть published, drafts или all'))
    } else {
      if ((show === 'drafts' || show === 'all') && !isAdmin) {
        return ctx.response.status(403).send(await ElvenTools.publicErrorConstructor('Доступ запрещен.'))
      }
    }
    let by = ctx.request.input('by', 'created')
    by = by.toLowerCase()
    if (by !== 'created' && by !== 'published' && by !== 'title') {
      return ctx.response.status(400).send(await ElvenTools.publicErrorConstructor('by должен быть created, published или title'))
    } else {
      if (by === 'created') {
        if (!isAdmin) {
          return ctx.response.status(403).send(await ElvenTools.publicErrorConstructor('Доступ запрещен.'))
        } else {
          by = 'created_at'
        }
      } else if (by === 'published') {
        by = 'published_at'
      }
    }
    let start = ctx.request.input('start', 'up')
    start = start.toLowerCase()
    if (start !== 'up' && start !== 'down') {
      return ctx.response.status(400).send(await ElvenTools.publicErrorConstructor('start должен быть up или down'))
    } else {
      if (start === 'up') {
        start = 'DESC'
      } else if (start === 'down') {
        start = 'ASC'
      }
    }
    // VALIDATION END //
    let page = ctx.request.input('page', 1)
    if (show === 'published') {
      const articles = await Article.query().where('is_published', 'true').orderBy(by, start).paginate(page, pageSize)
      return ctx.response.status(200).send(articles)
    } else if (show === 'drafts') {
      const articles = await Article.query().where('is_published', 'false').orderBy(by, start).paginate(page, pageSize)
      return ctx.response.status(200).send(articles)
    } else if (show === 'all') {
      const articles = await Article.query().orderBy(by, start).paginate(page, pageSize)
      return ctx.response.status(200).send(articles)
    }
    return ctx.response.status(500).send(await ElvenTools.publicErrorConstructor('Произошла ошибка.'))
  }

  // GET url/:id
  public async show(ctx: HttpContextContract) {
    const article = await Article.find(ctx.params.id)
    if (article) {
      return ctx.response.status(200).send(article)
    } else {
      return ctx.response.notFound(await ElvenTools.publicErrorConstructor('Запись не найдена.'))
    }
  }

  // POST url/
  public async store(ctx: HttpContextContract) {
    let article = new Article()
    try {
      Object.assign(article, await ElvenValidators.articleValidateCreate(ctx.request))
    } catch (error) {
      return ctx.response.badRequest(error)
    }
    const user = ctx['user']
    try {
      await user.related('articles').save(article)
      return ctx.response.status(200).send(article)
    } catch (error) {
      return ctx.response.internalServerError(await ElvenTools.publicErrorConstructor('Не удалось создать запись.'))
    }
  }

  // PUT OR PATCH url/:id
  public async update(ctx: HttpContextContract) {
    let article = await Article.find(ctx.params.id)
    if (!article) {
      return ctx.response.notFound(await ElvenTools.publicErrorConstructor('Запись не найдена.'))
    }
    try {
      Object.assign(article, await ElvenValidators.articleValidateUpdate(ctx.request, article))
    } catch (error) {
      return ctx.response.badRequest(error)
    }
    const user = ctx['user']
    try {
      await user.related('articles').save(article)
      return ctx.response.status(200).send(article)
    } catch (error) {
      console.log(error)
      return ctx.response.internalServerError(await ElvenTools.publicErrorConstructor('Не удалось сохранить запись.'))
    }
  }

  // DELETE url/:id
  public async destroy(ctx: HttpContextContract) {
    const article = await Article.find(ctx.params.id)
    if (!article) {
      return ctx.response.notFound(await ElvenTools.publicErrorConstructor('Запись не найдена.'))
    }
    try {
      await article.delete()
      return ctx.response.status(200).send(await ElvenTools.publicSuccessConstructor('Запись удалена.'))
    } catch (error) {
      return ctx.response.internalServerError(await ElvenTools.publicErrorConstructor('Ошибка при удалении записи.'))
    }
  }

}
