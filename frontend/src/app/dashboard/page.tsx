"use client";

import IncidentCard from "../../components/incidentCard";
import { Suspense, useEffect, useState } from "react";
import Loading from "./loading";

export default function DashboardPage() {
    const [incidents, setIncidents] = useState([] as any[]);

    // useEffect(() => {
    //     fetch("/api/incidents?resolved=false")
    //         .then(res => res.json())
    //         .then(data => {
    //             setIncidents(data);
    //         });
    // }, []);

    return (
        <div>
            <Suspense fallback={<Loading />}>
                {incidents.length < 0 && incidents.map((incident: any) => <IncidentCard incident={incident} />)}
                <Loading />
            </Suspense>
        </div>
    );
}