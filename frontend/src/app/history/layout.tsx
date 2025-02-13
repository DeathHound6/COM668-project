import { Metadata } from "next";

export const metadata: Metadata = {
    title: "Incident History",
    openGraph: {
        title: "Incident History",
    }
};

export default function PageLayout({ children }: Readonly<{ children: React.ReactNode }>) {
    return (
        <div>
            {children}
        </div>
    );
}