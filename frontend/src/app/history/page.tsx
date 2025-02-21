"use client";

import type { APIError, Incident } from "../../interfaces";
import IncidentCard from "../../components/incidentCard";
import { useEffect, useState } from "react";
import {
    Button,
    Col,
    FormCheck,
    Pagination,
    Row,
    Spinner
} from "react-bootstrap";
import { GetIncidents } from "../../actions/incidents";
import ToastContainerComponent from "../../components/toastContainer";

export default function HistoryPage() {
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
                "resolved": true,
                "page": page
            };
            GetIncidents({ params })
                .then(
                    (data) => {
                        setMaxPage(data.meta.pages);
                        setIncidents(data.data);
                        setLoaded(true);
                    },
                    handleError
                );

        }
        fetchData();
    }, [myTeams, page]);

    return (
        <main>
            <div className="mx-5 mt-3" style={{textAlign: "center"}}>
                <div className="mx-auto" style={{fontSize: 20}}><b>Incident History</b></div>
                <div style={{textAlign: "justify"}} className="mb-2">
                    <FormCheck inline label="My teams only" type="switch" checked={myTeams} onChange={(e) => setMyTeams(e.target.checked)} />
                </div>
                {
                    !loaded
                        ? (<Spinner role="status" animation="border" className="my-auto mx-auto" />)
                        : (
                            <div>
                                {
                                    incidents.length == 0
                                        ? (
                                            <div className="mx-auto mt-5">
                                                <h1 style={{fontSize: 40}}><b>No Incidents</b></h1>
                                                <br />
                                                <p style={{fontSize: 20}}>There are no past incidents</p>
                                                <p style={{fontSize: 20}}>Please ensure that applied filters are correct</p>
                                                <Button href="/dashboard" className="mt-4">View Unresolved Incidents</Button>
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
                                <Pagination className="mt-3 mx-auto max-w-40">
                                    <Pagination.First onClick={() => setPage(1)} disabled={maxPage == 0} />
                                    <Pagination.Prev onClick={() => setPage((prev) => prev - 1)} disabled={page == 1} />
                                    <Pagination.Ellipsis hidden={page < 3} />

                                    <Pagination.Item hidden={maxPage <= 1} active={page == 1}>{page == 1 ? 1 : page - 1}</Pagination.Item>
                                    <Pagination.Item active={(page != 1 && page != maxPage) || (page == 1 && maxPage < 3)}>{page}</Pagination.Item>
                                    <Pagination.Item hidden={maxPage < 3} active={page == maxPage}>{page == maxPage ? maxPage : page + 1}</Pagination.Item>

                                    <Pagination.Ellipsis hidden={page > maxPage - 3} />
                                    <Pagination.Next onClick={() => setPage((prev) => prev + 1)} disabled={page == maxPage || maxPage == 0} />
                                    <Pagination.Last onClick={() => setPage(maxPage)} disabled={maxPage == 0} />
                                </Pagination>
                            </div>
                        )
                }
            </div>

            <ToastContainerComponent
                errors={errors}
                successMessages={[]}
                setErrors={setErrors}
                setSuccessToastMessages={() => {}} />
        </main>
    );
}