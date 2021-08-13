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
        // generate preview
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
                // if block text length > 408 we cut text and quit
                text = text.slice(0, 408) + '...'
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
      const notFound = new E_NOTFOUND(['articles'], 'Article not found.')
      return ctx.response.status(404).send(EL_ErrorCollector.singleError(notFound))
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
      // if(error.code){
      //   while(error.code === '23505'){
      //     // code 23505 mean what id exists
      //     // it happens when you import data direct to database
      //     // for example: in you imported data last id = 42, but in ORM or etc counter id = 1 (as I understand)
      //     // and he try create article, but in database article with id = 1 exists, and he gives errors
      //     // if you try to save again, he try save it, but with id = 2, etc. In this cycle he saving and saving again, to moment, when ORM counter go to actual id in imported data.
      //     // its not good solution, but it seems to work
      //     // and idk about performance. What if you have 100000 articles?
      //     // but what if really all id's busy? idk, maybe this cycle goes to infinity
      //     // uncomment this if you import data to DB, and now you try to create article, but he gives you 500 error. Maybe it helps.
      //     try {
      //       await user.related('articles').save(article)
      //       return ctx.response.status(200).send(article)
      //     } catch (error){
      //       if(error.code){
      //         if(error.code !== '23505'){
      //           break
      //         }
      //       }
      //     }
      //   }
      // }
      const unknown = new E_UNKNOWN(['articles'], 'Error while creating article.')
      return ctx.response.status(500).send(EL_ErrorCollector.singleError(unknown))
    }
  }

  // PUT OR PATCH url/:id
  public async update(ctx: HttpContextContract) {
    let article = await Article.find(ctx.params.id)
    if (!article) {
      const notFound = new E_NOTFOUND(['articles'], 'Article not found.')
      return ctx.response.status(404).send(EL_ErrorCollector.singleError(notFound))
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
      const unknown = new E_UNKNOWN(['articles'], 'Error while saving article.')
      return ctx.response.status(500).send(EL_ErrorCollector.singleError(unknown))
    }
  }

  // DELETE url/:id
  public async destroy(ctx: HttpContextContract) {
    const article = await Article.find(ctx.params.id)
    if (!article) {
      const notFound = new E_NOTFOUND(['articles'], 'Article not found.')
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
