import {HttpContextContract} from '@ioc:Adonis/Core/HttpContext'
import EL_Auth from "App/Common/Elven/_TOOLS/EL_Auth"
import UserValidator from "App/Common/Elven/_VALIDATORS/UserValidator"
import {EL_ErrorCollector} from "App/Common/Elven/_ERRORS/EL_ErrorCollector";
import {E_AUTH_INCORRECT, E_UNKNOWN} from "App/Common/Elven/_ERRORS/EL_Errors";

export default class AuthController {

  public async login({request, response}: HttpContextContract) {
    const {type} = request.all()
    if (!type) {
      return response.status(400).send('')
    }
    let isAdminLogin = type === 'admin'
    let username, password
    try {
      const data = await UserValidator.validateCredentials(request)
      username = data.username
      password = data.password
    } catch (errors){
      return response.status(400).send(errors)
    }
    try {
      const token = await EL_Auth.login(username, password, isAdminLogin, request)
      return response.status(200).send({token: token})
    } catch (error){
      if(error === 'PIPE_TOKEN_SAVING_ERROR'){
        const unknown = new E_UNKNOWN(['AUTH'], 'Server error during auth.')
        return response.status(500).send(EL_ErrorCollector.singleError(unknown))
      } else {
        const wrong = new E_AUTH_INCORRECT(['AUTH'], 'Wrong username or password.')
        return response.status(403).send(EL_ErrorCollector.singleError(wrong))
      }
    }
  }

  public async logout({request, response}: HttpContextContract) {
    await EL_Auth.logout(request)
    return response.status(200).send('')
  }
}
