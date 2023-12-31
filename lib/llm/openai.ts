import { LLM, Message, ChatCompletion, LLMConstructorArgs } from './llm.ts';

export class OpenAI implements LLM {
	messages: Message[] = [];

	constructor({ messages }: LLMConstructorArgs = {}) {
		this.messages = messages ?? [];
	}

	async createChatCompletion(): Promise<ChatCompletion> {
		return await new Promise((res) =>
			res({
				choices: [
					{
						finish_reason: 'stop',
						message: {
							role: 'tool',
							content: null,
							tool_calls: [
								{
									type: 'function',
									function: {
										name: 'run_command',
										arguments: JSON.stringify({
											command: 'ls',
											rationale: 'We need to see what this folder contains',
										}),
									},
								},
							],
						},
					},
				],
			})
		);
	}
}
