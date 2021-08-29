import Env from '@ioc:Adonis/Core/Env'
import {HttpContextContract} from '@ioc:Adonis/Core/HttpContext'
import Article from "App/Models/Elven/Article"
import ArticleValidator from "App/Common/Elven/_VALIDATORS/ArticleValidator"
import {EL_ErrorCollector} from "App/Common/Elven/_ERRORS/EL_ErrorCollector"
import {E_NOTFOUND, E_UNKNOWN} from "App/Common/Elven/_ERRORS/EL_Errors";


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
    let isAdmin = false
    if(ctx.user){
      if(ctx.user.role === 'admin'){
        isAdmin = true
      }
    }
    let validatedParams
    try {
      validatedParams = ArticleValidator.requestParams(ctx.request, isAdmin)
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
        articles = await this.generatePreview(articles)
      }
    }
    return ctx.response.status(200).send(articles)
  }

  private async generatePreview(articles){
    for (let i = 0; articles.length > i; i++) {
      // articles[i] = one article
      let content = JSON.parse(articles[i].content)
      for(let i = 0; content.blocks.length > i; i++){
        // content.blocks[i] = one block in article content
        const block = content.blocks[i]
        const blockType = block.type
        if(blockType === 'paragraph'){
          let text = block.data.text
          if (text.length > 408) {
            // if block text length > 304 we cut text and quit
            text = text.slice(0, 304) + '...'
            content.blocks[i].data.text = text
          }
          // cut all blocks except paragraph
          content.blocks = {0: content.blocks[i]}
          break
        } else {
          delete content.blocks[i]
        }
      }
      articles[i].content = content
    }
    return Promise.resolve(articles)
  }

  // GET url/:id
  public async show(ctx: HttpContextContract) {
    const article = await Article.find(ctx.params.id)
    if (article) {
      return ctx.response.status(200).send(article)
    } else {
      const notFound = new E_NOTFOUND(['articles'])
      return ctx.response.status(404).send(EL_ErrorCollector.singleError(notFound))
    }
  }

  // POST url/
  public async store(ctx: HttpContextContract) {
    let article = new Article()
    try {
      const validated = ArticleValidator.whenCreate(ctx.request)
      Object.assign(article, validated)
    } catch (errors) {
      return ctx.response.status(400).send(errors)
    }
    try {
      const user = ctx.user
      if (!user){
        throw 'PIPE_MISSING_USER'
      }
      await user.related('articles').save(article)
      return ctx.response.status(200).send(article)
    } catch (error) {
      const unknown = new E_UNKNOWN(['articles'], 'Error while creating article.')
      return ctx.response.status(500).send(EL_ErrorCollector.singleError(unknown))
    }
  }

  // PUT OR PATCH url/:id
  public async update(ctx: HttpContextContract) {
    let article = await Article.find(ctx.params.id)
    if (!article) {
      const notFound = new E_NOTFOUND(['articles'])
      return ctx.response.status(404).send(EL_ErrorCollector.singleError(notFound))
    }
    try {
      Object.assign(article, ArticleValidator.whenUpdate(ctx.request, article))
    } catch (errors) {
      return ctx.response.status(400).send(errors)
    }
    try {
      const user = ctx.user
      if(!user){
        throw 'PIPE_MISSING_USER'
      }
      await user.related('articles').save(article)
      return ctx.response.status(200).send(article)
    } catch (error) {
      const unknown = new E_UNKNOWN(['articles'], 'Error while saving article.')
      return ctx.response.status(500).send(EL_ErrorCollector.singleError(unknown))
    }
  }

  // DELETE url/:id
  public async destroy(ctx: HttpContextContract) {
    const article = await Article.find(ctx.params.id)
    if (!article) {
      const notFound = new E_NOTFOUND(['articles'])
      return ctx.response.status(404).send(EL_ErrorCollector.singleError(notFound))
    }
    try {
      await article.delete()
      return ctx.response.status(200).send('Article deleted.')
    } catch (error) {
      const unknown = new E_UNKNOWN(['articles'], 'Error while deleting article.')
      return ctx.response.status(500).send(EL_ErrorCollector.singleError(unknown))
    }
  }

}
