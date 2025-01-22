"use client";

import type { HostMachine } from "../interfaces/hosts";
import type { Incident } from "../interfaces/incident";
import {
    Card,
    CardBody,
    CardText,
    CardSubtitle,
    ListGroup,
    ListGroupItem,
    CardHeader,
    CardFooter,
    CardLink
} from "react-bootstrap";

export default function IncidentCard({ incident }: Readonly<{ incident: Incident }>) {
    return (
        <Card>
            <CardBody>
                <CardSubtitle>Created at {new Date(incident.createdAt).toLocaleString()}</CardSubtitle>
                <CardText>{incident.summary}</CardText>
                <CardHeader>Hosts Affected</CardHeader>
                <ListGroup>
                    {incident.hostsAffected.map((host: HostMachine) => (
                        <ListGroupItem key={`${incident.uuid}-${host.uuid}`}>
                            <CardLink href={`/hosts/${host.uuid}`}>{host.hostname}</CardLink>
                        </ListGroupItem>
                    ))}
                </ListGroup>
            </CardBody>
            <CardFooter>
                { incident.resolvedAt ? `Resolved at ${incident.resolvedAt} by ${incident.resolvedBy?.name}` : "Unresolved" }
            </CardFooter>
        </Card>
    );
}