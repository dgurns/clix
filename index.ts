import { systemMessage, userPrompt } from './utils.ts';
import { OpenAI } from './lib/llm/openai.ts';

console.log('');

systemMessage(`Welcome to Clix! I can help you run commands on your computer. What would you like to do?
For example, "Reorganize my desktop" or "Initialize a new git repository"`);

const initialRequest = await userPrompt();

if (!initialRequest) {
	systemMessage('Kbye');
	Deno.exit(0);
}

const llm = new OpenAI({
	messages: [{ role: 'user', content: initialRequest }],
});

const llmRes = await llm.createChatCompletion();

const toolCall = llmRes.choices[0].message.tool_calls[0];

if (toolCall) {
	const args = JSON.parse(toolCall.function.arguments);

	systemMessage(
		`${args.rationale}
Let's run: %c${args.command}

%cWould you like to run this command? (y)es / (n)o`,
		'color: yellow',
		'color: red'
	);

	const shouldRun = await userPrompt();
	if (shouldRun === 'y') {
		const command = new Deno.Command(args.command);
		const { stdout, stderr } = await command.output();
		const out = new TextDecoder().decode(stdout);
		const err = new TextDecoder().decode(stderr);
		if (out) {
			systemMessage(`Running %c${args.command}`, 'color: yellow');
			console.log(out);
		} else if (err) {
			systemMessage(`Error: ${err}`);
		}
	} else {
		systemMessage('Ok, what would you like to do instead?');
		const newRequest = await userPrompt();
	}
}
