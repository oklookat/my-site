import { BaseCommand } from '@adonisjs/core/build/standalone'
import User from "App/Models/Elven/User"

export default class Superuser extends BaseCommand {


  public static commandName = 'elven:superuser'

  public static description = 'Создать аккаунт админа'

  public static settings = {
    loadApp: true,
    stayAlive: false,
  }

  public async run () {
    const username = await this.prompt.ask('Имя пользователя:')
    const user = await User.findBy('username', username)
    if(user){
      const isDelete = await this.prompt.confirm('Пользователь с таким именем уже существует. Удалить пользователя?')
      if(isDelete){
        await user.delete()
          .catch(error =>{
            this.logger.error(error)
            return
          })
        const isCreateNew = await this.prompt.confirm('Создать нового пользователя?')
        if(!isCreateNew){
          return
        }
      } else{
        return
      }
    }
    const password = await this.prompt.ask('Пароль:')
    await this.createUser(username, password)
      .then(() =>{
        this.logger.success('Пользователь создан!')
      })
      .catch(error =>{
        this.logger.error(error)
      })
    return
  }

  public async createUser(username: string, password: string){
    const user = {
      role: 'admin',
      username: username,
      password: password,
    }
    return await User.create(user)
      .then((document) =>{
        return Promise.resolve(document)
      })
      .catch(error =>{
        return Promise.reject(error)
      })
  }
}
