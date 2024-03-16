export type OpenAIResponse = {
    role: string;
    content: string;
};

export type Response = {
    openai?: OpenAIResponse;
    model?: string;
    // Google and Azure response types
};

export type Provider = {
    name: string;
    models: string[];
};
