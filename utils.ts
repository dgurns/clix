export function systemMessage(message: string, ...args: string[]) {
	if (args) {
		return console.log(`🤖 ${message}\n`, ...args);
	}
	return console.log(`🤖 ${message}\n`);
}

export function exit() {
	console.log('');
	systemMessage('Kbye');
	Deno.exit(0);
}

export async function userPrompt(prefix?: string) {
	let res = await prompt(prefix ? `${prefix} > ` : '> ');

	if (res === null) {
		exit();
	}

	if (res === '') {
		const resV2 = await prompt('Write something, or press Enter to exit > ');
		if (!resV2) {
			exit();
		}
		res = resV2;
	}

	console.log('');
	return res;
}
