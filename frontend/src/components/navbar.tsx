"use client";

import { getMe } from "../actions/auth";
import { redirect, RedirectType, usePathname } from "next/navigation";
import { useEffect, useState } from "react";
import { Navbar, NavbarText, NavItem, NavLink } from "react-bootstrap";

export default function NavbarComponent() {
    const pathname = usePathname();
    const [user, setUser] = useState({} as any);
    const [authedFor, setAuthedFor] = useState(null as any);

    useEffect(() => {
        // Handle the case where the JWT token is expired on page load
        let userinfo = localStorage.getItem("u");
        const expireTimestamp = localStorage.getItem("e");
        if (window.location.pathname.toLowerCase() != "/login" && (userinfo == null || expireTimestamp == null))
            redirect("/login", RedirectType.replace);

        if (expireTimestamp != null && parseInt(expireTimestamp) < Date.now()) {
            localStorage.removeItem("u");
            localStorage.removeItem("e");
            if (window.location.pathname.toLowerCase() != "/login")
                redirect("/login", RedirectType.replace);
        }

        // if the user is on the login page and already logged in, redirect them to the dashboard
        if (window.location.pathname.toLowerCase() == "/login" && userinfo != null)
            redirect("/dashboard", RedirectType.replace);
        setUser(userinfo != null ? JSON.parse(userinfo) : null);
        
        const authCompleteFor = localStorage.getItem("auth");
        if (authCompleteFor == null)
            return;
        getMe()
        .then(user => {
            setAuthedFor(authCompleteFor);
            setUser(user);
            localStorage.removeItem("auth");
        })
        .catch(e => {
            console.error(e);
            localStorage.removeItem("u");
            localStorage.removeItem("e");
            redirect("/login", RedirectType.replace);
        });
    }, [pathname]);

    function authSlack() {
        localStorage.setItem("auth", "Slack");
        redirect("https://localhost:5000/authorise/slack", RedirectType.replace);
    }

    return (
        <Navbar className="mt-1 mx-2 p-2 border-b">
            <NavItem className="p-2 m-1 border rounded" hidden={user == null}>
                <NavLink href="/dashboard">Dashboard</NavLink>
            </NavItem>
            <hr style={{rotate: "90"}} className="my-1" />
            <NavItem className="p-2 m-1 border rounded" hidden={user == null || user["admin"] == false}>
                <NavLink href="/settings">Settings</NavLink>
            </NavItem>
            <NavItem className="p-2 m-1 border rounded" hidden={user == null || user["admin"] == false}>
                <NavLink href="/hosts">Host Inventory</NavLink>
            </NavItem>
            <NavItem className="me-auto"></NavItem>
            { user != null
              ? (<NavItem className="p-2 m-1">
                    <NavbarText>Signed in as <span style={{color: "grey"}}>{user["name"]}</span></NavbarText>
                </NavItem>)
              : (<NavItem className="p-2 m-1">
                    <NavbarText>Not logged in</NavbarText>
                </NavItem>)
            }
            <NavItem hidden={user == null}>
                <NavLink href="#auth" style={{alignItems: "center", color: "#fff", backgroundColor: "#4A154B", border: 0, borderRadius: "48px", display: "inline-flex", fontFamily: "Lato, sans-serif", fontSize: "16px", fontWeight: 600, height: "48px", justifyContent: "center", textDecoration: "none", width: "190px"}} onClick={authSlack}>
                    <svg xmlns="http://www.w3.org/2000/svg" style={{height: "20px", width: "20px", marginRight: "12px"}} viewBox="0 0 122.8 122.8">
                        <path d="M25.8 77.6c0 7.1-5.8 12.9-12.9 12.9S0 84.7 0 77.6s5.8-12.9 12.9-12.9h12.9v12.9zm6.5 0c0-7.1 5.8-12.9 12.9-12.9s12.9 5.8 12.9 12.9v32.3c0 7.1-5.8 12.9-12.9 12.9s-12.9-5.8-12.9-12.9V77.6z" fill="#e01e5a" />
                        <path d="M45.2 25.8c-7.1 0-12.9-5.8-12.9-12.9S38.1 0 45.2 0s12.9 5.8 12.9 12.9v12.9H45.2zm0 6.5c7.1 0 12.9 5.8 12.9 12.9s-5.8 12.9-12.9 12.9H12.9C5.8 58.1 0 52.3 0 45.2s5.8-12.9 12.9-12.9h32.3z" fill="#36c5f0" />
                        <path d="M97 45.2c0-7.1 5.8-12.9 12.9-12.9s12.9 5.8 12.9 12.9-5.8 12.9-12.9 12.9H97V45.2zm-6.5 0c0 7.1-5.8 12.9-12.9 12.9s-12.9-5.8-12.9-12.9V12.9C64.7 5.8 70.5 0 77.6 0s12.9 5.8 12.9 12.9v32.3z" fill="#2eb67d" />
                        <path d="M77.6 97c7.1 0 12.9 5.8 12.9 12.9s-5.8 12.9-12.9 12.9-12.9-5.8-12.9-12.9V97h12.9zm0-6.5c-7.1 0-12.9-5.8-12.9-12.9s5.8-12.9 12.9-12.9h32.3c7.1 0 12.9 5.8 12.9 12.9s-5.8 12.9-12.9 12.9H77.6z" fill="#ecb22e" />
                    </svg>
                    Authorise Slack
                </NavLink>
            </NavItem>
        </Navbar>
    );
}