export type AIResponse = {
	role: string;
	content: string;
};

export type Response = {
	openai?: AIResponse;
	google?: AIResponse;
	model?: string;
	// Google and Azure response types
};

export type Provider = {
	name: string;
	models: string[];
};
