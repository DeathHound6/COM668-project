import { APIError, type ErrorResponse, type GetManyAPIResponse, type Settings } from "../interfaces";
import { handleUnauthorized } from "./api";

export async function GetSetting({ uuid }: { uuid: string }): Promise<Settings> {
    const response = await fetch(`/api/providers/${uuid}`);
    if (!response.ok) {
        handleUnauthorized({ res: response });
        const data: ErrorResponse = JSON.parse(await response.text());
        throw new APIError(data.error, response.status);
    }
    return JSON.parse(await response.text());
}

export async function GetSettings({ providerType, page, pageSize }: { providerType: "alert"|"log", page?: number, pageSize?: number}): Promise<GetManyAPIResponse<Settings>> {
    const query = new URLSearchParams();
    if (page) query.set("page", page.toString());
    if (pageSize) query.set("pageSize", pageSize.toString());
    if (providerType) query.set("provider_type", providerType);

    const response = await fetch(`/api/providers${query.size > 0 ? `?${query.toString()}` : ""}`);
    if (!response.ok) {
        handleUnauthorized({ res: response });
        const data: ErrorResponse = JSON.parse(await response.text());
        throw new APIError(data.error, response.status);
    }
    return JSON.parse(await response.text());
}

export async function CreateSetting({ name, providerType }: { name: string, providerType: "alert"|"log" }): Promise<string> {
    const response = await fetch(`/api/providers?provider_type=${providerType}`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ name })
    });
    if (!response.ok) {
        handleUnauthorized({ res: response });
        const data: ErrorResponse = JSON.parse(await response.text());
        throw new APIError(data.error, response.status);
    }
    const parts = (response.headers.get("Location") as string).split("/");
    return parts[parts.length - 1];
}

export async function UpdateSetting(setting: Settings): Promise<boolean> {
    const response = await fetch(`/api/providers/${setting.uuid}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(setting),
    });
    if (!response.ok) {
        handleUnauthorized({ res: response });
        const data: ErrorResponse = JSON.parse(await response.text());
        throw new APIError(data.error, response.status);
    }
    return true;
}

export async function DeleteSetting({ uuid }: { uuid: string }): Promise<boolean> {
    const response = await fetch(`/api/providers/${uuid}`, {
        method: "DELETE",
    });
    if (!response.ok) {
        handleUnauthorized({ res: response });
        const data: ErrorResponse = JSON.parse(await response.text());
        throw new APIError(data.error, response.status);
    }
    return true;
}
