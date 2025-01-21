export interface User {
    uuid: string;
    name: string;
    email: string;
    teams: Team[];
    slackID: string;
    admin: boolean;
}

export interface Team {
    uuid: string;
    name: string;
}