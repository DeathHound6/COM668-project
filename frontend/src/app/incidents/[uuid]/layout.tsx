import { Metadata } from "next";

export const metadata: Metadata = {
    title: "Incident Details",
    openGraph: {
        title: "Incident Details",
    }
};

export default function PageLayout({ children }: Readonly<{ children: React.ReactNode }>) {
    return (
        <div>
            {children}
        </div>
    );
}