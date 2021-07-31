export class EL_Errors {

  public static async privateError(code: number, type: string, message: string) {
    return {httpCode: code, errorType: type, message: message}
  }

  public static publicError(message: string) {
    return {error: message}
  }

  public static async publicSuccess(message: string) {
    return {success: message}
  }
}


export class EL_ErrorCollector {
  private errorsArray: object []
  private isCriticalTriggered = false

  public hasErrors() {
    if (this.errorsArray) {
      return this.errorsArray.length > 0
    } else {
      return false
    }
  }

  public addError(errorType: 'GENERIC' | 'VALIDATION' | 'CRITICAL', httpCode: number, message: string): boolean {
    // CRITICAL error ignores other errors and send only self
    if (this.isCriticalTriggered) {
      return false
    } else if (!this.isCriticalTriggered && errorType === 'CRITICAL') {
      this.isCriticalTriggered = true
    }
    this.errorsArray.push({
      type: 'ERROR',
      subType: errorType,
      httpCode: httpCode,
      message: message
    })
    return true
  }

  public getErrors(): { errors: object[] } {
    return {errors: this.errorsArray}
  }
}
