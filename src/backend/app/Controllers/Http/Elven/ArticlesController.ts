import Env from '@ioc:Adonis/Core/Env'
import {HttpContextContract} from '@ioc:Adonis/Core/HttpContext'
import Article from "App/Models/Elven/Article"
import ArticleValidator from "App/Common/Elven/_VALIDATORS/ArticleValidator"


const pageSize = Env.get('PAGINATION_SIZE') // default: 16

export default class ArticlesController {

  // GET url/
  // params:
  // page = number
  // show = published, drafts, all
  // by = created, updated, published
  // start = newest (DESC), oldest (ASC)
  // preview = true (content < 480 symbols), false (gives you full articles)
  public async index(ctx: HttpContextContract) {
    const isAdmin = ctx['user'] && ctx['user'].role === 'admin'
    let validatedParams
    try {
      validatedParams = await ArticleValidator.requestParams(ctx.request, isAdmin)
    } catch (errors) {
      return ctx.response.status(400).send(errors)
    }
    let {page, show, by, start, preview} = validatedParams
    let articles
    if (show === 'published') {
      articles = await Article.query().where('is_published', 'true').orderBy(by, start).paginate(page, pageSize)
    } else if (show === 'drafts') {
      articles = await Article.query().where('is_published', 'false').orderBy(by, start).paginate(page, pageSize)
    } else if (show === 'all') {
      articles = await Article.query().orderBy(by, start).paginate(page, pageSize)
    }
    if (articles) {
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
    }
    return ctx.response.status(200).send(articles)
  }

  // GET url/:id
  public async show(ctx: HttpContextContract) {
    const article = await Article.find(ctx.params.id)
    if (article) {
      return ctx.response.status(200).send(article)
    } else {
      return ctx.response.status(404).send('Article not found.')
    }
  }

  // POST url/
  public async store(ctx: HttpContextContract) {
    let article = new Article()
    try {
      const validated = await ArticleValidator.whenCreate(ctx.request)
      Object.assign(article, validated)
    } catch (errors) {
      return ctx.response.status(400).send(errors)
    }
    const user = ctx['user']
    try {
      await user.related('articles').save(article)
      return ctx.response.status(200).send(article)
    } catch (error) {
      return ctx.response.status(500).send('Error while creating article.')
    }
  }

  // PUT OR PATCH url/:id
  public async update(ctx: HttpContextContract) {
    let article = await Article.find(ctx.params.id)
    if (!article) {
      return ctx.response.status(404).send('Article not found.')
    }
    try {
      Object.assign(article, await ArticleValidator.whenUpdate(ctx.request, article))
    } catch (errors) {
      return ctx.response.status(400).send(errors)
    }
    const user = ctx['user']
    try {
      await user.related('articles').save(article)
      return ctx.response.status(200).send(article)
    } catch (error) {
      return ctx.response.internalServerError('Error while saving article.')
    }
  }

  // DELETE url/:id
  public async destroy(ctx: HttpContextContract) {
    const article = await Article.find(ctx.params.id)
    if (!article) {
      return ctx.response.status(404).send('Article not found.')
    }
    try {
      await article.delete()
      return ctx.response.status(200).send('Article deleted.')
    } catch (error) {
      return ctx.response.status(500).send('Error while deleting article.')
    }
  }

}
