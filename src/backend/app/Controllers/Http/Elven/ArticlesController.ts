import Env from '@ioc:Adonis/Core/Env'
import {HttpContextContract} from '@ioc:Adonis/Core/HttpContext'
import Article from "App/Models/Elven/Article"
import ArticleValidator from "App/Common/Elven/_VALIDATORS/ArticleValidator"
import EL_Errors from "App/Common/Elven/_TOOLS/EL_Errors"


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
    // VALIDATION START //
    let show = ctx.request.input('show', 'published')
    show = show.toLowerCase()
    if (show !== 'published' && show !== 'drafts') {
      return ctx.response.status(400).send(await EL_Errors.publicError('«show» must be published, drafts or all.'))
    } else {
      if ((show === 'drafts' || show === 'all') && !isAdmin) {
        return ctx.response.status(403).send(await EL_Errors.publicError('Access denied.'))
      }
    }
    let by = ctx.request.input('by', 'published')
    by = by.toLowerCase()
    if (by !== 'created' && by !== 'published' && by !== 'updated') {
      return ctx.response.status(400).send(await EL_Errors.publicError('«by» must be created, published or updated.'))
    } else {
      if (by === 'created' || by === 'updated') {
        if (!isAdmin) {
          return ctx.response.status(403).send(await EL_Errors.publicError('Access denied.'))
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
      return ctx.response.status(400).send(await EL_Errors.publicError('«start» must be newest or oldest.'))
    } else {
      if (start === 'newest') {
        start = 'DESC'
      } else if (start === 'oldest') {
        start = 'ASC'
      }
    }
    let preview = ctx.request.input('preview', 'false')
    if(preview !== 'false' && preview !== 'true'){
      return ctx.response.status(400).send(await EL_Errors.publicError('«preview» must be true or false.'))
    } else if(preview === 'true'){
      preview = true
    } else if(preview === 'false'){
      preview = false
    }
    let page = ctx.request.input('page', 1)
    if(page < 1){
      return ctx.response.status(400).send(await EL_Errors.publicError('«page» cannot be less than one.'))
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
      return ctx.response.status(500).send(await EL_Errors.publicError('Error while getting articles.'))
    }
  }

  // GET url/:id
  public async show(ctx: HttpContextContract) {
    const article = await Article.find(ctx.params.id)
    if (article) {
      return ctx.response.status(200).send(article)
    } else {
      return ctx.response.notFound(await EL_Errors.publicError('Article not found.'))
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
      return ctx.response.internalServerError(await EL_Errors.publicError('Error while creating article.'))
    }
  }

  // PUT OR PATCH url/:id
  public async update(ctx: HttpContextContract) {
    let article = await Article.find(ctx.params.id)
    if (!article) {
      return ctx.response.notFound(await EL_Errors.publicError('Article not found.'))
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
      return ctx.response.internalServerError(await EL_Errors.publicError('Error while saving article.'))
    }
  }

  // DELETE url/:id
  public async destroy(ctx: HttpContextContract) {
    const article = await Article.find(ctx.params.id)
    if (!article) {
      return ctx.response.notFound(await EL_Errors.publicError('Article not found.'))
    }
    try {
      await article.delete()
      return ctx.response.status(200).send(await EL_Errors.publicError('Article deleted.'))
    } catch (error) {
      return ctx.response.internalServerError(await EL_Errors.publicError('Error while deleting article.'))
    }
  }

}
