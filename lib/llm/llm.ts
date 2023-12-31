export interface Message {
	role: 'system' | 'user' | 'tool';
	content: string | null;
}

export interface ChatCompletion {
	choices: Array<{
		finish_reason: 'stop' | 'tool_calls' | 'length' | 'content_filter';
		message: {
			role: 'system' | 'user' | 'tool';
			content: string | null;
			tool_calls: Array<{
				type: 'function';
				function: {
					name: string;
					arguments: string;
				};
			}>;
		};
	}>;
}

export interface LLMConstructorArgs {
	messages?: Message[];
}

export abstract class LLM {
	messages: Message[] = [];

	constructor({ messages }: LLMConstructorArgs = {}) {
		this.messages = messages ?? [];
	}

	abstract createChatCompletion: () => Promise<ChatCompletion>;
}
