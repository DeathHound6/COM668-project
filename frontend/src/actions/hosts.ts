import type { GetManyAPIResponse, ErrorResponse } from "../interfaces/api";
import type { HostMachine } from "../interfaces/hosts";
import { APIError } from "../interfaces/error";
import { handleUnauthorized } from "./api";

export async function GetHost({ uuid }: { uuid: string }): Promise<HostMachine> {
    const response = await fetch(`/api/hosts/${uuid}`);
    if (response.status != 200) {
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
    if (response.status != 200) {
        handleUnauthorized({ res: response });
        const data: ErrorResponse = JSON.parse(await response.text());
        throw new APIError(data.error, response.status);
    }
    return JSON.parse(await response.text());
}

export async function CreateHost({ hostname, os, ip4, ip6, teamID }: { hostname: string, os: string, ip4: string, ip6: string, teamID: string }): Promise<string> {
    const response = await fetch(`/api/hosts`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ hostname, os, ip4, ip6, teamID })
    });
    if (response.status != 201) {
        handleUnauthorized({ res: response });
        const data: ErrorResponse = JSON.parse(await response.text());
        throw new APIError(data.error, response.status);
    }
    const parts = (response.headers.get("Location") as string).split("/");
    return parts[parts.length-1];
}

export async function UpdateHost({ uuid, body }: {uuid: string, body: { hostname: string, os: string, ip4?: string, ip6?: string, teamID: string }}): Promise<void> {
    const response = await fetch(`/api/hosts/${uuid}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(body)
    });
    if (response.status != 200) {
        handleUnauthorized({ res: response });
        const data: ErrorResponse = JSON.parse(await response.text());
        throw new APIError(data.error, response.status);
    }
}

export async function DeleteHost(uuid: string): Promise<void> {
    const response = await fetch(`/api/hosts/${uuid}`, {
        method: "DELETE"
    });
    if (response.status != 204) {
        handleUnauthorized({ res: response });
        const data: ErrorResponse = JSON.parse(await response.text());
        throw new APIError(data.error, response.status);
    }
}
