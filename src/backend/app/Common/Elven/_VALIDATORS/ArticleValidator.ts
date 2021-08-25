import {RequestContract} from "@ioc:Adonis/Core/Request"
import Article from "App/Models/Elven/Article"
import {EL_ErrorCollector} from "App/Common/Elven/_ERRORS/EL_ErrorCollector"
import {
  E_AUTH_FORBIDDEN,
  E_VALIDATION_EMPTY,
  E_VALIDATION_MINMAX,
  E_VALIDATION_MUSTBE
} from "App/Common/Elven/_ERRORS/EL_Errors";

export default class ArticleValidator {

  public static whenCreate(request: RequestContract) {
    const errorCollector = new EL_ErrorCollector()
    let {is_published, title, content} = request.all()
    if (!is_published) {
      is_published = false
    }
    if (is_published !== 'true' && is_published !== 'false' && is_published !== true && is_published !== false) {
      const mustbe = new E_VALIDATION_MUSTBE(['is_published'], ['true', 'false'])
      errorCollector.addError(mustbe)
    } else if (is_published === 'true') {
      is_published = true
    } else if (is_published === 'false') {
      is_published = false
    }
    if (!title) {
      title = 'Untitled'
    } else {
      title = title.toString()
      if (title.length > 124) {
        const minmax = new E_VALIDATION_MINMAX(['title'], 0, 124)
        errorCollector.addError(minmax)
      }
    }
    if (!content) {
      const empty = new E_VALIDATION_EMPTY(['content'])
      errorCollector.addError(empty)
    }
    if (errorCollector.hasErrors()) {
      throw errorCollector.getErrors()
    }
    return {is_published: is_published, title: title, content: content}
  }

  public static whenUpdate(request: RequestContract, foundArticle: Article) {
    const errorCollector = new EL_ErrorCollector()
    let {is_published, title, content} = request.all()
    if (is_published !== true && is_published !== false) {
      is_published = foundArticle.is_published
    }
    if (!title) {
      title = foundArticle.title
    } else {
      if (title.length > 124) {
        const minmax = new E_VALIDATION_MINMAX(['title'], 0, 124)
        errorCollector.addError(minmax)
      }
    }
    if (!content) {
      content = foundArticle.content
    }
    if (errorCollector.hasErrors()) {
      throw errorCollector.getErrors()
    }
    return {is_published: is_published, title: title, content: content}
  }


  public static requestParams(request: RequestContract, isAdmin: boolean) {
    const errorCollector = new EL_ErrorCollector()
    let show = request.input('show', 'published')
    show = show.toLowerCase()
    if (show !== 'published' && show !== 'drafts') {
      const mustbe = new E_VALIDATION_MUSTBE(['show'], ['drafts', 'all'])
      errorCollector.addError(mustbe)
    } else {
      if ((show === 'drafts' || show === 'all') && !isAdmin) {
        const forbidden = new E_AUTH_FORBIDDEN(['show'])
        errorCollector.addError(forbidden)
      }
    }
    let by = request.input('by', 'published')
    by = by.toLowerCase()
    if (by !== 'created' && by !== 'published' && by !== 'updated') {
      const mustbe = new E_VALIDATION_MUSTBE(['by'], ['created', 'published', 'updated'])
      errorCollector.addError(mustbe)
    } else {
      if (by === 'created' || by === 'updated') {
        if (!isAdmin) {
          const forbidden = new E_AUTH_FORBIDDEN(['by'])
          errorCollector.addError(forbidden)
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
    let start = request.input('start', 'newest')
    start = start.toLowerCase()
    if (start !== 'newest' && start !== 'oldest') {
      const mustbe = new E_VALIDATION_MUSTBE(['start'], ['newest', 'oldest'])
      errorCollector.addError(mustbe)
    } else {
      if (start === 'newest') {
        start = 'DESC'
      } else if (start === 'oldest') {
        start = 'ASC'
      }
    }
    let preview = request.input('preview', 'true')
    preview = preview.toLowerCase()
    if (preview !== 'false' && preview !== 'true') {
      const mustbe = new E_VALIDATION_MUSTBE(['preview'], ['true', 'false'])
      errorCollector.addError(mustbe)
    } else if (preview === 'true') {
      preview = true
    } else if (preview === 'false') {
      preview = false
    }
    let page = request.input('page', 1)
    if (page < 1) {
      const minmax = new E_VALIDATION_MINMAX(['page'], 1)
      errorCollector.addError(minmax)
    }
    if (errorCollector.hasErrors()) {
      throw errorCollector.getErrors()
    }
    return {page: page, show: show, by: by, start: start, preview: preview}
  }

}
