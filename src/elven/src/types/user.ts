export interface IUser {
	is_admin: boolean
	username: string
	last_ip: string | null
	last_agent: string | null
}