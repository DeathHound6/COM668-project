"use client";

import { redirect, RedirectType, usePathname } from "next/navigation";
import { useEffect, useState } from "react";
import {
    Col,
    Dropdown,
    DropdownDivider,
    DropdownItem,
    DropdownItemText,
    DropdownMenu,
    DropdownToggle,
    Nav,
    Navbar,
    NavItem,
    NavLink,
    Row
} from "react-bootstrap";
import { Slack } from "react-bootstrap-icons";

export default function NavbarComponent() {
    const pathname = usePathname();
    const [user, setUser] = useState({} as any);

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
    }, [pathname]);

    function authSlack() {
        localStorage.setItem("auth", "Slack");
        redirect("https://localhost:5000/authorise/slack", RedirectType.replace);
    }

    return (
        <Navbar className="mt-1 mx-2 p-2 border-b">
            <Nav activeKey={pathname}>
                <NavItem className="p-2 m-1 border rounded" hidden={user == null}>
                    <NavLink href="/dashboard" eventKey="/dashboard">Dashboard</NavLink>
                </NavItem>
                <NavItem className="p-2 m-1 border rounded" hidden={user == null || user["admin"] == false}>
                    <NavLink href="/settings" eventKey="/settings">Settings</NavLink>
                </NavItem>
                <NavItem className="p-2 m-1 border rounded" hidden={user == null || user["admin"] == false}>
                    <NavLink href="/hosts" eventKey="/hosts">Host Inventory</NavLink>
                </NavItem>
            </Nav>
            <NavItem className="me-auto"></NavItem>
            <NavItem hidden={user == null}>
                <Dropdown>
                    <DropdownToggle>User</DropdownToggle>
                    <DropdownMenu align="end" style={{minWidth: "13rem"}}>
                        { user == null
                            ? (<DropdownItemText>Not logged in</DropdownItemText>)
                            : (<DropdownItemText>Signed in as <span style={{color: "grey"}}>{user["name"]}</span></DropdownItemText>)
                        }
                        { user != null && (
                            <>
                                <DropdownDivider />
                                <DropdownItem onClick={authSlack}>
                                    <Row xs={4}>
                                        <Col xs={1}>
                                            <Slack />
                                        </Col>
                                        <Col xs={3}>
                                            Connect with Slack
                                        </Col>
                                    </Row>
                                </DropdownItem>
                            </>
                        )}
                    </DropdownMenu>
                </Dropdown>
            </NavItem>
        </Navbar>
    );
}