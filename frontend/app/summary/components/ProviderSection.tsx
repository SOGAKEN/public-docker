import React from "react";
import { Provider } from "./../../../types";

type ProviderSectionProps = {
	provider: Provider;
	selectedModels: { [key: string]: string[] };
	onModelChange: (provider: string, model: string) => void;
};

const ProviderSection: React.FC<ProviderSectionProps> = ({
	provider,
	selectedModels,
	onModelChange,
}) => {
	return (
		<div className="mb-9">
			<h3 className="mb-2 text-lg font-medium">{provider.name}</h3>
			<ul className="space-y-2">
				{provider.models.map((model) => (
					<li key={model}>
						<div className="flex items-center">
							<input
								type="checkbox"
								id={`${provider.name}-${model}`}
								value={model}
								checked={selectedModels[
									provider.name
								]?.includes(model)}
								onChange={() =>
									onModelChange(provider.name, model)
								}
								className="form-checkbox h-5 w-5 text-indigo-600"
							/>
							<label
								htmlFor={`${provider.name}-${model}`}
								className="ml-2"
							>
								{model}
							</label>
						</div>
					</li>
				))}
			</ul>
		</div>
	);
};

export default ProviderSection;
