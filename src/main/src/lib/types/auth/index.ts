/** auth request body */
export type body = {
    username: string
    password: string
    type: 'cookie' | 'direct'
}