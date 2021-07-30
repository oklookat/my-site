import Env from '@ioc:Adonis/Core/Env'
import {HttpContextContract} from "@ioc:Adonis/Core/HttpContext"
import Application from "@ioc:Adonis/Core/Application"
import {cuid} from '@ioc:Adonis/Core/Helpers'
import File from "App/Models/Elven/File"
import EL_Errors from "App/Common/Elven/_TOOLS/EL_Errors"
import EL_Files from "App/Common/Elven/_TOOLS/EL_Files"
import path from 'path'

const pageSize = Env.get('PAGINATION_SIZE') // default: 16
const _uploadsPath = Application.tmpPath(`uploads`)
const _tempPath = Application.tmpPath(`uploads/temp`)

export default class FilesController {

  // GET url/
  // params:
  // page = number
  // start = newest (DESC), oldest (ASC)
  public async index(ctx: HttpContextContract) {
    // VALIDATION START //
    let start = ctx.request.input('start', 'newest')
    start = start.toLowerCase()
    if (start !== 'newest' && start !== 'oldest') {
      return ctx.response.status(400).send(await EL_Errors.publicError('«start» must be newest or oldest.'))
    } else {
      if (start === 'newest') {
        start = 'DESC'
      } else if (start === 'oldest') {
        start = 'ASC'
      }
    }
    let page = ctx.request.input('page', 1)
    if(page < 1){
      return ctx.response.status(400).send(await EL_Errors.publicError('«page» cannot be less than one.'))
    }
    // VALIDATION END //
    const files = await File.query().orderBy('created_at', start).paginate(page, pageSize)
    return ctx.response.status(200).send(files)
  }


  // POST url/
  public async store(ctx: HttpContextContract) {
    const file = ctx.request.file('file', {
      size: '128mb'
    })
    if (!file) {
      return ctx.response.status(400).send(EL_Errors.publicError('Request does not contains file.'))
    }
    if (!file.isValid) {
      return ctx.response.status(400).send(EL_Errors.publicError('With file something is bad. Maybe his broken?'))
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
          return ctx.response.status(500).send(EL_Errors.publicError('Error while uploading file.'))
        })
      return ctx.response.status(200).send(foundFile)
    }
    const foldersGen = EL_Files.generateDirectoriesByHash(fileHash)
    const _newFolder = `${_uploadsPath}/${foldersGen}`
    let isFolderExists
    try {
      isFolderExists = await EL_Files.directoryExists(_newFolder)
    } catch (error) {
      return ctx.response.status(500).send(EL_Errors.publicError('Error while uploading file.'))
    }
    if (!isFolderExists) {
      const error = await EL_Files.createDirectory(_newFolder)
      if (error) {
        return ctx.response.status(500).send(EL_Errors.publicError('Error while uploading file.'))
      }
    }
    try {
      await EL_Files.move(`${_tempPath}/${newFileName}`, `${_uploadsPath}/${foldersGen}/${newFileName}`)
    } catch (error){
      return ctx.response.status(500).send(EL_Errors.publicError('Error while uploading file.'))
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
      return ctx.response.notFound(await EL_Errors.publicError('File not found.'))
    }
    try {
      await EL_Files.deleteFile(`${_uploadsPath}/${file.path}`)
      await EL_Files.deleteEmptyDirsRecursive(_uploadsPath, file.path)
      await file.delete()
      return ctx.response.status(200).send(await EL_Errors.publicError('File deleted.'))
    } catch (error) {
      console.log(error)
      return ctx.response.internalServerError(await EL_Errors.publicError('Error while deleting file.'))
    }
  }

}
