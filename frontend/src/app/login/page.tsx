"use client";

import { GetMe } from "../../actions/users";
import { useState } from "react";
import {
    FloatingLabel,
    FormControl,
    Form,
    Button,
    InputGroup
} from "react-bootstrap";
import {
    Eye,
    EyeSlash
} from "react-bootstrap-icons";
import { z } from "zod";
import { redirect, RedirectType } from "next/navigation";
import ToastContainerComponent from "../../components/toastContainer";

const loginSchema = z.object({
    email: z.string().trim().min(1, "email address is required").email("invalid email address"),
    password: z.string().trim().min(1, "password is required")
});

export default function LoginPage() {
    const [pending, setPending] = useState(false);
    const [errors, setErrors] = useState([] as string[]);
    const [passwordType, setPasswordType] = useState("password");

    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");

    async function login() {
        setPending(true);
        const validatedFields = loginSchema.safeParse({email, password});
        if (!validatedFields.success) {
            const newErrors = validatedFields.error?.flatten().fieldErrors ?? { email: [], password: [] };
            const existingErrors = [
                ...newErrors.email ?? [],
                ...newErrors.password ?? []
            ];
            setErrors((prev) => [...prev, ...existingErrors]);
            setPending(false);
            return;
        }
    
        const response = await fetch("/api/users/login", {
            method: 'POST',
            body: JSON.stringify({
                email: email,
                password: password
            })
        });
        if (!response.ok) {
            const errorResponse = JSON.parse(await response.text());
            setPending(false);
            setErrors((prev) => [...prev, errorResponse.error]);
            return;
        }
    
        try
        {
            await GetMe();
            localStorage.setItem("e", (Date.now() + 86400000).toString());
        }
        catch(e)
        {
            setErrors((prev) => [...prev, (e as Error).message]);
            setPending(false);
            return;
        }
        finally
        {
            setPending(false);
        }
        redirect("/dashboard", RedirectType.replace);
    }

    return (
        <main>
            <div style={{width: "33%", alignItems: "center", textAlign: "center"}} className="mx-auto mt-40">
                <h1 className="mb-4" style={{fontSize: 24}}>Login</h1>
                <Form>
                    <FloatingLabel controlId="email" label="Email Address" className="mb-3">
                        <FormControl type="email" name="email" disabled={pending} value={email} onChange={(e) => setEmail(e.target.value)} />
                    </FloatingLabel>
                    <InputGroup className="mb-3">
                        <FloatingLabel controlId="password" label="Password">
                            <FormControl type={passwordType} name="password" disabled={pending} value={password} onChange={(e) => setPassword(e.target.value)} />
                        </FloatingLabel>
                        <InputGroup.Text onClick={() => pending ? null : setPasswordType(passwordType == "text" ? "password" : "text")} style={{cursor: "pointer"}}>
                            {passwordType == "text" ? <EyeSlash /> : <Eye />}
                        </InputGroup.Text>
                    </InputGroup>
                    <Button onClick={login} variant="outline-primary" disabled={pending}>Login</Button>
                </Form>
            </div>

            <ToastContainerComponent
                errors={errors}
                successMessages={[]}
                setErrors={setErrors}
                setSuccessToastMessages={()=>{}} />
        </main>
    )
}
