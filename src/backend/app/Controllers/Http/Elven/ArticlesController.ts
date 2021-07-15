import {HttpContextContract} from '@ioc:Adonis/Core/HttpContext'
import Article from "App/Models/Elven/Article";
import ElvenValidators from "App/Common/Elven/_VALIDATORS/ElvenValidators"
import ElvenTools from "App/Common/Elven/_TOOLS/ElvenTools"

const pageSize = 16

export default class ArticlesController {

  // GET url/
  public async index(ctx: HttpContextContract) {
    let page = ctx.request.input('page', 1)
    let order = ctx.request.input('order', 'created_at')
    let direction = ctx.request.input('direction', 'desc')
    let articles
    try {
      if (ctx['user'] && ctx['user'].role === 'admin') {
        articles = await Article.query().orderBy(order, direction).paginate(page, pageSize)
      } else {
        articles = await Article.query().where('is_published', 'true').orderBy(order, direction).paginate(page, pageSize)
      }
    } catch (error) {
      return ctx.response.forbidden(await ElvenTools.publicErrorConstructor('При получении записей произошла ошибка. Обратитесь к администратору.'))
    }
    return ctx.response.status(200).send(articles)
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
    await this.postOrUpdate(ctx, true)
  }

  // PUT OR PATCH url/:id
  public async update(ctx: HttpContextContract) {
    await this.postOrUpdate(ctx, false)
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
      return ctx.response.badGateway(await ElvenTools.publicErrorConstructor('Ошибка при удалении записи. Обратитесь к администратору.'))
    }
  }


  private async postOrUpdate(ctx: HttpContextContract, isPost: boolean) {
    let data
    try {
      data = await ElvenValidators.articleValidate(ctx.request)
    } catch (error) {
      return ctx.response.badRequest(error)
    }
    let article
    if (isPost) {
      article = new Article()
    } else {
      article = await Article.find(ctx.params.id)
      if (!article) {
        return ctx.response.notFound(await ElvenTools.publicErrorConstructor('Запись не найдена.'))
      }
    }
    article.is_published = data.is_published
    if (data.thumbnail) {
      article.thumbnail = data.thumbnail
    }
    if (data.title) {
      article.title = data.title
    }
    article.content = data.content
    const user = ctx['user']
    try {
      await user.related('articles').save(article)
      return ctx.response.status(200).send(article)
    } catch (error) {
      console.log(error)
      return ctx.response.badGateway(await ElvenTools.publicErrorConstructor('Не удалось сохранить запись.'))
    }
  }
}
