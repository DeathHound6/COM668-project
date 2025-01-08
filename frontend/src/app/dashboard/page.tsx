"use client";

import { Suspense, useEffect, useState } from "react";

export default function DashboardPage() {
    const [incidents, setIncidents] = useState([]);

    useEffect(() => {
        fetch("/api/incidents")
            .then(res => res.json())
            .then(setIncidents);
    }, []);

    return (
        <div>
            <Suspense>
            </Suspense>
        </div>
    );
}