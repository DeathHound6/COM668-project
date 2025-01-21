import type { User } from "./user";
import type { HostMachine } from "./hosts";

export interface Incident {
    uuid: string;
    hostsAffected: HostMachine[];
    summary: string;
    createdAt: string;
    resolvedAt: string | undefined;
    resolvedBy: User | undefined;
}