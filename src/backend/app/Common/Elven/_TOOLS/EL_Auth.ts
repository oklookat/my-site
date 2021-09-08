import Env from "@ioc:Adonis/Core/Env"
import {RequestContract} from "@ioc:Adonis/Core/Request"
import User from "App/Models/Elven/User"
import Token from "App/Models/Elven/Token"
import CryptoJS from "crypto-js"
import EL_Random from "App/Common/Elven/_TOOLS/EL_Random"

const bcrypt = require("bcrypt")

const secret = Env.get('AES_SECRET')

// * PIPE errors using only in backend and may be used for debug, don't show this errors to users

export interface IUserAndToken {
  user: User
  token: Token
}

class EL_Auth {

  public static async login(username: string, password: string, request: RequestContract): Promise<string> {
    const user = await User.findBy('username', username)
    if (!user) {
      return Promise.reject('PIPE_USER_NOT_FOUND')
    }
    const isPassword = await bcrypt.compare(password, user.password)
    if (!isPassword) {
      return Promise.reject('PIPE_WRONG_PASSWORD')
    }
    // GENERATE TOKEN START //
    const length = EL_Random.randInt(8, 24)
    const trash = EL_Random.randString(length, 'hex')
    const userID = user.id
    const userData = {
      id: userID,
      trash: trash // mix trash in token data, for unique tokens
    }
    const encryptedToken = CryptoJS.AES.encrypt(JSON.stringify(userData), secret).toString()
    const decryptedBytes = await CryptoJS.AES.decrypt(encryptedToken, secret)
    const decryptedString = decryptedBytes.toString()
    let token = new Token()
    token.token = decryptedString
    token = await EL_Auth.tokenWriteAuthAgents(request, token)
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
      const token = this.grabToken(request)
      const userAndToken = await this.getUserAndTokenByToken(token)
      await userAndToken.token.delete()
      return Promise.resolve()
    } catch (error) {
      return Promise.reject(error)
    }
  }

  public static grabToken(request: RequestContract): string {
    const cookie = request.cookie('token', null)
    if (cookie) {
      return request.cookie('token')
    }
    try {
      return this.getTokenFromHeader(request.header('Authorization'))
    } catch (error) {
      throw error
    }
  }

  private static getTokenFromHeader(authHeader: string | undefined): string {
    if (authHeader && authHeader.startsWith("Elven ")) {
      return authHeader.substring(6, authHeader.length)
    } else {
      throw 'PIPE_NO_AUTH_HEADER'
    }
  }

  public static async getUserAndTokenByRequest(request: RequestContract): Promise<IUserAndToken> {
    try {
      const token = this.grabToken(request)
      const userAndToken = await EL_Auth.getUserAndTokenByToken(token)
      return Promise.resolve({user: userAndToken.user, token: userAndToken.token})
    } catch (error) {
      return Promise.reject(error)
    }
  }

  public static async getUserAndTokenByToken(token: string): Promise<IUserAndToken> {
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
    if (!decrypted.id || !decrypted.trash) {
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

  private static tokenWriteAuthAgents(request: RequestContract, token: Token): Token {
    // this method don't saving token, you need save it manual
    const ip = request.ip()
    const agent = request.header('User-Agent')
    token.auth_ip = ip
    token.auth_agent = agent
    token.last_ip = ip
    token.last_agent = agent
    return token
  }

  public static async tokenWriteLastAgents(request: RequestContract, token: Token): Promise<Token> {
    token.last_ip = request.ip()
    token.last_agent = request.header('User-Agent')
    await token.save().catch(() => {})
    return Promise.resolve(token)
  }
}

export default EL_Auth
