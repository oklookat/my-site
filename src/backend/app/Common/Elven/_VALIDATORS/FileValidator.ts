import {RequestContract} from "@ioc:Adonis/Core/Request"
import {EL_ErrorCollector} from "App/Common/Elven/_ERRORS/EL_ErrorCollector"
import {E_VALIDATION_MINMAX, E_VALIDATION_MUSTBE} from "App/Common/Elven/_ERRORS/EL_Errors"

export default class FileValidator {

  public static async requestParams(request: RequestContract) {
    const errorCollector = new EL_ErrorCollector()
    let start = request.input('start', 'newest')
    start = start.toLowerCase()
    if (start !== 'newest' && start !== 'oldest') {
      const mustbe = new E_VALIDATION_MUSTBE(['start'], ['newest', 'oldest'])
      errorCollector.addError(mustbe)
    } else {
      if (start === 'newest') {
        start = 'DESC'
      } else if (start === 'oldest') {
        start = 'ASC'
      }
    }
    let page = request.input('page', 1)
    if (page < 1) {
      const minmax = new E_VALIDATION_MINMAX(['page'], 1)
      errorCollector.addError(minmax)
    }
    if (errorCollector.hasErrors()) {
      return Promise.reject(errorCollector.getErrors())
    }
    return Promise.resolve({page: page, start: start})
  }

}
