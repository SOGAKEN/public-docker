import React from "react";
import { Response } from "./../../../types";

type ResponseSectionProps = {
	responses: Response[];
};

const ResponseSection: React.FC<ResponseSectionProps> = ({ responses }) => {
	return (
		<div className="mt-8">
			{responses.map((response, index) => (
				<div key={index} className="mb-4">
					{/* <h2 className="text-xl font-bold mb-2">Response:</h2> */}
					{response.openai && (
						<div className="bg-gray-100 p-4 rounded-md">
							<p className="text-lg font-semibold mb-2">Model:</p>
							<p>{response.model}</p>
							<p className="text-lg font-semibold mt-4 mb-2">
								Content:
							</p>
							<p>{response.openai.content}</p>
						</div>
					)}
					{response.google && (
						<div className="bg-gray-100 p-4 rounded-md">
							<p className="text-lg font-semibold mb-2">Model:</p>
							<p>{response.model}</p>
							<p className="text-lg font-semibold mt-4 mb-2">
								Content:
							</p>
							<p>{response.google.content}</p>
						</div>
					)}
					{response.claude && (
						<div className="bg-gray-100 p-4 rounded-md">
							<p className="text-lg font-semibold mb-2">Model:</p>
							<p>{response.model}</p>
							<p className="text-lg font-semibold mt-4 mb-2">
								Content:
							</p>
							<p>{response.claude.content}</p>
						</div>
					)}
					{/* Google and Azure response handling */}
				</div>
			))}
		</div>
	);
};

export default ResponseSection;
