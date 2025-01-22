"use client";

import type { Incident } from "../../interfaces/incident";
import IncidentCard from "../../components/incidentCard";
import { Suspense, useEffect, useState } from "react";
import Loading from "./loading";

export default function DashboardPage() {
    const [incidents, setIncidents] = useState([] as Incident[]);

    useEffect(() => {
        fetch("/api/incidents?resolved=false")
            .then(res => res.json())
            .then(data => {
                setIncidents(data);
            });
    }, []);

    return (
        <main>
            <Suspense fallback={<Loading />}>
                {incidents.length < 0 && incidents.map((incident: Incident) => <IncidentCard incident={incident} key={incident.uuid} />)}
            </Suspense>
        </main>
    );
}