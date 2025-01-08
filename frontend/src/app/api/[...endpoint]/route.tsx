import { NextRequest } from "next/server";

export async function GET(request: NextRequest, { params }: { params: { endpoint: string[] } }) {
    return await fetch(`https://com668-backend:5000/${params.endpoint.join("/")}`, {
        method: "GET",
        headers: request.headers
    });
}

export async function POST(request: NextRequest, { params }: { params: { endpoint: string[] } }) {
    return await fetch(`https://com668-backend:5000/${params.endpoint.join("/")}`, {
        method: "POST",
        body: await request.text(),
        headers: request.headers
    });
}

export async function PATCH(request: NextRequest, { params }: { params: { endpoint: string[] } }) {
    return await fetch(`https://com668-backend:5000/${params.endpoint.join("/")}`, {
        method: "PATCH",
        body: await request.text(),
        headers: request.headers
    });
}

export async function PUT(request: NextRequest, { params }: { params: { endpoint: string[] } }) {
    return await fetch(`https://com668-backend:5000/${params.endpoint.join("/")}`, {
        method: "PUT",
        body: await request.text(),
        headers: request.headers
    });
}

export async function DELETE(request: NextRequest, { params }: { params: { endpoint: string[] } }) {
    return await fetch(`https://com668-backend:5000/${params.endpoint.join("/")}`, {
        method: "DELETE",
        headers: request.headers
    });
}