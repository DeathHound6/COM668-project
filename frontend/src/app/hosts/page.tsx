"use client";

import { Button, Card, CardBody, Col, Row, Spinner } from "react-bootstrap";
import { HostMachine } from "../../interfaces/hosts";
import { Suspense, useEffect, useState } from "react";

export default function HostsPage() {
    const [hosts, setHosts] = useState([] as HostMachine[]);

    const [showCreateModal, setShowCreateModal] = useState(false);

    const [apiError, setAPIError] = useState(null as string | null);

    useEffect(() => {
        fetch("/api/hosts")
        .then(
            async(res) => {
                const data = await res.json();
                if (!res.ok)
                    return setAPIError(data.message);
                console.log(data);
                setHosts(data.data);
            },
            (err) => {
                setAPIError((err as Error).message);
            }
        )
    }, []);

    return (
        <div className="m-2">
            <Row className="mt-3">
                <Col style={{textAlign: "right"}}>
                    <Button variant="secondary" onClick={() => setShowCreateModal(true)}>Add Host</Button>
                </Col>
            </Row>

            <Suspense fallback={<Spinner animation="border" role="status" />}>
                <Row xs={2} md={4}>
                    {
                        hosts.length > 0 && hosts.map((host) => (
                            <Col key={host.uuid}>
                                <Card>
                                    <CardBody>
                                        <Card.Title>{host.hostname}</Card.Title>
                                        <Card.Text>
                                            {host.ip4}
                                        </Card.Text>
                                    </CardBody>
                                </Card>
                            </Col>
                        ))
                    }
                </Row>
            </Suspense>
        </div>
    );
}