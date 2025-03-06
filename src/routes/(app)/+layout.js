export const prerender = true;
export const trailingSlash = "always";
import { getUser } from "$lib/api";
export async function load() {
	try {
		let user = await getUser();
		return {
			user: user
		};
	} catch (error) {
		return {
			user: null
		};
	}
}
