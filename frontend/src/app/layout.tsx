"use client";

import Navbar from "../components/navbar";
import { redirect, RedirectType, usePathname } from "next/navigation";
import { ToastContainer, Toast, ToastHeader, ToastBody } from "react-bootstrap";
import { useEffect, useState } from "react";
import { GetMe } from "../actions/users";
import "bootstrap/dist/css/bootstrap.min.css";
import "./globals.css";

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
    const url = usePathname();
    const [authedFor, setAuthedFor] = useState(null as string | null);

    useEffect(() => {
        async function fetchData() {
            const authCompleteFor = localStorage.getItem("auth");
            if (authCompleteFor == null)
                return;
            try {
                await GetMe();
                setAuthedFor(authCompleteFor);
                localStorage.removeItem("auth");
            } catch(e) {
                localStorage.removeItem("u");
                localStorage.removeItem("e");
                redirect("/login", RedirectType.replace);
            }
        }
        fetchData();
    }, [url]);

    return (
        <html lang="en">
            <body>
                <Navbar />
                {children}
                { url != "/login" && (
                    <ToastContainer position="bottom-end" className="p-3">
                        { authedFor && (
                            <Toast autohide delay={5000} onClose={() => setAuthedFor(null)} bg="success">
                                <ToastHeader>Authorised</ToastHeader>
                                <ToastBody>Successfully authorised for {authedFor}</ToastBody>
                            </Toast>
                        )}
                    </ToastContainer>
                )}
            </body>
        </html>
    );
}
