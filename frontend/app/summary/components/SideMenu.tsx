import React from "react";
import { Provider } from "./../../../types";
import ProviderSection from "./ProviderSection";
import Image from "next/image";

type SideMenuProps = {
    providers: Provider[];
    selectedModels: { [key: string]: string[] };
    onModelChange: (provider: string, model: string) => void;
};

const SideMenu: React.FC<SideMenuProps> = ({
    providers,
    selectedModels,
    onModelChange,
}) => {
    return (
        <div className="w-80 bg-indigo-600 text-white p-8 overflow-y-auto">
            <h1 className="text-2xl font-bold mb-16 flex items-center">
                <Image
                    src="./../../log.svg"
                    width="30"
                    height="30"
                    alt="main log"
                />
                <span className="ml-1">Summary AI Tool</span>
            </h1>
            {providers.map((provider) => (
                <ProviderSection
                    key={provider.name}
                    provider={provider}
                    selectedModels={selectedModels}
                    onModelChange={onModelChange}
                />
            ))}
        </div>
    );
};

export default SideMenu;
