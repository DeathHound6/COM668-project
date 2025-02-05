import type { User } from "../interfaces/user";
import type { ErrorResponse } from "../interfaces/api";
import { APIError } from "../interfaces/error";
import { handleUnauthorized } from "./api";

export async function GetMe(): Promise<User> {
    const response = await fetch("/api/me");
    if (!response.ok) {
        handleUnauthorized({ res: response });
        const data: ErrorResponse = JSON.parse(await response.text());
        throw new APIError(data.error, response.status);
    }
    const userinfo: User = JSON.parse(await response.text());
    localStorage.setItem("u", JSON.stringify(userinfo));
    return userinfo;
}

