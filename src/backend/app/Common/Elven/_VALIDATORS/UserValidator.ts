import {RequestContract} from "@ioc:Adonis/Core/Request"
import validator from "validator"
import User from "App/Models/Elven/User"
import {EL_ErrorCollector} from "App/Common/Elven/_ERRORS/EL_ErrorCollector"
import {E_VALIDATION_ALLOWED, E_VALIDATION_EMPTY, E_VALIDATION_MINMAX} from "App/Common/Elven/_ERRORS/EL_Errors";

export default class UserValidator {

  public static async validateCredentials(request: RequestContract) {
    const errorCollector = new EL_ErrorCollector()
    const {username, password} = request.all()

    const isUsernameEmpty = validator.isEmpty(username)
    const isPasswordEmpty = validator.isEmpty(password)
    if (isUsernameEmpty || isPasswordEmpty) {
      if (isUsernameEmpty) {
        const empty = new E_VALIDATION_EMPTY(['username'])
        errorCollector.addError(empty)
      } if (isPasswordEmpty) {
        const empty = new E_VALIDATION_EMPTY(['password'])
        errorCollector.addError(empty)
      }
    }
    let isUsername = validator.isAlphanumeric(username)
    if (!isUsername) {
      const allowed = new E_VALIDATION_ALLOWED(['username'], ['alphanumeric'], 'Allowed only alphanumeric symbols.')
      errorCollector.addError(allowed)
    }
    isUsername = validator.isLength(username, {min: 4, max: 24})
    if (!isUsername) {
      const minmax = new E_VALIDATION_MINMAX(['username'], 4, 24)
      errorCollector.addError(minmax)
    }
    const isPass = validator.isLength(password, {min: 8, max: 64})
    if (!isPass) {
      const minmax = new E_VALIDATION_MINMAX(['password'], 8, 64)
      errorCollector.addError(minmax)
    }
    if (errorCollector.hasErrors()) {
      return Promise.reject(errorCollector.getErrors())
    }
    return Promise.resolve({username: username, password: password})
  }

  public static async validateReg(user: User) { // todo: перевести на английский, и совместить с validateCredentials(?)
    let isUsername = validator.isLength(user.username, {min: 4, max: 24})
    if (!isUsername) {
      const err = new Error('Имя пользователя должно быть больше 4, и меньше 24 символов.')
      return Promise.reject(err)
    }
    isUsername = validator.isAlphanumeric(user.username)
    if (!isUsername) {
      const err = new Error('Имя пользователя должно быть без странных символов, и только на английском языке.')
      return Promise.reject(err)
    }
    const isPass = validator.isLength(user.password, {min: 8, max: 64})
    if (!isPass) {
      const err = new Error('Пароль должен быть больше 8, и меньше 64 символов.')
      return Promise.reject(err)
    }
    return Promise.resolve(true)
  }

}
