import {HttpContextContract} from '@ioc:Adonis/Core/HttpContext'
import {RequestContract} from '@ioc:Adonis/Core/Request'
import AuthMaster from "App/Common/Elven/Auth/AuthMaster"
import UserValidator from "App/Common/Elven/_VALIDATORS/UserValidator"

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
        return response.status(500).send('Произошла странная ошибка.')
      }
      const type = error.type
      const isWrongLogin = type === 'WRONG_PASSWORD' || type === 'USER_NOT_FOUND' || type === 'VALIDATION_ERROR'
      if (isWrongLogin) {
        return response.status(401).send('Неверный логин или пароль.')
      }
      return response.status(500).send('Применена магия вне Хогвартса, или сервер сошел с ума. Попробуйте очистить данные сайта.')
    }
  }

  public async logout({request, response}: HttpContextContract) {
    await AuthMaster.logout(request)
      .catch(() =>{
        return response.status(500).send('Применена магия вне Хогвартса, или сервер сошел с ума. Попробуйте очистить данные сайта.')
      })
    return response.status(200).send('Успешный выход.')
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
