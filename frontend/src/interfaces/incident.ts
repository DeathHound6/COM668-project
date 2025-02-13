import type { User, Team } from "./user";
import type { HostMachine } from "./hosts";

export interface Incident {
    uuid: string;
    description: string;
    comments: IncidentComment[];
    hostsAffected: HostMachine[];
    summary: string;
    createdAt: string;
    resolvedAt: string | undefined;
    resolvedBy: User | undefined;
    resolutionTeams: Team[];
}

export interface IncidentComment {
    uuid: string;
    comment: string;
    commentedAt: string;
    commentedBy: User;
}
