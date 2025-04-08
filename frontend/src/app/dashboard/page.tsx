"use client";

import type { APIError, Incident } from "../../interfaces";
import IncidentCard from "../../components/incidentCard";
import { useEffect, useState } from "react";
import {
    Button,
    Col,
    FormCheck,
    Row,
    Spinner
} from "react-bootstrap";
import { GetIncidents } from "../../actions/incidents";
import ToastContainerComponent from "../../components/toastContainer";
import Paginator from "../../components/paginator";

export default function DashboardPage() {
    const [loaded, setLoaded] = useState(false);
    const [incidents, setIncidents] = useState([] as Incident[]);

    const [page, setPage] = useState(1);
    const [maxPage, setMaxPage] = useState(1);
    const [myTeams, setMyTeams] = useState(false);

    const [errors, setErrors] = useState([] as string[]);

    function handleError(err: APIError) {
        if ([400, 404, 500].includes(err.status))
            setErrors((prev) => [...prev, err.message]);
        setLoaded(true);
    }

    useEffect(() => {
        setLoaded(false);
        async function fetchData() {
            const params = {
                "myTeams": myTeams,
                "resolved": false,
                "page": page
            };
            const incidents = await GetIncidents({ params }).catch(handleError);
            if (!incidents) return;
            setMaxPage(incidents.meta.pages);
            setIncidents(incidents.data);
            setLoaded(true);
        }
        fetchData();
    }, [myTeams, page]);

    return (
        <main>
            <div className="mx-5 mt-3" style={{textAlign: "center"}}>
                <div className="mx-auto" style={{fontSize: 20}}><b>Active Incidents</b></div>
                <div style={{textAlign: "justify"}} className="mb-2">
                    <FormCheck inline label="My teams only" type="switch" checked={myTeams} onChange={(e) => setMyTeams(e.target.checked)} />
                </div>
                {
                    !loaded
                        ? (<Spinner role="status" animation="border" className="my-auto mx-auto" />)
                        : (
                            <>
                                {
                                    incidents.length == 0
                                        ? (
                                            <div className="mx-auto mt-5">
                                                <h1 style={{fontSize: 40}}><b>No Incidents</b></h1>
                                                <br />
                                                <p style={{fontSize: 20}}>There are currently no unresolved incidents</p>
                                                <p style={{fontSize: 20}}>Please ensure that applied filters are correct</p>
                                                <Button href="/history" className="mt-4">View Incident History</Button>
                                            </div>
                                        )
                                        : (
                                            <Row xs={1} md={2} lg={4}>
                                                {
                                                    incidents.map((incident: Incident) => (
                                                        <Col key={incident.uuid}>
                                                            <IncidentCard incident={incident} />
                                                        </Col>
                                                    ))
                                                }
                                            </Row>
                                        )
                                }
                                <Paginator page={page} maxPage={maxPage} setPage={setPage} />
                            </>
                        )
                }
            </div>

            <ToastContainerComponent
                errors={errors}
                setErrors={setErrors}
                successMessages={[]}
                setSuccessToastMessages={() => {}}
            />
        </main>
    );
}