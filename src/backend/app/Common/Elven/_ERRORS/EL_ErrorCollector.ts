import {
  E_AUTH_FORBIDDEN,
  E_AUTH_INCORRECT,
  E_CUSTOM,
  E_NOTFOUND,
  E_UNKNOWN,
  E_VALIDATION_MINMAX
} from "App/Common/Elven/_ERRORS/EL_Errors";

export class EL_ErrorCollector {
  private errorsArray: object [] = []

  public hasErrors() {
    if (this.errorsArray) {
      return this.errorsArray.length > 0
    } else {
      return false
    }
  }

  public addError(error: E_UNKNOWN | E_CUSTOM | E_AUTH_INCORRECT
    | E_AUTH_FORBIDDEN | E_NOTFOUND | E_VALIDATION_MINMAX): boolean {
    this.errorsArray.push(error)
    return true
  }

  public static singleError(error: E_UNKNOWN | E_CUSTOM | E_AUTH_INCORRECT
    | E_AUTH_FORBIDDEN | E_NOTFOUND | E_VALIDATION_MINMAX): object {
    return {errors: [error]}
  }

  public getErrors(): { errors: object[] } {
    if (!this.hasErrors()) {
      const unknown = new E_UNKNOWN(['SYSTEM'], 'Auto-generated error. Seems like developer forgot add error. Hmm.')
      this.addError(unknown)
    }
    return {errors: this.errorsArray}
  }
}
