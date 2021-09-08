import {HttpContextContract} from '@ioc:Adonis/Core/HttpContext'
import {RequestContract} from "@ioc:Adonis/Core/Request"
import EL_Auth, {IUserAndToken} from "App/Common/Elven/_TOOLS/EL_Auth"
import User from "App/Models/Elven/User"
import Token from "App/Models/Elven/Token"
import {E_AUTH_FORBIDDEN} from "App/Common/Elven/_ERRORS/EL_Errors";
import {EL_ErrorCollector} from "App/Common/Elven/_ERRORS/EL_ErrorCollector";
import Logger from "@ioc:Adonis/Core/Logger";

interface IPerm {
  readonlyMethods: string[]
  accessTypes: string[]
}

const perms: IPerm = {
  readonlyMethods: ['GET', 'OPTIONS', 'HEAD'],
  accessTypes: ['readOnly', 'adminOnly']
}

interface IAccess {
  access: boolean,
  user: User | null,
  token: Token | null
}

export default class Perm {

  public async handle(ctx: HttpContextContract, next: () => Promise<void>, perm: string[]) {
    if (perm.length < 1 || !perms.accessTypes.includes(perm[0])) {
      Logger.fatal("[perm middleware] wrong access type in args")
    }
    let isAccessGranted: boolean
    const accessType = perm[0]
    try {
      const data = await Perm.check(ctx.request, accessType)
      isAccessGranted = data.access
      if (data.user && data.token) {
        ctx.user = data.user // after this middleware you be have ctx.user with user instance (if user logged in, of course)
        data.token = await EL_Auth.tokenWriteLastAgents(ctx.request, data.token)
      }
    } catch (error) {
      // all errors in this middleware be here. For debug paste console.log(error).
      isAccessGranted = false
    }
    if (!isAccessGranted) {
      const forbidden = new E_AUTH_FORBIDDEN(['AUTH'])
      return ctx.response.status(403).send(EL_ErrorCollector.singleError(forbidden))
    }
    return await next()
  }

  private static async check(request: RequestContract, perm: string): Promise<IAccess> {
    let instances: IUserAndToken | null = null
    try {
      instances = await EL_Auth.getUserAndTokenByRequest(request)
    } catch (error) {
    }
    let access: IAccess = {
      access: false,
      user: null,
      token: null
    }
    if(instances){
      access.user = instances.user
      access.token = instances.token
      switch (instances.user.role){
        case 'admin':
          access.access = true
          return Promise.resolve(access)
        default:
          break
      }
    }
    switch (perm) {
      case perms.accessTypes[0]: // readOnly
        const method = request.method()
        access.access = Perm.isReadOnlyMethod(method)
        break
      case perms.accessTypes[1]: // adminOnly
        access.access = false // access denied, because recently we checked is user an admin
        break
      default:
        access.access = false
    }
    return Promise.resolve(access)
  }

  private static isReadOnlyMethod(method: string): boolean {
    return perms.readonlyMethods.includes(method)
  }

}
