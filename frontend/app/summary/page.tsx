"use client";

import React, { useState } from "react";
import { saveAs } from "file-saver";
import Papa from "papaparse";
import { Provider, Response } from "./../../types";
import SideMenu from "./components/SideMenu";
import MainContent from "./components/MainContent";

const SummaryPage = () => {
	const [selectedModels, setSelectedModels] = useState<{
		[key: string]: string[];
	}>({});
	const [content, setContent] = useState("");
	const [responses, setResponses] = useState<Response[]>([]);
	const [isLoading, setIsLoading] = useState(false);

	const url = process.env.NEXT_PUBLIC_URL;
	const providers: Provider[] = JSON.parse(
		process.env.NEXT_PUBLIC_PROVIDERS || "[]"
	);

	const handleModelChange = (provider: string, model: string) => {
		if (selectedModels[provider]?.includes(model)) {
			setSelectedModels({
				...selectedModels,
				[provider]: selectedModels[provider].filter((m) => m !== model),
			});
		} else {
			setSelectedModels({
				...selectedModels,
				[provider]: [...(selectedModels[provider] || []), model],
			});
		}
	};

	const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault();

		// チェックボックスが選択されているかどうかを確認
		const isModelSelected = Object.values(selectedModels).some((models) => models.length > 0);

		if (!isModelSelected) {
			alert("左メニューよりチェックボックスにチェックを入れてください。");
			return;
		}

		setIsLoading(true);

		try {
			const promises = Object.entries(selectedModels).flatMap(
				([provider, models]) =>
					models.map(async (model) => {
						const data = {
							data: {
								[provider]: [
									{
										model: model,
										messages: [
											{
												role: "system",
												content:
													"あなたは優秀な教師です。",
											},
											{
												role: "user",
												content:
													"下記文章を要約してください。",
											},
											{ role: "user", content: content },
										],
									},
								],
							},
						};

						const res = await fetch(`${url}/api/summary`, {
							method: "POST",
							headers: {
								"Content-Type": "application/json",
							},
							body: JSON.stringify(data),
						});

						if (res.ok) {
							const result: Response = await res.json();
							return result;
						} else {
							console.error("Request failed");
							return null;
						}
					})
			);

			const results = await Promise.all(promises);
			setResponses(
				results.filter((result): result is Response => result !== null)
			);
		} catch (error) {
			console.error("An error occurred", error);
		} finally {
			setIsLoading(false);
		}
	};

	const exportToCsv = () => {
		const data = responses.map((response) => ({
			Request: content,
			Model: response.model || "",
			Response: response.openai?.content || "",
		}));

		const csv = Papa.unparse(data, { delimiter: "," });
		const blob = new Blob(["\ufeff", csv], { type: "text/csv;charset=utf-8" });

		// 現在の日時を取得
		const now = new Date();
		const timestamp = `${now.getFullYear()}${padZero(now.getMonth() + 1)}${padZero(now.getDate())}${padZero(now.getHours())}${padZero(now.getMinutes())}${padZero(now.getSeconds())}`;

		saveAs(blob, `summary_${timestamp}.csv`);
	};

	// ゼロパディング用のヘルパー関数
	const padZero = (num: number) => {
		return num.toString().padStart(2, "0");
	};

	return (
		<div className="flex h-screen">
			<SideMenu
				providers={providers}
				selectedModels={selectedModels}
				onModelChange={handleModelChange}
			/>
			<MainContent
				content={content}
				responses={responses}
				onContentChange={(e) => setContent(e.target.value)}
				onSubmit={handleSubmit}
				onExportCsv={exportToCsv}
				isLoading={isLoading}
			/>
		</div>
	);
};

export default SummaryPage;
