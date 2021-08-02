import {RequestContract} from "@ioc:Adonis/Core/Request";
import validator from "validator";
import Article from "App/Models/Elven/Article";
import {EL_ErrorCollector} from "App/Common/Elven/_TOOLS/EL_Errors"
import {
  EL_IError_VALIDATION_FORBIDDEN,
  EL_IError_VALIDATION_MINMAX,
  EL_IError_VALIDATION_MUSTBE
} from "App/Common/Elven/_TOOLS/EL_Interfaces";

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
      if (title.length > 124) { // todo: оформить в соответствии с EL_Errors
        errorCollector.addError('VALIDATION', 400, '«title» must be less than 124 symbols.')
      }
    }
    if (!content) {
      errorCollector.addError('VALIDATION', 400, '«content» can not be empty.')
    }
    is_published = is_published.toString()
    if (!validator.isBoolean(is_published)) { // todo: вместо проверки на bool, проверять isPublished === true || isPublished === false
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
        const minmax: EL_IError_VALIDATION_MINMAX = {
          errorCode: 'E_VALIDATION_MINMAX',
          issuer: 'title',
          max: '124'
        }
        errorCollector.addError(minmax)
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
      const mustbe: EL_IError_VALIDATION_MUSTBE = {
        errorCode: 'E_VALIDATION_MUSTBE',
        issuer: 'show',
        available: ['drafts', 'all']
      }
      errorCollector.addError(mustbe)
    } else {
      if ((show === 'drafts' || show === 'all') && !isAdmin) {
        const forbidden: EL_IError_VALIDATION_FORBIDDEN = {
          errorCode: 'E_VALIDATION_FORBIDDEN',
          issuer: ['drafts', 'all']
        }
        errorCollector.addError(forbidden)
      }
    }
    let by = request.input('by', 'published')
    by = by.toLowerCase()
    if (by !== 'created' && by !== 'published' && by !== 'updated') {
      const mustbe: EL_IError_VALIDATION_MUSTBE = {
        errorCode: 'E_VALIDATION_MUSTBE',
        issuer: 'by',
        available: ['created', 'published', 'updated']
      }
      errorCollector.addError(mustbe)
    } else {
      if (by === 'created' || by === 'updated') {
        if (!isAdmin) {
          const forbidden: EL_IError_VALIDATION_FORBIDDEN = {
            errorCode: 'E_VALIDATION_FORBIDDEN',
            issuer: 'by',
          }
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
      const mustbe: EL_IError_VALIDATION_MUSTBE = {
        errorCode: 'E_VALIDATION_MUSTBE',
        issuer: 'start',
        available: ['newest', 'oldest']
      }
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
      const mustbe: EL_IError_VALIDATION_MUSTBE = {
        errorCode: 'E_VALIDATION_MUSTBE',
        issuer: 'preview',
        available: ['true', 'false']
      }
      errorCollector.addError(mustbe)
    } else if (preview === 'true') {
      preview = true
    } else if (preview === 'false') {
      preview = false
    }
    let page = request.input('page', 1)
    if (page < 1) {
      const minmax: EL_IError_VALIDATION_MINMAX = {
        errorCode: 'E_VALIDATION_MINMAX',
        issuer: 'page',
        min: '1'
      }
      errorCollector.addError(minmax)
    }

    if (errorCollector.hasErrors()) {
      return Promise.reject(errorCollector.getErrors())
    }
    return Promise.resolve({page: page, show: show, by: by, start: start, preview: preview})
  }

}
