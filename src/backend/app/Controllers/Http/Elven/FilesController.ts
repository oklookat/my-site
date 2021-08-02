import Env from '@ioc:Adonis/Core/Env'
import {HttpContextContract} from "@ioc:Adonis/Core/HttpContext"
import Application from "@ioc:Adonis/Core/Application"
import {cuid} from '@ioc:Adonis/Core/Helpers'
import File from "App/Models/Elven/File"
import EL_Files from "App/Common/Elven/_TOOLS/EL_Files"
import FileValidator from "App/Common/Elven/_VALIDATORS/FileValidator"

const pageSize = Env.get('PAGINATION_SIZE') // default: 16
const _uploadsPath = Application.tmpPath(`uploads`)
const _tempPath = Application.tmpPath(`uploads/temp`)

export default class FilesController {

  // GET url/
  // params:
  // page = number
  // start = newest (DESC), oldest (ASC)
  public async index(ctx: HttpContextContract) {
    let validatedParams
    try {
      validatedParams = FileValidator.requestParams(ctx.request)
    } catch (errors) {
      return ctx.response.status(400).send(errors)
    }
    const files = await File.query().orderBy('created_at', validatedParams.start).paginate(validatedParams.page, pageSize)
    return ctx.response.status(200).send(files)
  }


  // POST url/
  public async store(ctx: HttpContextContract) {
    const file = ctx.request.file('file', {
      size: '128mb'
    })
    console.log(file?.size)
    if (!file) {
      return ctx.response.status(400).send('Request does not contains file.')
    }
    if (!file.isValid) {
      return ctx.response.status(400).send('With file something is bad. Maybe his broken?')
    }
    const user = ctx['user']
    const extension = file.extname
    const originalFileName = file.clientName
    const newFileName = `${cuid()}.${extension}`
    const _tempPathOfFile = `${_tempPath}/${newFileName}`
    await file.move(_tempPath, {
      name: newFileName
    })
    const fileHash = await EL_Files.getHash(_tempPathOfFile)
    const foundFile = await File.findBy('hash', fileHash)
    if (foundFile) {
      // if file exists
      await EL_Files.deleteFile(_tempPathOfFile) // delete file from temp folder
        .catch(() => {
          return ctx.response.status(500).send('Error while uploading file.')
        })
      return ctx.response.status(200).send(foundFile)
    }
    const foldersGen = EL_Files.generateDirectoriesByHash(fileHash)
    const _newFolder = `${_uploadsPath}/${foldersGen}`
    let isFolderExists
    try {
      isFolderExists = await EL_Files.directoryExists(_newFolder)
    } catch (error) {
      return ctx.response.status(500).send('Error while uploading file.')
    }
    if (!isFolderExists) {
      const error = await EL_Files.createDirectory(_newFolder)
      if (error) {
        return ctx.response.status(500).send('Error while uploading file.')
      }
    }
    try {
      await EL_Files.move(`${_tempPath}/${newFileName}`, `${_uploadsPath}/${foldersGen}/${newFileName}`)
    } catch (error) {
      return ctx.response.status(500).send('Error while uploading file.')
    }
    const newFile = new File()
    newFile.hash = fileHash
    newFile.path = `${foldersGen}/${newFileName}`
    newFile.original_name = originalFileName
    newFile.name = newFileName
    newFile.extension = extension
    newFile.size = file.size.toString()
    const savedFile = await user.related('files').save(newFile)
    return ctx.response.status(200).send(savedFile)
  }


  // DELETE url/:id
  public async destroy(ctx: HttpContextContract) {
    const file = await File.find(ctx.params.id)
    if (!file) {
      return ctx.response.notFound('File not found.')
    }
    try {
      await EL_Files.deleteFile(`${_uploadsPath}/${file.path}`)
      EL_Files.deleteEmptyDirsRecursive(_uploadsPath, file.path)
      await file.delete()
      return ctx.response.status(200).send('File deleted.')
    } catch (error) {
      console.log(error)
      return ctx.response.internalServerError('Error while deleting file.')
    }
  }

}
