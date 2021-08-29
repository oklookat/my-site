import fs from "fs"
import Path = require('path')
import Application from "@ioc:Adonis/Core/Application";

const crypto = require('crypto')

export default class EL_Files {

  public static _uploadsPath = Path.normalize(Application.tmpPath(`uploads`))
  public static _tempPath = Path.normalize(Application.tmpPath(`uploads/temp`))

  public static async deleteFile(path: string): Promise<null> {
    path = Path.normalize(path)
    return new Promise((resolve, reject) => {
      fs.unlink(path, (_error) => {
        if (_error) {
          reject(_error)
        }
        resolve(null)
      })
    })
  }

  public static async deleteDirIfEmpty(path: string): Promise<null> {
    path = Path.normalize(path)
    return new Promise((resolve, reject) => {
      fs.readdir(path, (_error, _files) => {
        if (_error) {
          reject(_error)
        }
        if(!_files){
          reject('PIPE_NO_FILES')
        } else if(_files.length > 0){
          reject('PIPE_DIR_NOT_EMPTY')
        }
        fs.rm(path, {recursive: true}, (_error) => {
          if (_error) {
            reject(_error)
          }
        })
        resolve(null)
      })
    })
  }

  public static deleteEmptyDirsRecursive(basePath: string, relativePathToFile: string) {
    // #1 - path like 'D:\Test\'
    // #2 - path like '123\456\789\music.flac' or '123\456\789\'
    // in summary: 'D:\Test\123\456\789\music.flac'
    // deletes the file, then goes up to the basePath, along the way deleting relativePath directories if they are empty
    relativePathToFile = relativePathToFile.replace(Path.basename(relativePathToFile), '') // remove filename from path (music.flac)
    const pathSep = Path.sep
    if(pathSep === '/'){
      // replace path slashes to POSIX
      relativePathToFile = relativePathToFile.replace(/\\/g, pathSep)
    } else {
      // to other, like win
      relativePathToFile = relativePathToFile.replace(/\//g, pathSep)
    }
    let pathArr = relativePathToFile.split(pathSep) // split path by path separator like ['123', '456', '789', ''] (last empty because maybe slash)
    let pathsForDelete: string[] = []
    for (const piece of pathArr) {
      if(piece.length < 1){
        continue
      }
      // make concat and collect paths for recursive. Example:
      // it.1 D:\Test\123\
      // it.2 D:\Test\123\456\
      // it.3 D:\Test\123\456\789\
      let _p = `${basePath}${pathSep}${piece}`
      // then add to array
      pathsForDelete.push(_p)
    }
    // reverse array, because for now we have a wrong order of paths
    pathsForDelete = pathsForDelete.reverse()
    for (const path of pathsForDelete) {
      // it.1 D:\Test\123\456\789\
      // it.2 D:\Test\123\456\
      // it.3 D:\Test\123\
      //console.log(path)
      EL_Files.deleteDirIfEmpty(path)
        .catch(() => {})
    }
    return null
  }

  public static async move(from: string, to: string): Promise<null> {
    from = Path.normalize(from)
    to = Path.normalize(to)
    return new Promise((resolve, reject) => {
      fs.rename(from, to, (_error) => {
        if (_error) {
          reject(_error)
        }
        resolve(null)
      })
    })
  }

  public static async directoryExists(path: string): Promise<boolean> {
    path = Path.normalize(path)
    return new Promise((resolve, reject) => {
      fs.stat(path, (_error) => {
        if (!_error) {
          resolve(true)
        } else if (_error && _error.errno === -4058) {
          resolve(false)
        } else {
          reject(_error)
        }
      })
    })
  }

  public static async createDirectory(path: string): Promise<null> {
    path = Path.normalize(path)
    return new Promise((resolve, reject) => {
      fs.mkdir(path, {recursive: true}, (_error) =>{
        if(_error){
          reject(_error)
        } else {
          resolve(null)
        }
      })
    })
  }

  public static async getHash(path: string): Promise<string> {
    path = Path.normalize(path)
    return new Promise((resolve, reject) => {
      // https://github.com/kodie/md5-file/blob/master/index.js
      const output = crypto.createHash('md5') // can be sha512
      const input = fs.createReadStream(path, {encoding: 'utf-8'})
      input.on('error', (err) => {
        reject(err)
      })
      output.once('readable', () => {
        resolve(output.read().toString('hex'))
      })
      input.pipe(output)
    })
  }

  public static generateDirectoriesByHash(hash: string): string {
    // MD5: hash 0 - 2 to first folder name, hash 2 - 4 to second folder name, and 4 - hash.length to third folder name
    const hash_firstTwo1 = hash.slice(0, 2)
    const hash_firstTwo2 = hash.slice(2, 4)
    const hash_other = hash.slice(4, hash.length)
    return `${hash_firstTwo1}/${hash_firstTwo2}/${hash_other}`
  }

}
