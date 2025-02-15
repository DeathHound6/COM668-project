import { Metadata } from "next";

export const metadata: Metadata = {
    title: "Host Details",
    openGraph: {
        title: "Host Details",
    }
};

export default function PageLayout({ children }: Readonly<{ children: React.ReactNode }>) {
    return (
        <div>
            {children}
        </div>
    );
}