//femabras/frontend/src/middleware.ts
import { NextRequest, NextResponse } from "next/server";

const locales = ["en", "pt", "fr"];
const defaultLocale = "en";

function getPreferredLocale(request: NextRequest): string {
  const acceptLanguage = request.headers.get("accept-language") || "";
  if (acceptLanguage.includes("pt")) return "pt";
  if (acceptLanguage.includes("fr")) return "fr";
  return defaultLocale;
}

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;

  if (
    pathname.startsWith("/_next") ||
    pathname.startsWith("/api") ||
    pathname.includes(".") ||
    pathname === "/favicon.ico"
  ) {
    return;
  }

  const pathnameHasLocale = locales.some(
    (locale) => pathname.startsWith(`/${locale}/`) || pathname === `/${locale}`,
  );

  if (pathnameHasLocale) return NextResponse.next();

  const locale = getPreferredLocale(request);
  request.nextUrl.pathname = `/${locale}${pathname}`;
  return NextResponse.redirect(request.nextUrl);
}

export const config = {
  matcher: [
    "/((?!api|_next/static|_next/image|favicon.ico|sitemap.xml|robots.txt|fonts).*)",
  ],
};
