export type OpenAIResponse = {
	role: string;
	content: string;
};

export type Response = {
	openai?: OpenAIResponse;
	model?: string;
	google?: OpenAIResponse;
	// Google and Azure response types
};

export type Provider = {
	name: string;
	models: string[];
};
