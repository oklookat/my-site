import {HttpContextContract} from '@ioc:Adonis/Core/HttpContext'
import {RequestContract} from "@ioc:Adonis/Core/Request"
import EL_Auth from "App/Common/Elven/_TOOLS/EL_Auth"
import User from "App/Models/Elven/User"
import Token from "App/Models/Elven/Token"
import {E_AUTH_FORBIDDEN} from "App/Common/Elven/_ERRORS/EL_Errors";
import {EL_ErrorCollector} from "App/Common/Elven/_ERRORS/EL_ErrorCollector";

const readOnlyMethods = ['GET', 'OPTIONS', 'HEAD']
const access = { // if admin, access categories will be ignored
  readOnly: 'readOnly',
}

export default class Perm {

  public async handle(ctx: HttpContextContract, next: () => Promise<void>, perm: string[] = ['null']) {
    let isAllow
    try {
      const firstPerm = perm[0]
      const data = await Perm.permChecker(ctx.request, firstPerm)
      isAllow = data.access
      if(data.user && data.token){
        ctx['user'] = data.user // after this middleware you be have ctx.user with user instance (if user logged in, of course)
        await EL_Auth.tokenWriteLastAgents(ctx.request, data.token)
      }
    } catch (error) {
      // all errors in this middleware be here. For debug paste console.log(error).
      isAllow = false
    }
    if (!isAllow) {
      const forbidden = new E_AUTH_FORBIDDEN(['AUTH'])
      return ctx.response.status(403).send(EL_ErrorCollector.singleError(forbidden))
    }
    return await next()
  }

  private static async permChecker(request: RequestContract, perm: string): Promise<{access: boolean, user?: User | null, token?: Token | null}> {
    try{
      const data = await EL_Auth.isAdmin(request)
      if (data.isAdmin) {
        return Promise.resolve({access: data.isAdmin, user: data.user, token: data.token})
      }
    } catch (error){

    }
    if (perm === access.readOnly) {
      const method = request.method()
      const isReadOnly = await Perm.readOnly(method)
      return Promise.resolve({access: isReadOnly, user: null})
    }
    return Promise.resolve({access: false, user: null})
  }

  private static async readOnly(method: string): Promise<boolean> {
    if (readOnlyMethods.includes(method)) {
      return Promise.resolve(true)
    } else {
      return Promise.resolve(false)
    }
  }

}
