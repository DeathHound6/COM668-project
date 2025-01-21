import type { Team } from "./user";

export interface HostMachine {
    uuid: string;
    os: string;
    hostname: string;
    ip4: string;
    ip6: string;
    team: Team;
}