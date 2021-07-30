import fs from "fs"
import mkdirp from 'mkdirp'
import Path = require('path')

const crypto = require('crypto')

export default class EL_Files {

  public static async deleteFile(path: string) {
    let error
    await fs.unlink(path, (err) => {
      if (err) {
        error = err
      }
    })
    if (error) {
      return Promise.reject(error)
    }
    return Promise.resolve()
  }

  public static async deleteDirIfEmpty(path: string) {
    let error
    await fs.rm(path, {recursive: true}, (err) => {
      if (err) {
        error = err
      }
    })
    if (!error) {
      return Promise.resolve(null)
    } else {
      return Promise.reject(error)
    }
  }

  public static async deleteEmptyDirsRecursive(basePath: string, relativePathToFile: string) {
    /**
     * @param string path like 'D:\Test\'
     * @param string path like '123\456\789\music.flac' or '123\456\789\'
     * in summary: 'D:\Test\123\456\789\music.flac'
     * @return Promise
     * deletes the file, then goes up to the basePath, along the way deleting relativePath directories if they are empty
     */
    relativePathToFile = relativePathToFile.replace(Path.basename(relativePathToFile), '') // remove filename from path (music.flac)
    let pathArr = relativePathToFile.split('/') // split path by slash like ['123', '456', '789', ''] (last empty because slash)
    let goDeep = '' // need for recursive concat
    let pathsForDelete: string[] = []
    for (const piece of pathArr) {
      if (piece.length > 0) { // check is array element not empty
        // make concat and collect paths for recursive. Example:
        // it.1 D:\Test\123\
        // it.2 D:\Test\123\456\
        // it.3 D:\Test\123\456\789\
        goDeep = goDeep + '/' + piece
        // then add to array, and then reverse this array, because for now we have a wrong order of paths
        let _p = `${basePath}/${goDeep}`
        _p = Path.normalize(_p)
        pathsForDelete.push(_p)
      }
    }
    pathsForDelete = pathsForDelete.reverse()
    for (const path of pathsForDelete) {
      // it.1 D:\Test\123\456\789\
      // it.2 D:\Test\123\456\
      // it.3 D:\Test\123\
      await EL_Files.deleteDirIfEmpty(path)
    }
  }

  public static async move(from: string, to: string) {
    let error
    await fs.rename(from, to, (err) => {
      if (err) {
        error = err
      }
    })
    if (error) {
      return Promise.reject(error)
    } else {
      return Promise.resolve(null)
    }
  }

  public static async directoryExists(path: string) {
    let error
    let exists = false
    await fs.stat(path, async (err) => {
      if (!err) {
        exists = true
      } else if (err && err.errno === -4058) {
        exists = false
      } else {
        error = err
      }
    })
    if (!error && exists) {
      return Promise.resolve(true)
    }
    if (!error && !exists) {
      return Promise.resolve(false)
    }
    if (error) {
      return Promise.reject(error)
    }
  }

  public static async createDirectory(path: string) {
    return await mkdirp(path)
      .then(() => {
        return Promise.resolve(null)
      })
      .catch(error => {
        return Promise.reject(error)
      })
  }

  public static async getHash(path: string): Promise<string> {
    // SHA-512
    return new Promise((resolve, reject) => {
      // https://github.com/kodie/md5-file/blob/master/index.js but sha512
      const output = crypto.createHash('sha512')
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
    // hash 0 - 16 to first folder name
    // hash 16 - 32 to second folder name
    // and 32 - hash.length to third folder name
    const hash_firstTwo1 = hash.slice(0, 16)
    const hash_firstTwo2 = hash.slice(16, 32)
    const hash_other = hash.slice(32, hash.length)
    return `${hash_firstTwo1}/${hash_firstTwo2}/${hash_other}`
  }

}
