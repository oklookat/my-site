export default class EL_Errors {
  public static async privateError(type: string, message: string){
    return {type: type, message: message}
  }

  public static publicError(message: string){
    return {error: message}
  }

  public static async publicSuccess(message: string){
    return {success: message}
  }
}
