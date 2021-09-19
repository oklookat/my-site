import {HttpContextContract} from '@ioc:Adonis/Core/HttpContext'
import EL_Auth from "App/Common/Elven/_TOOLS/EL_Auth"
import UserValidator from "App/Common/Elven/_VALIDATORS/UserValidator"
import {EL_ErrorCollector} from "App/Common/Elven/_ERRORS/EL_ErrorCollector";
import {
  E_AUTH_INCORRECT,
  E_UNKNOWN,
  E_VALIDATION_ALLOWED,
} from "App/Common/Elven/_ERRORS/EL_Errors";


export default class AuthController {

  public async login({request, response}: HttpContextContract) {
    let {username, password, type} = request.all()
    if(!type || type !== 'cookie' && type !== 'direct'){
      const allowed = new E_VALIDATION_ALLOWED(['type'], ['cookie', 'direct'])
      return response.status(400).send(EL_ErrorCollector.singleError(allowed))
    }
    try {
      const data = UserValidator.validateCredentials(username, password)
      username = data.username
      password = data.password
    } catch (errors) {
      return response.status(400).send(errors)
    }
    try {
      const token = await EL_Auth.login(username, password, request)
      switch (type){
        case 'cookie':
          return response.status(200).cookie('token', token)
        case 'direct':
          return response.status(200).send({token: token})
      }
    } catch (error) {
      if (error === 'PIPE_TOKEN_SAVING_ERROR') {
        const unknown = new E_UNKNOWN(['AUTH'], '')
        return response.status(500).send(EL_ErrorCollector.singleError(unknown))
      } else {
        const wrong = new E_AUTH_INCORRECT(['AUTH'])
        return response.status(403).send(EL_ErrorCollector.singleError(wrong))
      }
    }
  }

  public async logout({request, response}: HttpContextContract) {
    try {
      await EL_Auth.logout(request)
    } catch (err) {
    }
    return response.status(200).send('')
  }

  public async check({response}: HttpContextContract){
    return response.status(200).send('')
  }

}
