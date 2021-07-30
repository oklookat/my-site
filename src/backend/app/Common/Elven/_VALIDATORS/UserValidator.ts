import {RequestContract} from "@ioc:Adonis/Core/Request"
import validator from "validator"
import User from "App/Models/Elven/User"
import EL_Errors from "App/Common/Elven/_TOOLS/EL_Errors";

export default class UserValidator {
  public static async validateCredentials(request: RequestContract) {
    const {username, password} = request.all()
    if (!username) {
      const err = await EL_Errors.privateError('VALIDATION_ERROR', 'Имя пользователя не может быть пустым.')
      return Promise.reject(err)
    }
    if (!password) {
      const err = await EL_Errors.privateError('VALIDATION_ERROR', 'Пароль не может быть пустым.')
      return Promise.reject(err)
    }
    const isPass = validator.isLength(password, {min: 8, max: 64})
    if (!isPass) {
      const err = await EL_Errors.privateError('VALIDATION_ERROR', 'Пароль должен быть больше 8 и меньше 64 символов.')
      return Promise.reject(err)
    }
    let isUsername = validator.isLength(username, {min: 4, max: 24})
    if (!isUsername) {
      const err = await EL_Errors.privateError('VALIDATION_ERROR', 'Имя пользователя должно быть больше 4 и меньше 24 символов.')
      return Promise.reject(err)
    }
    isUsername = validator.isAlphanumeric(username)
    if (!isUsername) {
      const err = await EL_Errors.privateError('VALIDATION_ERROR', 'Имя пользователя должно быть без странных символов и только на английском языке.')
      return Promise.reject(err)
    }
    return Promise.resolve({username: username, password: password})
  }

  public static async validateReg(user: User) {
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
