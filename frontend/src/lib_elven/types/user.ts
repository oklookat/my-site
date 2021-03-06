/** information about authorized user */
export type User = {
	is_admin: boolean;
	username: string;
};

/** request body for changing username or password */
export type UserChange = {
	what: 'username' | 'password';
	password: string;
	newValue: string;
};

/** validate username */
export function usernameValidate(username: string): boolean {
	const lenValid = !(username.length < 4 || username.length > 24);
	if (!lenValid) {
		return false;
	}
	const reg1 = new RegExp('^[a-zA-Z0-9]*$');
	const isAlphanumberic = username.search(reg1);
	const isValid = isAlphanumberic !== -1;
	return isValid;
}

/** validate password */
export function passwordValidate(password: string): boolean {
	const lenValid = !(password.length < 8 || password.length > 64);
	if (!lenValid) {
		return false;
	}
	const reg1 = new RegExp(`^[a-zA-Z0-9\-+~"'\x60(){\[}|:;,.!=@#$%^&?№*\\\\/<>]*$`);
	const alphanumWithSymbols = password.search(reg1);
	const isValid = alphanumWithSymbols !== -1;
	return isValid;
}
