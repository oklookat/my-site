declare module '@ioc:Adonis/Core/HttpContext' {

  import User from "App/Models/Elven/User";

  interface HttpContextContract {
    user: User | null
  }
}
