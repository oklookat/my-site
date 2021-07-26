import {RequestContract} from "@ioc:Adonis/Core/Request";
import validator from "validator";
import Article from "App/Models/Elven/Article";
import ErrorConstructors from "App/Common/Elven/_TOOLS/ErrorConstructors";

export default class ArticleValidator{

  public static async whenCreate(request: RequestContract) {
    try{
      let {is_published, thumbnail, title, content} = request.all()
      if (!is_published) {
        is_published = false
      }
      if(!thumbnail){
        thumbnail = null
      }
      if (!title) {
        title = 'Без названия'
      } else {
        title = title.toString()
        if(title.length > 124){
          const err = await ErrorConstructors.privateError('VALIDATION_ERROR', 'Заголовок должен быть меньше 124 символов.')
          return Promise.reject(err)
        }
      }
      if (!content) {
        const err = await ErrorConstructors.privateError('VALIDATION_ERROR', 'Контент не может быть пустым.')
        return Promise.reject(err)
      }
      is_published = is_published.toString()
      if (!validator.isBoolean(is_published)) {
        const err = await ErrorConstructors.privateError('VALIDATION_ERROR', 'isPublished должен иметь тип bool.')
        return Promise.reject(err)
      }
      if(typeof content !== 'object' && content !== null){
        const err = await ErrorConstructors.privateError('VALIDATION_ERROR', 'content не является объектом.')
        return Promise.reject(err)
      }
      return Promise.resolve({is_published: is_published, thumbnail: thumbnail, title: title, content: content})
    }
    catch (error){
      console.log(error)
      const err = await ErrorConstructors.privateError('VALIDATION_ERROR', error)
      return Promise.reject(err)
    }
  }

  public static async whenUpdate(request: RequestContract, foundArticle: Article){
    let {is_published, thumbnail, title, content} = request.all()
    if(is_published !== true && is_published !== false){
      is_published = foundArticle.is_published
    }
    if(!thumbnail){
      thumbnail = foundArticle.thumbnail
    }
    if(!title){
      title = foundArticle.title
    } else{
      if(title.length > 124){
        const err = await ErrorConstructors.privateError('VALIDATION_ERROR', 'Заголовок должен быть меньше 124 символов.')
        return Promise.reject(err)
      }
    }
    if(!content){
      content = foundArticle.content
    }
    return Promise.resolve({is_published: is_published, thumbnail: thumbnail, title: title, content: content})
  }

}
