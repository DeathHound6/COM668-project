import type { HostMachine } from "./hosts";
import type { Incident } from "./incident";
import type { Settings } from "./settings";
import type { Team } from "./user";

export interface ErrorResponse {
    error: string;
}

export interface GetManyAPIResponse<T=HostMachine|Incident|Team|Settings> {
    data: T[];
    meta: {
        total: number;
        pages: number;
        page: number;
        pageSize: number;
    };
}