import { NextRequest } from "next/server";

const host = "com668-backend:5000";

export async function GET(request: NextRequest, { params }: { params: { endpoint: string[] } }) {
    const { endpoint } = await params;
    const queryStrings = Array.from(request.nextUrl.searchParams.entries()).map(([key, value]) => `${key}=${value}`);
    const query = queryStrings.length > 0 ? `?${queryStrings.join("&")}` : "";
    return await fetch(`https://${host}/${endpoint.join("/")}${query}`, {
        method: "GET",
        headers: request.headers
    });
}

export async function POST(request: NextRequest, { params }: { params: { endpoint: string[] } }) {
    const { endpoint } = await params;
    const queryStrings = Array.from(request.nextUrl.searchParams.entries()).map(([key, value]) => `${key}=${value}`);
    const query = queryStrings.length > 0 ? `?${queryStrings.join("&")}` : "";
    return await fetch(`https://${host}/${endpoint.join("/")}${query}`, {
        method: "POST",
        body: await request.text(),
        headers: request.headers
    });
}

export async function PATCH(request: NextRequest, { params }: { params: { endpoint: string[] } }) {
    const { endpoint } = await params;
    const queryStrings = Array.from(request.nextUrl.searchParams.entries()).map(([key, value]) => `${key}=${value}`);
    const query = queryStrings.length > 0 ? `?${queryStrings.join("&")}` : "";
    return await fetch(`https://${host}/${endpoint.join("/")}${query}`, {
        method: "PATCH",
        body: await request.text(),
        headers: request.headers
    });
}

export async function PUT(request: NextRequest, { params }: { params: { endpoint: string[] } }) {
    const { endpoint } = await params;
    const queryStrings = Array.from(request.nextUrl.searchParams.entries()).map(([key, value]) => `${key}=${value}`);
    const query = queryStrings.length > 0 ? `?${queryStrings.join("&")}` : "";
    return await fetch(`https://${host}/${endpoint.join("/")}${query}`, {
        method: "PUT",
        body: await request.text(),
        headers: request.headers
    });
}

export async function DELETE(request: NextRequest, { params }: { params: { endpoint: string[] } }) {
    const { endpoint } = await params;
    const queryStrings = Array.from(request.nextUrl.searchParams.entries()).map(([key, value]) => `${key}=${value}`);
    const query = queryStrings.length > 0 ? `?${queryStrings.join("&")}` : "";
    return await fetch(`https://${host}/${endpoint.join("/")}${query}`, {
        method: "DELETE",
        headers: request.headers
    });
}