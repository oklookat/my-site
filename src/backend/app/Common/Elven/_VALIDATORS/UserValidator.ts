import validator from "validator"
import {EL_ErrorCollector} from "App/Common/Elven/_ERRORS/EL_ErrorCollector"
import {E_VALIDATION_ALLOWED, E_VALIDATION_EMPTY, E_VALIDATION_MINMAX} from "App/Common/Elven/_ERRORS/EL_Errors";

export default class UserValidator {

  public static validateCredentials(username: string, password: string) {
    const errorCollector = new EL_ErrorCollector()
    const isUsernameEmpty = validator.isEmpty(username)
    const isPasswordEmpty = validator.isEmpty(password)
    if (isUsernameEmpty || isPasswordEmpty) {
      let failed = 0
      if (isUsernameEmpty) {
        const empty = new E_VALIDATION_EMPTY(['username'])
        errorCollector.addError(empty)
        failed++
      } if (isPasswordEmpty) {
        const empty = new E_VALIDATION_EMPTY(['password'])
        errorCollector.addError(empty)
        failed++
      }
      if(failed == 2){
        throw errorCollector.getErrors()
      }
    }
    let isUsername = validator.isAlphanumeric(username)
    if (!isUsername) {
      const allowed = new E_VALIDATION_ALLOWED(['username'], ['alphanumeric'])
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
      throw errorCollector.getErrors()
    }
    return {username: username, password: password}
  }
}
