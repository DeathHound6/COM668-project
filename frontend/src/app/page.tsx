"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function Home() {
    const router = useRouter();
    useEffect(() => {
        const token = localStorage.getItem("j");
        if (token == null)
            router.push("/login");
        else
            router.push("/dashboard");
    }, []);
    return (<></>);
}
