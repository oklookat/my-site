import Env from "@ioc:Adonis/Core/Env"
import {RequestContract} from "@ioc:Adonis/Core/Request"
import User from "App/Models/Elven/User"
import Token from "App/Models/Elven/Token"
import CryptoJS from "crypto-js"
import ErrorConstructors from "App/Common/Elven/_TOOLS/ErrorConstructors";
import MakeRandom from "App/Common/Elven/_TOOLS/MakeRandom";

const bcrypt = require("bcrypt")

const secret = Env.get('AES_SECRET')


class AuthMaster {

  public static async login(username: string, password: string, adminLogin: boolean, request: RequestContract): Promise<string> {
    const user = await User.findBy('username', username)
    if (!user) {
      const error = await ErrorConstructors.privateError('USER_NOT_FOUND', 'Пользователь не найден.')
      return Promise.reject(error)
    }
    if (adminLogin && user.role !== 'admin') {
      const error = await ErrorConstructors.privateError('USER_NOT_ADMIN', 'Пользователь не админ.')
      return Promise.reject(error)
    }
    const isPassword = await bcrypt.compare(password, user.password)
    if (!isPassword) {
      const error = ErrorConstructors.privateError('WRONG_PASSWORD', 'Неверный пароль.')
      return Promise.reject(error)
    }
    const length = await MakeRandom.randInt(8, 24)
    const trash = await MakeRandom.randString(length, 'hex')
    const userID = user.id // ObjectID
    const userData = {
      id: userID,
      trash: trash
    }
    const encrypted = await CryptoJS.AES.encrypt(JSON.stringify(userData), secret).toString()
    const decryptedBytes = await CryptoJS.AES.decrypt(encrypted, secret)
    const decryptedString = decryptedBytes.toString()
    let token = new Token()
    token.token = decryptedString
    token = await AuthMaster.tokenWriteRegAgents(request, token)
    try {
      await user.related('tokens').save(token)
    } catch (error) {
      const err = await ErrorConstructors.privateError('TOKEN_SAVING_ERROR', 'Ошибка при сохранении токена.')
      return Promise.reject(err)
    }
    return Promise.resolve(encrypted)
  }

  public static async logout(request: RequestContract) {
    try {
      const token = await this.getAuthHeader(request)
      const userAndToken = await this.getUserAndTokenInstances(token)
      await userAndToken.token.delete()
      return Promise.resolve(true)
    } catch (error) {
      return Promise.reject(error)
    }
  }

  public static async getAuthHeader(request: RequestContract): Promise<string> {
    let authHeader = request.header('Authorization')
    if (authHeader && authHeader.startsWith("Elven ")) {
      const token = authHeader.substring(6, authHeader.length)
      return Promise.resolve(token)
    }
    const error = await ErrorConstructors.privateError('NO_AUTH_HEADER', 'Authorization header не предоставлен.')
    return Promise.reject(error)
  }

  public static async getUserAndTokenInstances(token: string) {
    let decryptedBytes
    let decryptedString
    let decrypted
    try {
      // дешифруем токен, получаем дешифрованный объект и строку с байтами
      decryptedBytes = await CryptoJS.AES.decrypt(token, secret)
      decryptedString = decryptedBytes.toString()
      decrypted = await JSON.parse(decryptedBytes.toString(CryptoJS.enc.Utf8))
    } catch (error) {
      const err = await ErrorConstructors.privateError('TOKEN_DAMAGED', 'Токен поврежден.')
      return Promise.reject(err)
    }
    if (!decrypted.id) {
      // если в дешифрованном токене нет ID
      const err = await ErrorConstructors.privateError('TOKEN_MISSING_DATA', 'В токене нет необходимых данных.')
      return Promise.reject(err)
    }
    // получаем сущность токена по поиску его байтов в базе
    const tokenInstance = await Token.findBy('token', decryptedString)
    if (!tokenInstance) {
      const err = await ErrorConstructors.privateError('TOKEN_NOT_FOUND', 'Токен не найден.')
      return Promise.reject(err)
    }
    const userID = decrypted.id
    const tokenUserID = tokenInstance.user_id
    if (tokenUserID !== userID) {
      // если ID владельца токена из БД не совпадает с ID владельца из присланного дешифрованного токена
      const err = await ErrorConstructors.privateError('TOKEN_STRANGE_OWNER', 'Токен не принадлежит этому пользователю.')
      return Promise.reject(err)
    }
    // ищем владельца токена
    const user = await User.find(userID)
    if (!user) {
      const err = await ErrorConstructors.privateError('TOKEN_OWNER_NOT_FOUND', 'Владелец токена не найден.')
      return Promise.reject(err)
    }
    // на этом моменте можно считать токен верным, и далее делать что нужно
    return Promise.resolve({user: user, token: tokenInstance})
  }

  public static async isAdmin(request: RequestContract) {
    try {
      const token = await AuthMaster.getAuthHeader(request)
      const userAndToken = await AuthMaster.getUserAndTokenInstances(token)
      if (userAndToken.user.role === 'admin') {
        return Promise.resolve({isAdmin: true, user: userAndToken.user, token: userAndToken.token})
      }
      return Promise.resolve({isAdmin: false, user: userAndToken.user, token: userAndToken.token})
    } catch (error) {
      return Promise.reject(error)
    }
  }

  public static async tokenWriteLastAgents(request, token: Token){
    token.last_ip = request.ip()
    token.last_agent = request.header('User-Agent')
    await token.save()
    return Promise.resolve(true)
  }

  private static async tokenWriteRegAgents(request, token: Token){
    token.auth_ip = request.ip()
    token.auth_agent = request.header('User-Agent')
    token.last_ip = token.auth_ip
    token.last_agent = token.auth_agent
    return Promise.resolve(token)
  }
}

export default AuthMaster
