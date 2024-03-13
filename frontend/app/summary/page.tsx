// frontend/app/summary/page.tsx
"use client";

import { useState } from "react";

interface OpenAIResponse {
	role: string;
	content: string;
}

interface Response {
	openai?: OpenAIResponse;
	// Google and Azure response types
}

const SummaryPage = () => {
	const [provider, setProvider] = useState("openai");
	const [model, setModel] = useState("gpt-3.5-turbo");
	const [content, setContent] = useState("");
	const [response, setResponse] = useState<Response | null>(null);

	const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault();

		const data = {
			data: {
				[provider]: [
					{
						model: model,
						messages: [
							{ role: "system", content: "あなたは優秀な教師です。" },
							{ role: "user", content: "下記文章を要約してください。" },
							{ role: "user", content: content },
						],
					},
					{
						model: model,
						messages: [
							{ role: "system", content: "あなたは優秀な教師です。" },
							{ role: "user", content: "下記文章を英語にしてください。" },
							{ role: "user", content: content },
						],
					},
				],
			},
		};

		try {
			const res = await fetch("/api", {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify(data),
			});

			if (res.ok) {
				const result: Response = await res.json();
				setResponse(result);
			} else {
				console.error("Request failed");
			}
		} catch (error) {
			console.error("An error occurred", error);
		}
	};

	return (
		<div className="container mx-auto px-4 py-8">
			<h1 className="text-3xl font-bold mb-4">Summary Page</h1>
			<form onSubmit={handleSubmit} className="max-w-md mx-auto">
				<div className="mb-4">
					<label className="block text-gray-700 font-bold mb-2">
						Provider:
						<select
							value={provider}
							onChange={(e) => setProvider(e.target.value)}
							className="block w-full mt-1 rounded-md border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50"
						>
							<option value="openai">OpenAI</option>
							<option value="google">Google</option>
							<option value="azure">Azure</option>
						</select>
					</label>
				</div>
				<div className="mb-4">
					<label className="block text-gray-700 font-bold mb-2">
						Model:
						<select
							value={model}
							onChange={(e) => setModel(e.target.value)}
							className="block w-full mt-1 rounded-md border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50"
						>
							<option value="gpt-3.5-turbo">gpt-3.5-turbo</option>
							<option value="gpt-4">gpt-4</option>
						</select>
					</label>
				</div>
				<div className="mb-4">
					<label className="block text-gray-700 font-bold mb-2">
						Content:
						<textarea
							value={content}
							onChange={(e) => setContent(e.target.value)}
							className="block w-full mt-1 rounded-md border-gray-300 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50"
							rows={5}
						/>
					</label>
				</div>
				<button
					type="submit"
					className="px-4 py-2 font-bold text-white bg-indigo-500 rounded-md hover:bg-indigo-700 focus:outline-none focus:shadow-outline"
				>
					Submit
				</button>
			</form>
			{response && (
				<div className="mt-8">
					<h2 className="text-xl font-bold mb-2">Response:</h2>
					{provider === "openai" && response.openai && (
						<div className="bg-gray-100 p-4 rounded-md">
							<p className="text-lg font-semibold mb-2">Role:</p>
							<p>{response.openai.role}</p>
							<p className="text-lg font-semibold mt-4 mb-2">Content:</p>
							<p>{response.openai.content}</p>
						</div>
					)}
					{/* Google and Azure response handling */}
				</div>
			)}
		</div>
	);
};

export default SummaryPage;
