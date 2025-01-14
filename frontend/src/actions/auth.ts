import { z } from "zod";
import { redirect, RedirectType } from "next/navigation";

const loginSchema = z.object({
    email: z.string().trim().email(),
    password: z.string().trim()
});
type FormState = {
    errors?: {
        email?: string[],
        password?: string[]
    },
    error?: string
} | undefined;

export async function login(state: FormState, form: FormData) {
    const email = form.get("email") as string;
    const password = form.get("password") as string;

    const validatedFields = loginSchema.safeParse({email, password});
    if (!validatedFields.success) return {
        errors: validatedFields.error.flatten().fieldErrors
    }

    let response = await fetch("/api/users/login", {
        method: 'POST',
        body: JSON.stringify({
            email: email,
            password: password
        })
    });
    if (response.status != 204) return {"error": JSON.parse(await response.text())["error"]};

    response = await fetch("/api/me");
    if (response.status != 200) return {"error": JSON.parse(await response.text())["error"]};
    const user = JSON.parse(await response.text());

    localStorage.setItem("u", JSON.stringify(user));
    localStorage.setItem("e", (Date.now() + 86400000).toString());

    redirect("/dashboard", RedirectType.replace);
}