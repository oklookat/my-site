import {HttpContextContract} from '@ioc:Adonis/Core/HttpContext'
import {RequestContract} from '@ioc:Adonis/Core/Request'
import AuthMaster from "App/Common/Elven/Auth/AuthMaster"
import ErrorConstructors from "App/Common/Elven/_TOOLS/ErrorConstructors"
import UserValidator from "App/Common/Elven/_VALIDATORS/UserValidator";

export default class AuthController {

  public async login({request, response}: HttpContextContract) {
    const {type} = request.all()
    if (!type) {
      return response.badRequest()
    }
    let isAdminLogin = false
    if (type === 'admin') {
      isAdminLogin = true
    }
    try {
      const token = await this.authRiver(request, isAdminLogin)
      return response.status(200).send({token: token})
    } catch (error) {
      if (!error.type) {
        const err = ErrorConstructors.publicError('Произошла странная ошибка.')
        return response.status(500).send(err)
      }
      const type = error.type
      const isWrongLogin = type === 'WRONG_PASSWORD' || type === 'USER_NOT_FOUND' || type === 'VALIDATION_ERROR'
      if (isWrongLogin) {
        const err = ErrorConstructors.publicError('Неверный логин или пароль.')
        return response.forbidden(err)
      }
      const err = ErrorConstructors.publicError('Применена магия вне Хогвартса, или сервер сошел с ума. Попробуйте очистить данные сайта.')
      return response.forbidden(err)
    }
  }

  public async logout({request, response}: HttpContextContract) {
    await AuthMaster.logout(request)
      .catch(() =>{
        const err = ErrorConstructors.publicError('Применена магия вне Хогвартса, или сервер сошел с ума. Попробуйте очистить данные сайта.')
        return response.forbidden(err)
      })
    return response.status(200).send('')
  }

  private async authRiver(request: RequestContract, adminLogin: boolean) {
    let validated
    await UserValidator.validateCredentials(request)
      .then(data => {
        validated = data
      })
      .catch(error => {
        return Promise.reject(error)
      })
    return await AuthMaster.login(validated.username, validated.password, adminLogin, request)
      .then(token => {
        return Promise.resolve(token)
      })
      .catch(error => {
        return Promise.reject(error)
      })
  }


}
