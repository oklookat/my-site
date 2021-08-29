import Env from '@ioc:Adonis/Core/Env'
import {HttpContextContract} from "@ioc:Adonis/Core/HttpContext"
import {cuid} from '@ioc:Adonis/Core/Helpers'
import File from "App/Models/Elven/File"
import EL_Files from "App/Common/Elven/_TOOLS/EL_Files"
import FileValidator from "App/Common/Elven/_VALIDATORS/FileValidator"
import {EL_ErrorCollector} from "App/Common/Elven/_ERRORS/EL_ErrorCollector";
import {E_UNKNOWN, E_VALIDATION_EMPTY, E_VALIDATION_INVALID} from "App/Common/Elven/_ERRORS/EL_Errors";

const pageSize = Env.get('PAGINATION_SIZE') // default: 16
const _uploadsPath = EL_Files._uploadsPath
const _tempPath = EL_Files._tempPath

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
    // validation
    const file = ctx.request.file('file', {
      size: '256mb'
    })
    if (!file) {
      const empty = new E_VALIDATION_EMPTY(['files'])
      return ctx.response.status(400).send(EL_ErrorCollector.singleError(empty))
    }
    if (!file.isValid) {
      const invalid = new E_VALIDATION_INVALID(['files'], 'File is invalid. Maybe broken?')
      return ctx.response.status(400).send(invalid)
    }
    const user = ctx.user
    if(!user){
      return FilesController.send505(ctx, 'Error while saving file.')
    }
    // create file
    const extension = file.extname
    const originalFileName = file.clientName
    const newFileName = `${cuid()}.${extension}`
    await file.move(_tempPath, {
      name: newFileName
    })
    const _tempPathOfFile = `${_tempPath}/${newFileName}`
    const fileHash = await EL_Files.getHash(_tempPathOfFile)
    const foundFile = await File.findBy('hash', fileHash)
    if (foundFile) {
      // if file exists
      EL_Files.deleteFile(_tempPathOfFile) // delete file from temp folder
        .catch(() => {})
      return ctx.response.status(200).send(foundFile)
    }
    const foldersGen = EL_Files.generateDirectoriesByHash(fileHash)
    const _newFolder = `${_uploadsPath}/${foldersGen}`
    let isFolderExists
    try {
      isFolderExists = await EL_Files.directoryExists(_newFolder)
    } catch (error) {
      await FilesController.send505(ctx)
    }
    if (!isFolderExists) {
      try {
        await EL_Files.createDirectory(_newFolder)
      } catch (error){
        await FilesController.send505(ctx)
      }
    }
    try {
      await EL_Files.move(`${_tempPath}/${newFileName}`, `${_uploadsPath}/${foldersGen}/${newFileName}`)
    } catch (error) {
      await FilesController.send505(ctx)
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
      await EL_Files.deleteFile(`${_uploadsPath}/${file.path}`).catch(() =>{})
      EL_Files.deleteEmptyDirsRecursive(_uploadsPath, file.path)
      await file.delete()
      return ctx.response.status(200).send('File deleted.')
    } catch (error) {
      return FilesController.send505(ctx, 'Error while deleting file.')
    }
  }


  private static async send505(ctx: HttpContextContract, message = 'Error while uploading file.'){
    const unknown = new E_UNKNOWN(['files'], message)
    return ctx.response.status(500).send(EL_ErrorCollector.singleError(unknown))
  }

}
