import {RequestContract} from "@ioc:Adonis/Core/Request";
import {EL_ErrorCollector} from "App/Common/Elven/_TOOLS/EL_Errors"

export default class FileValidator {

  public static async requestParams(request: RequestContract) {
    const errorCollector = new EL_ErrorCollector()
    let start = request.input('start', 'newest')
    start = start.toLowerCase()
    if (start !== 'newest' && start !== 'oldest') {
      errorCollector.addError('VALIDATION', 400, '«start» must be newest or oldest.')
    } else {
      if (start === 'newest') {
        start = 'DESC'
      } else if (start === 'oldest') {
        start = 'ASC'
      }
    }
    let page = request.input('page', 1)
    if (page < 1) {
      errorCollector.addError('VALIDATION', 400, '«page» cannot be less than one.')
    }
    if (errorCollector.hasErrors()) {
      return Promise.reject(errorCollector.getErrors())
    }
    return Promise.resolve({page: page, start: start})
  }

}
