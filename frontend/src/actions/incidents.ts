import { APIError, type ErrorResponse, type GetManyAPIResponse, type Incident } from "../interfaces";
import { handleUnauthorized } from "./api";

export async function GetIncident({ uuid }: { uuid: string }): Promise<Incident> {
    const response = await fetch(`/api/incidents/${uuid}`);
    if (!response.ok) {
        handleUnauthorized({ res: response });
        const data: ErrorResponse = JSON.parse(await response.text());
        throw new APIError(data.error, response.status);
    }
    return JSON.parse(await response.text());
}

export async function GetIncidents({ params }: { params: object }): Promise<GetManyAPIResponse<Incident>> {
    const paramsString = Object.keys(params).length > 0 ?  `?${Object.entries(params).map(([key, value]) => `${key}=${value}`).join("&")}` : "";
    const response = await fetch(`/api/incidents${paramsString}`);
    if (!response.ok) {
        handleUnauthorized({ res: response });
        const data: ErrorResponse = JSON.parse(await response.text());
        throw new APIError(data.error, response.status);
    }
    return JSON.parse(await response.text());
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export async function UpdateIncident({ uuid, incident }: { uuid: string, incident: any }): Promise<boolean> {
    const response = await fetch(`/api/incidents/${uuid}`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(incident),
    });
    if (!response.ok) {
        handleUnauthorized({ res: response });
        const data: ErrorResponse = JSON.parse(await response.text());
        throw new APIError(data.error, response.status);
    }
    return true;
}

export async function PostComment({ uuid, comment }: { uuid: string, comment: string }): Promise<string> {
    const response = await fetch(`/api/incidents/${uuid}/comments`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ comment }),
    });
    if (!response.ok) {
        handleUnauthorized({ res: response });
        const data: ErrorResponse = JSON.parse(await response.text());
        throw new APIError(data.error, response.status);
    }
    const parts = (response.headers.get("Location") as string).split("/");
    return parts[parts.length-1];
}

export async function DeleteComment({ incidentUUID, commentUUID }: { incidentUUID: string, commentUUID: string }): Promise<boolean> {
    const response = await fetch(`/api/incidents/${incidentUUID}/comments/${commentUUID}`, {
        method: "DELETE",
    });
    if (!response.ok) {
        handleUnauthorized({ res: response });
        const data: ErrorResponse = JSON.parse(await response.text());
        throw new APIError(data.error, response.status);
    }
    return true;
}
