import { APIError } from "../interfaces/error";
import { redirect, RedirectType } from "next/navigation";

export function handleUnauthorized({ res, err }: { res?: Response | undefined, err?: APIError | undefined }) {
    if ((res != undefined && res.status == 401) || (err != undefined && err.status == 401))
    {
        localStorage.removeItem("u");
        localStorage.removeItem("e");
        redirect("/login", RedirectType.replace);
    }
}

export function formatDate(date: Date): string {
    const months = ["Jan", "Feb", "Mar", "Apr","May", "Jun", "Jul", "Aug","Sep", "Oct", "Nov", "Dec"];
    // prepend 0 to single digit hours and minutes
    return `${months[date.getMonth()]} ${date.getDate()} ${date.getFullYear()}, ${padStringWith0(date.getHours().toString(), 2, false)}:${padStringWith0(date.getMinutes().toString(), 2, false)}`;
}

function padStringWith0(str: string, size: number, padRight: boolean = true): string {
    if (str.length >= size)
        return str;
    while (str.length < size) {
        if (padRight)
            str += "0";
        else
            str = "0" + str;
    }
    return str;
}
