"use client";

import { redirect, RedirectType } from "next/navigation";
import { useEffect, useState } from "react";

export default function Navbar() {
    const [user, setUser] = useState({} as any);

    useEffect(() => {
        // Handle the case where the JWT token is expired on page load
        const userinfo = localStorage.getItem("u");
        if (userinfo == null)
            return;
        const userjson = JSON.parse(userinfo);
        // split by "." and get the second part, which is the payload
        // parse the payload as JSON and extract the expiry timestamp field
        const expireTimestamp = JSON.parse(Buffer.from(userjson["j"].split(".")[1], "base64").toString("utf-8"))["exp"];
        if (expireTimestamp > Date.now()) {
            localStorage.removeItem("u");
            redirect("/login", RedirectType.replace);
        }

        setUser(userjson["u"]);
    }, []);

    return (<div>{user["name"]}</div>);
}