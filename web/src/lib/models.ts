export interface File {
	path: string;
	name: string;
	content_type: string;
}

export interface Token {
	username: string;
	permissions: string[];
	activated: boolean;
}

export class User {
	username: string;
	permissions: string[];
	activated: boolean;

	public constructor(init?: Partial<User>) {
		Object.assign(this, init);
	}

	public hasPermission(permission: string): boolean {
		if (this.permissions === undefined) {
			return false;
		}
		if (this.permissions.indexOf(permission) === -1) {
			return false;
		}
		return true;
	}
}
