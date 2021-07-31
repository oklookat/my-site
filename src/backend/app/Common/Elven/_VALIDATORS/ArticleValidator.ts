import {RequestContract} from "@ioc:Adonis/Core/Request";
import validator from "validator";
import Article from "App/Models/Elven/Article";
import {EL_ErrorCollector} from "App/Common/Elven/_TOOLS/EL_Errors"

export default class ArticleValidator {

  public static async whenCreate(request: RequestContract) {
    const errorCollector = new EL_ErrorCollector()
    let {is_published, thumbnail, title, content} = request.all()
    if (!is_published) {
      is_published = false
    }
    if (!thumbnail) {
      thumbnail = null
    }
    if (!title) {
      title = 'Без названия'
    } else {
      title = title.toString()
      if (title.length > 124) {
        errorCollector.addError('VALIDATION', 400, '«title» must be less than 124 symbols.')
      }
    }
    if (!content) {
      errorCollector.addError('VALIDATION', 400, '«content» can not be empty.')
    }
    is_published = is_published.toString()
    if (!validator.isBoolean(is_published)) {
      errorCollector.addError('VALIDATION', 400, '«isPublished» must be type bool.')
    }
    if (typeof content !== 'object' && content !== null) {
      errorCollector.addError('VALIDATION', 400, '«content» must be an object.')
    }
    if(errorCollector.hasErrors()){
      return Promise.reject(errorCollector.getErrors())
    }
    return Promise.resolve({is_published: is_published, thumbnail: thumbnail, title: title, content: content})
  }

  public static async whenUpdate(request: RequestContract, foundArticle: Article) {
    const errorCollector = new EL_ErrorCollector()
    let {is_published, thumbnail, title, content} = request.all()
    if (is_published !== true && is_published !== false) {
      is_published = foundArticle.is_published
    }
    if (!thumbnail) {
      thumbnail = foundArticle.thumbnail
    }
    if (!title) {
      title = foundArticle.title
    } else {
      if (title.length > 124) {
        errorCollector.addError('VALIDATION', 400, '«title» must be less than 124 symbols.')
      }
    }
    if (!content) {
      content = foundArticle.content
    }

    if (errorCollector.hasErrors()) {
      return Promise.reject(errorCollector.getErrors())
    }
    return Promise.resolve({is_published: is_published, thumbnail: thumbnail, title: title, content: content})
  }


  public static async requestParams(request: RequestContract, isAdmin: boolean) {
    const errorCollector = new EL_ErrorCollector()
    let show = request.input('show', 'published')
    show = show.toLowerCase()
    if (show !== 'published' && show !== 'drafts') {
      errorCollector.addError('VALIDATION', 400, '«show» must be published, drafts or all.')
    } else {
      if ((show === 'drafts' || show === 'all') && !isAdmin) {
        errorCollector.addError('VALIDATION', 403, '«drafts» or «all» can see only admin.')
      }
    }
    let by = request.input('by', 'published')
    by = by.toLowerCase()
    if (by !== 'created' && by !== 'published' && by !== 'updated') {
      errorCollector.addError('VALIDATION', 400, '«by» must be created, published or updated.')
    } else {
      if (by === 'created' || by === 'updated') {
        if (!isAdmin) {
          errorCollector.addError('VALIDATION', 403, 'by «created» or «updated» can see only admin.')
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
      errorCollector.addError('VALIDATION', 400, '«start» must be newest or oldest.')
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
      errorCollector.addError('VALIDATION', 400, '«preview» must be true or false.')
    } else if (preview === 'true') {
      preview = true
    } else if (preview === 'false') {
      preview = false
    }
    let page = request.input('page', 1)
    if (page < 1) {
      errorCollector.addError('VALIDATION', 400, '«page» cannot be less than one.')
    }

    if (errorCollector.hasErrors()) {
      return Promise.reject(errorCollector.getErrors())
    }
    return Promise.resolve({page: page, show: show, by: by, start: start, preview: preview})
  }

}
