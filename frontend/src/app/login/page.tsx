"use client";

import { login } from "../../actions/auth";
import { useActionState } from "react";

export default function LoginPage() {
    const [state, action, pending] = useActionState(login, undefined);
    return (
        <div>
            <form action={action}>
                <input type="email" name="email" placeholder="Email" />
                <input type="password" name="password" placeholder="Password" />
                <button type="submit">Login</button>
            </form>
        </div>
    )
}