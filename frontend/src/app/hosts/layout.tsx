import { Metadata } from "next";

export const metadata: Metadata = {
    title: "Hosts",
    openGraph: {
        title: "Hosts",
    }
};

export default function PageLayout({ children }: Readonly<{ children: React.ReactNode }>) {
    return (
        <div>
            {children}
        </div>
    );
}