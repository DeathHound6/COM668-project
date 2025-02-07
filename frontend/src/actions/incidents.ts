import { APIError, type ErrorResponse, type GetManyAPIResponse, type Incident } from "../interfaces";
import { handleUnauthorized } from "./api";

export async function GetIncidents({ params }: { params: Object }): Promise<GetManyAPIResponse<Incident>> {
    const paramsString = Object.keys(params).length > 0 ?  `?${Object.entries(params).map(([key, value]) => `${key}=${value}`).join("&")}` : "";
    const response = await fetch(`/api/incidents${paramsString}`);
    if (!response.ok) {
        handleUnauthorized({ res: response });
        const data: ErrorResponse = JSON.parse(await response.text());
        throw new APIError(data.error, response.status);
    }
    return JSON.parse(await response.text());
}