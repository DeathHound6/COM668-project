import { Metadata } from "next";

export const metadata: Metadata = {
    title: "Dashboard",
    openGraph: {
        title: "Dashboard",
    }
};

export default function PageLayout({ children }: Readonly<{ children: React.ReactNode }>) {
    return (
        <div>
            {children}
        </div>
    );
}