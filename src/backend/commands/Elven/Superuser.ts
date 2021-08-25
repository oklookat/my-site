import {BaseCommand} from '@adonisjs/core/build/standalone'
import User from "App/Models/Elven/User"

export default class Superuser extends BaseCommand {

  public static commandName = 'elven:superuser'

  public static description = 'create superuser account'

  public static settings = {
    loadApp: true,
    stayAlive: false,
  }

  public async run() {
    const username = await this.prompt.ask('Username:')
    const user = await User.findBy('username', username)
    if (user) {
      const isDelete = await this.prompt.confirm('Username exists. Delete?')
      if (!isDelete) {
        return
      } else {
        try {
          await user.delete()
        } catch (error) {
          this.logger.error(error)
          return
        }
        const isCreateNew = await this.prompt.confirm('Create new user?')
        if (!isCreateNew) {
          return
        }
      }
    }
    const password = await this.prompt.secure('Password:')
    try {
      await Superuser.createUser(username, password)
      this.logger.success('User created!')
    } catch (error) {
      await this.errorRender(error)
    }
    return
  }

  private static async createUser(username: string, password: string) {
    const user = {
      role: 'admin',
      username: username,
      password: password,
    }
    try {
      const userInstance = await User.create(user)
      return Promise.resolve(userInstance)
    } catch (error) {
      return Promise.reject(error)
    }
  }

  private async errorRender(error) {
    if (error["errors"]) {
      const errors = error["errors"]
      let i = 1
      errors.forEach(err => {
        this.logger.error(`---- ERROR #${i} ----`)
        const issuers: Array<string> = err.issuers
        issuers.forEach(issuer => {
          this.logger.error(`issuer: ${issuer}`)
        })
        const errorCode = err.errorCode
        this.logger.error(`errorCode: ${errorCode}`)
        switch (errorCode) {
          case "E_VALIDATION_MINMAX":
            this.logger.error(`min: ${err.min} / max: ${err.max}`)
            break
          case "E_VALIDATION_ALLOWED":
            err.allowed.forEach(allowed => {
              this.logger.error(`allowed: ${allowed}`)
            })
            break
          default:
            if (err.message) {
              this.logger.error(err.message)
            } else {
              this.logger.error("error message not provided")
            }
        }
        this.logger.error(`---- ERROR #${i} ----\n`)
        i++
      })
    } else {
      if (error) {
        this.logger.error(error)
        return
      }
      this.logger.error(`unknown error`)
      return
    }
  }
}
