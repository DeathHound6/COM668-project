"use client";

import type { Incident } from "../../interfaces/incident";
import IncidentCard from "../../components/incidentCard";
import { Suspense, useEffect, useState } from "react";
import { Spinner } from "react-bootstrap";

export default function DashboardPage() {
    const [incidents, setIncidents] = useState([] as Incident[]);

    const [_, setAPIError] = useState(null as string | null);

    useEffect(() => {
        fetch("/api/incidents?resolved=false")
            .then(
                async(res) => {
                    const data = await res.json();
                    if (!res.ok)
                        return setAPIError(data.error);
                    setIncidents(data.data);
                },
                (err) => {
                    setAPIError((err as Error).message);
                }
            );
    }, []);

    return (
        <main>
            <Suspense fallback={<Spinner role="status" animation="border" />}>
                {
                    incidents.length < 0 && incidents.map((incident: Incident) => (
                        <IncidentCard incident={incident} key={incident.uuid} />)
                    )
                }
            </Suspense>
        </main>
    );
}