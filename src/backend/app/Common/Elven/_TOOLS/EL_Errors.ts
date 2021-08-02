import {
  EL_IError_UNKNOWN,
  EL_IError_VALIDATION_MUSTBE,
  EL_IError_VALIDATION_FORBIDDEN,
  EL_IError_AUTH_INCORRECT, EL_IError_VALIDATION_MINMAX, EL_IError_CUSTOM,
} from "App/Common/Elven/_TOOLS/EL_Interfaces"


export class EL_ErrorCollector {
  private errorsArray: object []

  public hasErrors() {
    if (this.errorsArray) {
      return this.errorsArray.length > 0
    } else {
      return false
    }
  }

  public addError(error: EL_IError_AUTH_INCORRECT | EL_IError_VALIDATION_MUSTBE | EL_IError_VALIDATION_FORBIDDEN | EL_IError_VALIDATION_MINMAX | EL_IError_UNKNOWN | EL_IError_CUSTOM): boolean {
    // CRITICAL error ignores other errors and send only self
    this.errorsArray.push(error)
    return true
  }

  public static singleError(error: EL_IError_AUTH_INCORRECT | EL_IError_VALIDATION_MUSTBE | EL_IError_VALIDATION_FORBIDDEN | EL_IError_VALIDATION_MINMAX | EL_IError_UNKNOWN | EL_IError_CUSTOM): object {
    return {error: error}
  }

  public getErrors(): { errors: object[] } {
    if (!this.hasErrors()) {
      const unknownError: EL_IError_UNKNOWN = {
        errorCode: 'E_UNKNOWN',
        message: `This is an auto-generated error. Something went wrong, but the admin didnt add it here. This is strange.`
      }
      this.addError(unknownError)
    }
    return {errors: [this.errorsArray]}
  }
}
