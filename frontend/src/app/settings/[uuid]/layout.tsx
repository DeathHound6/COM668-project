import { Metadata } from "next";

export const metadata: Metadata = {
    title: "Setting Details",
    openGraph: {
        title: "Setting Details",
    }
};

export default function PageLayout({ children }: Readonly<{ children: React.ReactNode }>) {
    return (
        <div>
            {children}
        </div>
    );
}