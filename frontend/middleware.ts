import { NextRequest, NextResponse } from 'next/server';
import { jwtVerify } from 'jose';

// トークン検証関数
async function verifyToken(token: string | undefined, secretKey: Uint8Array): Promise<boolean> {
	if (!token) return false;
	try {
		await jwtVerify(token, secretKey);
		return true; // トークンが有効
	} catch (error) {
		return false; // トークンが無効
	}
}

export async function middleware(request: NextRequest) {
	const token = request.cookies.get('token')?.value;
	const currentPath = request.nextUrl.pathname;
	const secretKey = new TextEncoder().encode(process.env.JWT_SECRET_KEY);

	// トークンの検証結果に基づいて処理
	const isValidToken = await verifyToken(token, secretKey);

	// ログインページの処理
	if (currentPath === '/login') {
		if (isValidToken) {
			// トークンが有効ならサマリーページにリダイレクト
			return NextResponse.redirect(new URL('/summary', request.url));
		}
		// トークンが無効または存在しない場合はログインページに留まる
		return NextResponse.next();
	}

	// サマリーページまたはトップページの処理
	if (currentPath === '/summary' || currentPath === '/') {
		if (!isValidToken) {
			// トークンが無効または存在しない場合はログインページにリダイレクト
			return NextResponse.redirect(new URL('/login', request.url));
		}
		// トップページにアクセスしてトークンが有効ならサマリーページにリダイレクト
		if (currentPath === '/') {
			return NextResponse.redirect(new URL('/summary', request.url));
		}
		// サマリーページにアクセスしてトークンが有効な場合はそのまま表示
		return NextResponse.next();
	}

	// 上記以外のパスでは特にリダイレクトせず次の処理へ
	return NextResponse.next();
}

export const config = {
	matcher: ['/', '/summary', '/login'], // ミドルウェアが適用されるパス
};
