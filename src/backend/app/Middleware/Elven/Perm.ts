import {HttpContextContract} from '@ioc:Adonis/Core/HttpContext'
import {RequestContract} from "@ioc:Adonis/Core/Request";
import AuthMaster from "App/Common/Elven/Auth/AuthMaster";
import User from "App/Models/Elven/User";
import Token from "App/Models/Elven/Token";
import EL_Errors from "App/Common/Elven/_TOOLS/EL_Errors";

const readOnlyMethods = ['GET', 'OPTIONS', 'HEAD']
const access = { // если админ, то категории доступа игнорируются
  readOnly: 'readOnly',
}

export default class Perm {

  public async handle(ctx: HttpContextContract, next: () => Promise<void>, perm: string[] = ['null']) {
    let isAllow
    try {
      const firstPerm = perm[0]
      const data = await this.permChecker(ctx.request, firstPerm)
      isAllow = data.access
      if(data.user && data.token){
        ctx['user'] = data.user // по итогу после middleware будет доступен ctx.user с userInstance
        await AuthMaster.tokenWriteLastAgents(ctx.request, data.token)
      }
    } catch (error) {
      // по итогу все ошибки в этом middleware приходят сюда, в случае дебага поставить console.log(error)
      isAllow = false
    }
    if (!isAllow) {
      const err = await EL_Errors.publicError('Доступ запрещен.')
      return ctx.response.forbidden(err)
    }
    return await next()
  }

  private async permChecker(request: RequestContract, perm: string): Promise<{access: boolean, user?: User | null, token?: Token | null}> {
    try{
      const data = await AuthMaster.isAdmin(request)
      if (data.isAdmin) {
        return Promise.resolve({access: data.isAdmin, user: data.user, token: data.token})
      }
    } catch (error){

    }
    if (perm === access.readOnly) {
      const method = request.method()
      const isReadOnly = await this.readOnly(method)
      return Promise.resolve({access: isReadOnly, user: null})
    }
    return Promise.resolve({access: false, user: null})
  }

  private async readOnly(method: string): Promise<boolean> {
    if (readOnlyMethods.includes(method)) {
      return Promise.resolve(true)
    } else {
      return Promise.resolve(false)
    }
  }

}
