import type { GetManyAPIResponse, ErrorResponse } from "../interfaces/api";
import type { Team } from "../interfaces/user";
import { APIError } from "../interfaces/error";
import { handleUnauthorized } from "./api";

export async function GetTeams({ page, pageSize }: { page?: number, pageSize?: number }): Promise<GetManyAPIResponse<Team>> {
    const query = new URLSearchParams();
    if (page) query.set("page", page.toString());
    if (pageSize) query.set("pageSize", pageSize.toString());

    const response = await fetch(`/api/teams${query.size > 0 ? `?${query.toString()}` : ""}`);
    if (response.status != 200) {
        handleUnauthorized(response);
        const data: ErrorResponse = JSON.parse(await response.text());
        throw new APIError(data.error, response.status);
    }
    const data: GetManyAPIResponse<Team> = JSON.parse(await response.text());
    return data;
}
