"use client";

import { login } from "../../actions/auth";
import { useActionState } from "react";
import { FloatingLabel, FormControl, Form, Button } from "react-bootstrap";

export default function LoginPage() {
    const [state, action, pending] = useActionState(login, undefined);
    return (
        <div style={{width: "33%", alignItems: "center", textAlign: "center"}} className="mx-auto mt-40">
            <h1 className="mb-4" style={{fontSize: 24}}>Login</h1>
            <Form action={action}>
                <FloatingLabel controlId="email" label="Email Address" className="mb-3">
                    <FormControl type="email" name="email" />
                </FloatingLabel>
                <FloatingLabel controlId="password" label="Password" className="mb-3">
                    <FormControl type="password" name="password" />
                </FloatingLabel>
                <Button type="submit" variant="outline-primary">Login</Button>
            </Form>
        </div>
    )
}