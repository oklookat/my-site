import Env from "@ioc:Adonis/Core/Env"
import {RequestContract} from "@ioc:Adonis/Core/Request"
import User from "App/Models/Elven/User"
import Token from "App/Models/Elven/Token"
import CryptoJS from "crypto-js"
import EL_Random from "App/Common/Elven/_TOOLS/EL_Random"

const bcrypt = require("bcrypt")

const secret = Env.get('AES_SECRET')

// * PIPE errors using only in backend and may be used for debug, don't show this errors to users

class EL_Auth {

  public static async login(username: string, password: string, isAdminLogin: boolean, request: RequestContract): Promise<string> {
    const user = await User.findBy('username', username)
    if (!user) {
      return Promise.reject('PIPE_USER_NOT_FOUND')
    }
    if (isAdminLogin && user.role !== 'admin') {
      return Promise.reject('PIPE_USER_NOT_ADMIN')
    }
    const isPassword = await bcrypt.compare(password, user.password)
    if (!isPassword) {
      return Promise.reject('PIPE_WRONG_PASSWORD')
    }
    // GENERATE TOKEN START //
    const length = EL_Random.randInt(8, 24)
    const trash = EL_Random.randString(length, 'hex')
    const userID = user.id // ObjectID
    const userData = {
      id: userID,
      trash: trash
    }
    const encryptedToken = CryptoJS.AES.encrypt(JSON.stringify(userData), secret).toString()
    const decryptedBytes = await CryptoJS.AES.decrypt(encryptedToken, secret)
    const decryptedString = decryptedBytes.toString()
    let token = new Token()
    token.token = decryptedString
    token = await EL_Auth.tokenWriteRegAgents(request, token)
    try {
      await user.related('tokens').save(token)
    } catch (error) {
      return Promise.reject('PIPE_TOKEN_SAVING_ERROR')
    }
    // GENERATE TOKEN END //
    return Promise.resolve(encryptedToken)
  }

  public static async logout(request: RequestContract) {
    try {
      const token = this.getAuthHeader(request.header('Authorization'))
      const userAndToken = await this.getUserAndTokenInstances(token)
      await userAndToken.token.delete()
      return Promise.resolve()
    } catch (error) {
      return Promise.reject(error)
    }
  }

  public static getAuthHeader(authHeader: string | undefined): string {
    if (authHeader && authHeader.startsWith("Elven ")) {
      return authHeader.substring(6, authHeader.length)
    }
    throw 'PIPE_NO_AUTH_HEADER'
  }

  public static async getUserAndTokenInstances(token: string) {
    let decryptedBytes
    let decryptedString
    let decrypted
    try {
      // decrypting token, get decrypted object and string with bytes
      decryptedBytes = await CryptoJS.AES.decrypt(token, secret)
      decryptedString = decryptedBytes.toString()
      decrypted = JSON.parse(decryptedBytes.toString(CryptoJS.enc.Utf8))
    } catch (error) {
      return Promise.reject('PIPE_TOKEN_DAMAGED')
    }
    if (!decrypted.id) {
      // if decrypted token missing id
      return Promise.reject('PIPE_TOKEN_MISSING_DATA')
    }
    // get token instance by searching his bytes in database
    const tokenInstance = await Token.findBy('token', decryptedString)
    if (!tokenInstance) {
      return Promise.reject('PIPE_TOKEN_NOT_FOUND')
    }
    const userID = decrypted.id
    const tokenUserID = tokenInstance.user_id
    if (tokenUserID !== userID) {
      // if token owner ID from DB not equal with decrypted token owner ID from the request
      return Promise.reject('PIPE_TOKEN_STRANGE_OWNER')
    }
    // find token owner
    const user = await User.find(userID)
    if (!user) {
      return Promise.reject('PIPE_TOKEN_OWNER_NOT_FOUND')
    }
    // in this moment token is valid
    return Promise.resolve({user: user, token: tokenInstance})
  }

  public static async isAdmin(request: RequestContract) {
    try {
      const token = EL_Auth.getAuthHeader(request.header('Authorization'))
      const userAndToken = await EL_Auth.getUserAndTokenInstances(token)
      return Promise.resolve({
        isAdmin: userAndToken.user.role === 'admin',
        user: userAndToken.user,
        token: userAndToken.token
      })
    } catch (error) {
      return Promise.reject(error)
    }
  }

  public static async tokenWriteLastAgents(request, token: Token) {
    token.last_ip = request.ip()
    token.last_agent = request.header('User-Agent')
    await token.save()
    return Promise.resolve(true)
  }

  private static async tokenWriteRegAgents(request, token: Token) {
    token.auth_ip = request.ip()
    token.auth_agent = request.header('User-Agent')
    token.last_ip = token.auth_ip
    token.last_agent = token.auth_agent
    return Promise.resolve(token)
  }
}

export default EL_Auth
