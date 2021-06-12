import {RequestContract} from "@ioc:Adonis/Core/Request"
import validator from "validator";
import ElvenTools from "App/Common/Elven/_TOOLS/ElvenTools"

class ElvenValidators {

  public static async userValidate(request: RequestContract){
    const {username, password} = request.all()
    if(!username){
      const err = await ElvenTools.errorConstructor('VALIDATION_ERROR', 'Имя пользователя не может быть пустым.')
      return Promise.reject(err)
    }
    if(!password){
      const err = await ElvenTools.errorConstructor('VALIDATION_ERROR', 'Пароль не может быть пустым.')
      return Promise.reject(err)
    }
    const isPass = validator.isLength(password, {min: 8, max: 64})
    if(!isPass){
      const err = await ElvenTools.errorConstructor('VALIDATION_ERROR', 'Пароль должен быть больше 8 и меньше 64 символов.')
      return Promise.reject(err)
    }
    let isUsername = validator.isLength(username, {min: 4, max: 24})
    if(!isUsername){
      const err = await ElvenTools.errorConstructor('VALIDATION_ERROR', 'Имя пользователя должно быть больше 4 и меньше 24 символов.')
      return Promise.reject(err)
    }
    isUsername = validator.isAlphanumeric(username)
    if(!isUsername){
      const err = await ElvenTools.errorConstructor('VALIDATION_ERROR', 'Имя пользователя должно быть без странных символов и только на английском языке.')
      return Promise.reject(err)
    }
    return Promise.resolve({username: username, password: password})
  }

  public static async articleValidate(request: RequestContract){
    let {is_published, thumbnail, title, content} = request.all()
    if(!is_published){
      const err = await ElvenTools.errorConstructor('VALIDATION_ERROR', 'isPublished не может быть пустым.')
      return Promise.reject(err)
    }
    if(!title){
      title = 'Без названия'
    }
    if(!content){
      const err = await ElvenTools.errorConstructor('VALIDATION_ERROR', 'Контент не может быть пустым.')
      return Promise.reject(err)
    }
    if(!validator.isBoolean(is_published)){
      const err = await ElvenTools.errorConstructor('VALIDATION_ERROR', 'isPublished должен иметь тип bool.')
      return Promise.reject(err)
    }
    if(!validator.isLength(content, {min: 1})){
      const err = await ElvenTools.errorConstructor('VALIDATION_ERROR', 'Контент должен быть больше 1 символа.')
      return Promise.reject(err)
    }
    return Promise.resolve({is_published: is_published, thumbnail: thumbnail, title: title, content: content})
  }

  public static async slugValidate(slug: string){
    if(!validator.isAscii(slug)){
      return Promise.resolve(false)
    }
    if(validator.isEmpty(slug, {ignore_whitespace: true})){
      return Promise.resolve(false)
    }
    return Promise.resolve(true)
  }
}

export default ElvenValidators
