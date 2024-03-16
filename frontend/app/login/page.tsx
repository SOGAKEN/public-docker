"use client";

import { useState } from "react";
import LoginForm from "../components/LoginForm";

const LoginPage = () => {
    const url = process.env.NEXT_PUBLIC_URL;
    const [error, setError] = useState("");

    const handleLogin = async (username: string, password: string) => {
        try {
            const response = await fetch(`${url}/api/login`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    username: username,
                    password: password,
                }),
            });

            if (response.ok) {
                const data = await response.json();
                const token = data.token;
                const expiresIn = data.expiresIn;

                document.cookie = `token=${token}; max-age=${expiresIn}; path=/; samesite=lax; secure`;
                window.location.href = "/summary";
            } else {
                setError("Invalid username or password");
            }
        } catch (error) {
            setError("An error occurred. Please try again.");
        }
    };

    return (
        <div className="flex min-h-full flex-col justify-center px-6 py-12 lg:px-8">
            <div className="sm:mx-auto sm:w-full sm:max-w-sm">
                {/* img */}
                <h2 className="mt-10 text-center text-2xl font-bold leading-9 tracking-tight text-gray-900">
                    Sign in to your account
                </h2>
            </div>
            <LoginForm onLogin={handleLogin} error={error} />
        </div>
    );
};

export default LoginPage;
