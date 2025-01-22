"use client";

import { login, FormState } from "../../actions/auth";
import { startTransition, useActionState, useEffect, useState } from "react";
import {
    FloatingLabel,
    FormControl,
    Form,
    Button,
    Toast,
    ToastHeader,
    ToastBody,
    ToastContainer,
    InputGroup
} from "react-bootstrap";
import {
    Eye,
    EyeSlash
} from "react-bootstrap-icons";

export default function LoginPage() {
    const [state, action, pending] = useActionState<FormState, FormData>(login, { error: undefined, errors: { email: undefined, password: undefined } });
    const [showAPIError, setShowAPIError] = useState(false);
    const [showEmailError, setShowEmailError] = useState([] as boolean[]);
    const [showPasswordError, setShowPasswordError] = useState([] as boolean[]);
    const [passwordType, setPasswordType] = useState("password");

    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");

    useEffect(() => {
        setShowAPIError(state.error != undefined);
        setShowEmailError(state.errors.email == undefined ? [] : new Array(state.errors.email?.length || 0).fill(true));
        setShowPasswordError(state.errors.password == undefined ? [] : new Array(state.errors.password?.length || 0).fill(true));
    }, [state.error, state.errors.email, state.errors.password]);

    function onFormSubmit(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault();
        const form = new FormData();
        form.append("email", email);
        form.append("password", password);
        startTransition(() => action(form));
    }

    function onCloseToast(index: number, type: string) {
        const errors = type == "email" ? [...showEmailError] : [...showPasswordError];
        errors[index] = false;
        if (type == "email")
            setShowEmailError(errors);
        else
            setShowPasswordError(errors);
    }

    return (
        <main>
            <div style={{width: "33%", alignItems: "center", textAlign: "center"}} className="mx-auto mt-40">
                <h1 className="mb-4" style={{fontSize: 24}}>Login</h1>
                <Form onSubmit={onFormSubmit}>
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
                    <Button type="submit" variant="outline-primary" disabled={pending}>Login</Button>
                </Form>
            </div>

            <ToastContainer position="bottom-end" className="p-3">
                { state.errors.email?.map((error: string, index: number) => (
                    showEmailError[index] && (
                        <Toast bg="danger" onClose={() => onCloseToast(index, "email")} key={`email-${index}`}>
                            <ToastHeader>Error</ToastHeader>
                            <ToastBody>{error}</ToastBody>
                        </Toast>
                    ))
                )}
                { state.errors.password?.map((error: string, index: number) => (
                    showPasswordError[index] && (
                        <Toast bg="danger" onClose={() => onCloseToast(index, "password")} key={`pw-${index}`}>
                            <ToastHeader>Error</ToastHeader>
                            <ToastBody>{error}</ToastBody>
                        </Toast>
                    ))
                )}
                { showAPIError && (
                    <Toast bg="danger" onClose={() => { setShowAPIError(false); }} key={"error"}>
                        <ToastHeader>Error</ToastHeader>
                        <ToastBody>{state.error}</ToastBody>
                    </Toast>
                )}
            </ToastContainer>
        </main>
    )
}
