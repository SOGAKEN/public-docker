import { NextRequest, NextResponse } from 'next/server';

export function middleware(request: NextRequest) {
	const token = request.cookies.get('token')?.value;
	const isLoginPage = request.nextUrl.pathname === '/login';
	const isSummaryPage = request.nextUrl.pathname === '/summary';

	if (request.nextUrl.pathname === '/') {
		if (token) {
			const summaryUrl = new URL('/summary', request.url);
			return NextResponse.redirect(summaryUrl);
		} else {
			const loginUrl = new URL('/login', request.url);
			return NextResponse.redirect(loginUrl);
		}
	}

	if (!token && !isLoginPage) {
		const loginUrl = new URL('/login', request.url);
		return NextResponse.redirect(loginUrl);
	}

	if (isLoginPage && token) {
		const summaryUrl = new URL('/summary', request.url);
		return NextResponse.redirect(summaryUrl);
	}
}

export const config = {
	matcher: ['/', '/summary', '/login'],
};
