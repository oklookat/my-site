import Env from '@ioc:Adonis/Core/Env'
import {HttpContextContract} from "@ioc:Adonis/Core/HttpContext"
import Application from "@ioc:Adonis/Core/Application"
import {cuid} from '@ioc:Adonis/Core/Helpers'
import File from "App/Models/Elven/File"
import md5File from 'md5-file'
import fs from 'fs'
import ErrorConstructors from "App/Common/Elven/_TOOLS/ErrorConstructors"
import mkdirp from 'mkdirp'

const pageSize = Env.get('PAGINATION_SIZE') // default: 16

export default class FilesController {

  // GET url/
  public async index(ctx: HttpContextContract) {
      const files = await File.query().paginate(1, pageSize)
      return ctx.response.status(200).send(files)
  }

  // // GET url/:id -- removed because nginx
  // public async show(ctx: HttpContextContract) {
  //   return ctx.response.attachment(
  //     Application.tmpPath('uploads', ctx.params.filename)
  //   )
  // }

  // POST url/
  public async store(ctx: HttpContextContract) {
    const file = ctx.request.file('file', {
      size: '128mb'
    })
    if (!file) {
      return ctx.response.status(400).send(ErrorConstructors.publicError('Запрос не содержит файл.'))
    }
    if (!file.isValid) {
      return ctx.response.status(400).send(ErrorConstructors.publicError('С файлом что-то не так. Возможно, он поврежден.'))
    }
    const user = ctx['user']
    const _tempPath = Application.tmpPath(`uploads/temp`)
    const _uploadsPath = Application.tmpPath(`uploads`)
    const extension = file.extname
    const originalFileName = file.clientName
    const newFileName = `${cuid()}.${extension}`
    const _tempPathOfFile = `${_tempPath}/${newFileName}`
    await file.move(_tempPath, {
      name: newFileName
    })
    const fileHash = await md5File(_tempPathOfFile)
      .then((hash) => {
        return hash
      })
    const foundFile = await File.findBy('hash', fileHash)
    if (foundFile) {
      // if file exists
      await fs.unlink(_tempPathOfFile, (err => { // delete file from temp folder
        if (err) {
          return ctx.response.status(500).send(ErrorConstructors.publicError('При загрузке файла произошла ошибка.'))
        }
      }))
      return ctx.response.status(200).send(foundFile)
    }
    const hash_firstTwo1 = fileHash.slice(0, 2)
    const hash_firstTwo2 = fileHash.slice(2, 4)
    const hash_other = fileHash.slice(4, fileHash.length)
    const foldersGen = `${hash_firstTwo1}/${hash_firstTwo2}/${hash_other}/`
    const _newFolder = `${_uploadsPath}/${foldersGen}`
    let isFolderExists = false
    await fs.stat(_newFolder, async (error) => {
      if (!error) {
        isFolderExists = true
      }
      if (error && error.errno === -4058) {
        isFolderExists = false
      } else {
        return ctx.response.status(500).send(ErrorConstructors.publicError('При загрузке файла произошла ошибка.'))
      }
    })
    if (!isFolderExists) {
      const result = await mkdirp(_newFolder)
        .then(() =>{
          return Promise.resolve(null)
        })
        .catch(error => {return Promise.reject(error)})
      if(result){
        return ctx.response.status(500).send(ErrorConstructors.publicError('При загрузке файла произошла ошибка.'))
      }
    }
    await fs.rename(`${_tempPath}/${newFileName}`, `${_uploadsPath}/${foldersGen}/${newFileName}`, (err) => {
      if (err) {
        throw err
      }
    })
    const newFile = new File()
    newFile.hash = fileHash
    newFile.path = `${foldersGen}/${newFileName}`
    newFile.original_name = originalFileName
    newFile.name = newFileName
    newFile.extension = extension
    newFile.size = file.size.toString()
    const savedFile = await user.related('files').save(newFile)
    console.log(savedFile)
    return ctx.response.status(200).send(savedFile)
  }


  // DELETE url/:id
  public async destroy(ctx: HttpContextContract) {

  }

}
