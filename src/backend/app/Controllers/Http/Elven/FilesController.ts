import {HttpContextContract} from "@ioc:Adonis/Core/HttpContext"
import Application from "@ioc:Adonis/Core/Application"
import { cuid } from '@ioc:Adonis/Core/Helpers'


export default class FilesController {

  // GET url/
  public async index(ctx: HttpContextContract) {

  }

  // GET url/:id
  public async show(ctx: HttpContextContract) {
    return ctx.response.attachment(
      Application.tmpPath('uploads', ctx.params.filename)
    )
  }

  // POST url/
  public async store(ctx: HttpContextContract) {
    const files = ctx.request.files('files', {
      size: '128mb'
    })
    if (!files) {
      return
    }


    for (let file of files) {
      const fileName = `${cuid()}.${file.extname}`
      await file.move(Application.tmpPath('uploads'), {
        name: fileName
      })
    }

  }

  // PUT OR PATCH url/:id
  public async update(ctx: HttpContextContract) {

  }

  // DELETE url/:id
  public async destroy(ctx: HttpContextContract) {

  }

}
