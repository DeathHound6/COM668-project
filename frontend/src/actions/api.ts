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
