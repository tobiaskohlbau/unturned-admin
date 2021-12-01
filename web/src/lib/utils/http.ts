import { Notify } from '$lib/store';

export class HttpResponse<T> extends Response {
	parsedBody?: T | string;
}

export async function http<T>(request: RequestInfo, parse = 'json'): Promise<HttpResponse<T>> {
	const response: HttpResponse<T> = await fetch(request);

	if (parse) {
		try {
			switch (parse) {
				case 'json':
					response.parsedBody = await response.json();
					break;
				case 'text':
					response.parsedBody = await response.text();
					break;
			}
		} catch (ex) {
			Notify({
				type: 'error',
				message: ex
			});
		}
	}

	if (!response.ok) {
		throw new Error(response.statusText);
	}
	return response;
}

export async function httpPost<Req extends BodyInit, Res>(
	path: string,
	body: Req,
	parse = 'json'
): Promise<HttpResponse<Res>> {
	return await http<Res>(
		new Request(path, {
			method: 'POST',
			body: body
		}),
		parse
	);
}

export async function httpPut<Req extends BodyInit, Res>(
	path: string,
	body: Req,
	parse = 'json'
): Promise<HttpResponse<Res>> {
	return await http<Res>(
		new Request(path, {
			method: 'PUT',
			body: body
		}),
		parse
	);
}

export async function httpDelete<T>(path: string, parse = 'json'): Promise<HttpResponse<T>> {
	return await http<T>(
		new Request(path, {
			method: 'DELETE'
		}),
		parse
	);
}
