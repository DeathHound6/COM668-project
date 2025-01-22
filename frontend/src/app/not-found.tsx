import Link from "next/link";

export default function NotFound() {
    return (
        <main style={{textAlign: "center", fontSize: 20}} className="mt-40">
            <h1 style={{ fontSize: 40 }}>Not Found</h1>
            <br />
            <p>It seems you&apos;ve gotten lost<br />
            There is nothing here</p>
            <br />
            <Link href="/dashboard" style={{textDecoration: "underline"}}>Go back to the dashboard</Link>
        </main>
    );
}