import { Metadata } from "next";

export const metadata: Metadata = {
    title: "Login",
    openGraph: {
        title: "Login",
    }
};

export default function PageLayout({ children }: Readonly<{ children: React.ReactNode }>) {
    return (
        <div>
            {children}
        </div>
    );
}