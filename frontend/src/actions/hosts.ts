import type { GetManyAPIResponse, ErrorResponse, HostMachine } from "../interfaces";
import { APIError } from "../interfaces/error";
import { handleUnauthorized } from "./api";

export async function GetHost({ uuid }: { uuid: string }): Promise<HostMachine> {
    const response = await fetch(`/api/hosts/${uuid}`);
    if (!response.ok) {
        handleUnauthorized({ res: response });
        const data: ErrorResponse = JSON.parse(await response.text());
        throw new APIError(data.error, response.status);
    }
    return JSON.parse(await response.text());
}

export async function GetHosts({ page, pageSize }: { page?: number, pageSize?: number }): Promise<GetManyAPIResponse<HostMachine>> {
    const query = new URLSearchParams();
    if (page) query.set("page", page.toString());
    if (pageSize) query.set("pageSize", pageSize.toString());

    const response = await fetch(`/api/hosts${query.size > 0 ? `?${query.toString()}` : ""}`);
    if (!response.ok) {
        handleUnauthorized({ res: response });
        const data: ErrorResponse = JSON.parse(await response.text());
        throw new APIError(data.error, response.status);
    }
    return JSON.parse(await response.text());
}

export async function CreateHost({ hostname, os, ip4, ip6, teamID }: { hostname: string, os: string, ip4: string|null, ip6: string|null, teamID: string }): Promise<string> {
    const response = await fetch(`/api/hosts`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ hostname, os, ip4, ip6, teamID })
    });
    if (!response.ok) {
        handleUnauthorized({ res: response });
        const data: ErrorResponse = JSON.parse(await response.text());
        throw new APIError(data.error, response.status);
    }
    const parts = (response.headers.get("Location") as string).split("/");
    return parts[parts.length-1];
}

export async function UpdateHost({ uuid, body }: {uuid: string, body: { hostname: string, os: string, ip4?: string, ip6?: string, teamID: string }}): Promise<undefined> {
    const response = await fetch(`/api/hosts/${uuid}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(body)
    });
    if (!response.ok) {
        handleUnauthorized({ res: response });
        const data: ErrorResponse = JSON.parse(await response.text());
        throw new APIError(data.error, response.status);
    }
    return undefined;
}

export async function DeleteHost(uuid: string): Promise<undefined> {
    const response = await fetch(`/api/hosts/${uuid}`, {
        method: "DELETE"
    });
    if (!response.ok) {
        handleUnauthorized({ res: response });
        const data: ErrorResponse = JSON.parse(await response.text());
        throw new APIError(data.error, response.status);
    }
    return undefined;
}
