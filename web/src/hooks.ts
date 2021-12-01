import type { MaybePromise } from '@sveltejs/kit/types/helper';
import type { ServerRequest, ServerResponse } from '@sveltejs/kit/types/hooks';

const publicPages = ['/login'];

export async function handle({
	request,
	resolve
}: {
	request: ServerRequest<Record<string, string>, Body>;
	resolve(request: ServerRequest<Record<string, string>, Body>): MaybePromise<ServerResponse>;
}): Promise<ServerResponse> {
	var backendURL = 'http://localhost:8080/api/user';

	if (import.meta.env.MODE == 'production') {
		backendURL = 'http://backend:8080/api/user';
	}

	const res = await fetch(backendURL, {
		headers: request.headers
	});

	if (res.status !== 200 && !publicPages.includes(request.path)) {
		// If you are not logged in and you are not on a public page,
		// it just redirects you to the main page, which is / in this case.
		return {
			status: 301,
			headers: {
				location: '/login'
			}
		};
	}

	if (res.status === 200) {
		const user = await res.json();
		request.locals.user = user;
		if (!publicPages.includes(request.path) && user && !user.Activated) {
			return {
				status: 301,
				headers: {
					location: '/login'
				}
			};
		}
	}

	const response = await resolve(request);

	return {
		...response
	};
}

export function getSession(request: ServerRequest<Record<string, string>, Body>): string {
	return request.locals.user;
}
