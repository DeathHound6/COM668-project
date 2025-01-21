import { z } from "zod";
import { redirect, RedirectType } from "next/navigation";

const loginSchema = z.object({
    email: z.string().trim().min(1, "email address is required").email("invalid email address"),
    password: z.string().trim().min(1, "password is required")
});
export type FormState = {
    errors: {
        email?: string[] | undefined,
        password?: string[] | undefined
    },
    error?: string | undefined
};

export async function login(state: FormState, form: FormData) {
    const email = form.get("email") as string;
    const password = form.get("password") as string;

    const validatedFields = loginSchema.safeParse({email, password});
    if (!validatedFields.success) return {
        errors: validatedFields.error.flatten().fieldErrors,
        error: undefined
    }

    const response = await fetch("/api/users/login", {
        method: 'POST',
        body: JSON.stringify({
            email: email,
            password: password
        })
    });
    if (response.status != 204) return {
        errors: {email: undefined, password: undefined},
        error: JSON.parse(await response.text())["error"] as string
    };

    try
    {
        await getMe();
        localStorage.setItem("e", (Date.now() + 86400000).toString());
        redirect("/dashboard", RedirectType.replace);
    }
    catch(e)
    {
        return {
            errors: {email: undefined, password: undefined},
            error: (e as Error).message
        };
    }
}

export async function getMe() {
    const response = await fetch("/api/me");
    if (response.status != 200) throw new Error(JSON.parse(await response.text())["error"]);
    const userinfo = JSON.parse(await response.text());
    localStorage.setItem("u", JSON.stringify(userinfo));
    return userinfo;
}
