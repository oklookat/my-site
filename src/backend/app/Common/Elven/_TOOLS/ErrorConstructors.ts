export default class ErrorConstructors{
  public static async privateError(type: string, message: string): Promise<{ type: string, message: string }> {
    return Promise.resolve({type: type, message: message})
  }

  public static async publicError(message: string): Promise<{ error: string }>{
    return Promise.resolve({error: message})
  }

  public static async publicSuccess(message: string): Promise<{ success: string }>{
    return Promise.resolve({success: message})
  }
}
