"use client";

import { formatDate } from "../actions/api";
import type { HostMachine, Incident } from "../interfaces";
import {
    Card,
    CardBody,
    ListGroup,
    ListGroupItem,
    CardHeader,
    CardFooter,
    CardText,
    Button
} from "react-bootstrap";

export default function IncidentCard(
    { incident }:
    { incident: Incident }
) {
    const hostsAffectedLimit = 5;
    return (
        <Card>
            <CardHeader>{incident.summary}</CardHeader>
            <CardBody>
                <CardText>{incident.description}</CardText>
                <h4 className="mt-2 mb-1 underline">Affected Servers</h4>
                <ListGroup>
                    {
                        // limit to the first 5 hosts
                        incident.hostsAffected.slice(0, hostsAffectedLimit).map((host: HostMachine) => (
                            <ListGroupItem key={`${incident.uuid}-${host.uuid}`} href={`/hosts/${host.uuid}`} active={false} target="_blank" action>{host.hostname}</ListGroupItem>
                        ))
                    }
                    {
                        // if there are more than 5 hosts, show a "+n more" item
                        incident.hostsAffected.length > hostsAffectedLimit && (
                            <ListGroupItem key={`${incident.uuid}-more`}>{`+${incident.hostsAffected.length - hostsAffectedLimit} more`}</ListGroupItem>
                        )
                    }
                </ListGroup>
                <Button href={`/incidents/${incident.uuid}`} className="mt-3">View Incident</Button>
            </CardBody>
            <CardFooter>
                <CardText>
                    Opened at {formatDate(new Date(incident.createdAt))}<br />
                    {
                        incident.resolvedAt
                            ? `Resolved at ${formatDate(new Date(incident.resolvedAt))} by ${incident.resolvedBy?.name}`
                            : "Unresolved"
                    }
                </CardText>
            </CardFooter>
        </Card>
    )
}