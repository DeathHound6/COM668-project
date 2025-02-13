import { Metadata } from "next";

export const metadata: Metadata = {
    title: "Settings",
    openGraph: {
        title: "Settings",
    }
};

export default function PageLayout({ children }: Readonly<{ children: React.ReactNode }>) {
    return (
        <div>
            {children}
        </div>
    );
}