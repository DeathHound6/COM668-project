"use client";

import { GetMe } from "../../../actions/users";
import { formatDate } from "../../../actions/api";
import { DeleteComment, GetIncident, PostComment, UpdateIncident } from "../../../actions/incidents";
import { type IncidentComment, type Incident, type User, type HostMachine, type Team, APIError } from "../../../interfaces";
import { useEffect, useState } from "react";
import {
    Button,
    Card,
    CardBody,
    CardFooter,
    CardHeader,
    CardLink,
    Col,
    FloatingLabel,
    FormCheck,
    FormControl,
    FormSelect,
    InputGroup,
    ListGroup,
    ListGroupItem,
    OverlayTrigger,
    Row,
    Spinner,
    Tooltip
} from "react-bootstrap";
import { Trash } from "react-bootstrap-icons";
import { GetTeams } from "../../../actions/teams";
import { GetHosts } from "../../../actions/hosts";
import ToastContainerComponent from "../../../components/toastContainer";

export default function IncidentPage({ params }: { params: Promise<{ uuid: string }> }) {
    const [pending, setPending] = useState(false);
    const [loaded, setLoaded] = useState(false);
    const [incident, setIncident] = useState(undefined as Incident | undefined);
    const [user, setUser] = useState(undefined as User | undefined);

    const [summary, setSummary] = useState("");
    const [description, setDescription] = useState("");
    const [teamID, setTeamID] = useState("");
    const [hostID, setHostID] = useState("");
    const [comment, setComment] = useState("");
    const [resolved, setResolved] = useState(false);
    const [comments, setComments] = useState([] as IncidentComment[]);

    const [teams, setTeams] = useState([] as Team[]);
    const [hosts, setHosts] = useState([] as HostMachine[]);

    const [errors, setErrors] = useState([] as string[]);
    const [successMessages, setSuccessMessages] = useState([] as string[]);

    function handleError(err: APIError, include404 = true) {
        const statuses = [400, 500];
        if (include404)
            statuses.push(404);
        if (statuses.includes(err.status))
            setErrors((prev) => [...prev, err.message]);
        setLoaded(true);
    }

    useEffect(() => {
        setPending(true);
        async function fetchData() {
            setLoaded(false);
            setPending(true);

            const userResponse = await GetMe().catch((e) => handleError(e, false));
            if (!userResponse)
                return;
            setUser(userResponse);

            setLoaded(false);
            const incidentResponse = await GetIncident({ uuid: (await params).uuid }).catch((e) => handleError(e, false));
            setPending(false);
            setLoaded(true);
            if (!incidentResponse)
                return;
            setIncident(incidentResponse);
            setResolved(incidentResponse.resolvedAt != undefined);
            setSummary(incidentResponse.summary);
            setDescription(incidentResponse.description);
            setComments(incidentResponse.comments);

            setLoaded(false);
            setPending(true);
            const teamsResponse = await GetTeams({ pageSize: 1000 }).catch(handleError);
            setPending(false);
            setLoaded(true);
            if (!teamsResponse)
                return;
            const newTeams = teamsResponse.data.filter((team: Team) => !incidentResponse.resolutionTeams.map((t: Team) => t.uuid).includes(team.uuid));
            setTeamID(newTeams.length > 0 ? newTeams[0].uuid : "");
            setTeams(teamsResponse.data);

            setLoaded(false);
            setPending(true);
            const hostsResponse = await GetHosts({ pageSize: 1000 }).catch(handleError);
            setPending(false);
            setLoaded(true);
            if (!hostsResponse)
                return;
            const newHosts = hostsResponse.data.filter((host: HostMachine) => !incidentResponse.hostsAffected.map((h: HostMachine) => h.uuid).includes(host.uuid));
            setHostID(newHosts.length > 0 ? newHosts[0].uuid : "");
            setHosts(hostsResponse.data);

            setLoaded(true);
            setPending(false);
        }
        fetchData();
// eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    async function deleteComment(index: number) {
        if (pending)
            return;
        // this shouldnt hit, just here to ensure typescript is happy
        if (incident == undefined) {
            setErrors((prev) => [...prev, "Incident not found"]);
            return;
        }
        const comment = comments[index];
        if (comment == undefined) {
            setErrors((prev) => [...prev, "Comment not found"]);
            return;
        }
        const deleteResponse = await DeleteComment({ incidentUUID: incident.uuid, commentUUID: comment.uuid }).catch(handleError);
        if (deleteResponse != undefined)
            return;
        setSuccessMessages((prev) => [...prev, "Comment deleted successfully"]);
        setComments(comments.filter((c: IncidentComment) => c.uuid != comment.uuid));
    }

    async function postComment() {
        if (pending)
            return;
        // this also shouldnt hit, just here to ensure typescript is happy
        if (incident == undefined) {
            setErrors((prev) => [...prev, "Incident not found"]);
            return;
        }
        // actual validity checks
        if (comment.length > 200 || comment.length == 0) {
            setErrors((prev) => [...prev, "Comment must be between 1 and 200 characters"]);
            return;
        }
        const postResponse = await PostComment({ uuid: incident.uuid, comment }).catch(handleError);
        if (postResponse == undefined)
            return;
        setSuccessMessages((prev) => [...prev, "Comment posted successfully"]);
        setComments((prev) => [{ uuid: postResponse, comment, commentedBy: user as User, commentedAt: new Date().toISOString() } as IncidentComment, ...prev]);
        setComment("");
    }

    async function updateIncident() {
        if (pending)
            return;
        // this shouldnt hit, just here to ensure typescript is happy
        if (incident == undefined) {
            setErrors((prev) => [...prev, "Incident not found"]);
            return;
        }
        // actual validity checks
        if (summary.length > 100 || summary.length == 0) {
            setErrors((prev) => [...prev, "Summary must be between 1 and 100 characters"]);
            return;
        }
        if (description.length > 500 || description.length == 0) {
            setErrors((prev) => [...prev, "Description must be between 1 and 500 characters"]);
            return;
        }
        const updated = {
            summary,
            description,
            resolutionTeams: incident.resolutionTeams.length > 0 ? [...incident.resolutionTeams.map((team: Team) => team.uuid)] : [],
            hostsAffected: incident.hostsAffected.length > 0 ? [...incident.hostsAffected.map((host: HostMachine) => host.uuid)] : [],
            resolved
        };
        const updateResponse = await UpdateIncident({ uuid: incident.uuid, incident: updated }).catch(handleError);
        if (updateResponse != undefined)
            return;
        setSuccessMessages((prev) => [...prev, "Incident updated successfully"]);
        setIncident({
            uuid: incident.uuid,
            summary,
            description,
            resolutionTeams: teams.filter((team: Team) => updated.resolutionTeams.includes(team.uuid)),
            hostsAffected: hosts.filter((host: HostMachine) => updated.hostsAffected.includes(host.uuid)),
            createdAt: incident.createdAt,
            resolvedAt: resolved ? new Date().toISOString() : undefined,
            resolvedBy: resolved ? user : undefined,
            comments
        } as Incident);
    }

    function IncidentCommentCard({ comment, index }: { comment: IncidentComment, index: number }) {
        // this shouldnt hit, just here to ensure typescript is happy
        if (user == undefined)
            return (<div></div>);
        return (
            <Card className="mb-3">
                <CardHeader>
                    <Row>
                        <Col className="ms-4">{comment.commentedBy.name}<span style={{color: "gray"}}>{comment.commentedBy.uuid == user.uuid ? " (You)": ""}</span></Col>
                        <Col xs={1}>
                            {
                                (user.admin || comment.commentedBy.uuid == user.uuid) && (
                                    <OverlayTrigger overlay={<Tooltip style={{color: pending ? "grey" : "red", cursor: "pointer"}}>Delete Comment</Tooltip>}>
                                        <Trash style={{color: "red", cursor: "pointer"}} onClick={() => pending ? null : deleteComment(index)} />
                                    </OverlayTrigger>
                                )
                            }
                        </Col>
                    </Row>
                </CardHeader>
                <CardBody>{comment.comment}</CardBody>
                <CardFooter>{formatDate(new Date(comment.commentedAt))}</CardFooter>
            </Card>
        )
    }

    return (
        <main style={{textAlign: "center"}}>
            {
                incident != undefined && (
                    <Row className="mt-3 mx-2 ms-auto max-w-96">
                        <Col style={{textAlign: "right"}}>
                            { /* <FormCheck type="switch" label="Notify Teams" inline /> */}
                            <FormCheck type="switch" label="Resolved" inline checked={resolved} onChange={(e) => setResolved(e.target.checked)} />
                            <Button onClick={() => updateIncident()} className="ms-2" disabled={pending}>Update Incident</Button>
                        </Col>
                    </Row>
            )}

            <div className="mt-2">
                {
                    !loaded
                        ? (<Spinner animation="border" role="status" className="my-auto mx-auto" />)
                        : (
                            <div>
                                {
                                    incident == undefined
                                        ? (
                                            <div className="text-center mt-5">
                                                <h1 className="underline" style={{fontSize: 32}}>Incident not found</h1>
                                                <p className="mt-2">There is no incident with the given identifier</p>
                                                <Button className="mt-2" href="/dashboard">Go to Dashboard</Button>
                                                <br />
                                                <Button className="mt-2" href="/history">View Incident History</Button>
                                            </div>
                                        )
                                        : (
                                            <Row className="mx-4">
                                                { /* Incident Information */ }
                                                <Col className="text-center" xs={3}>
                                                    <h1 className="underline mb-2" style={{fontSize: 24}}>Incident Details</h1>
                                                    <FloatingLabel controlId="summary" label={`Summary (${summary.length}/100 Characters)`}>
                                                        <FormControl value={summary} onChange={(e) => setSummary(e.target.value)} className="mt-2" isInvalid={summary.length > 100} />
                                                        <FormControl.Feedback type="invalid" tooltip>Summary must be between 1 and 100 characters</FormControl.Feedback>
                                                    </FloatingLabel>
                                                    <FloatingLabel controlId="description" label={`Description (${description.length}/500 Characters)`}>
                                                        <FormControl value={description} as="textarea" rows={4} onChange={(e) => setDescription(e.target.value)} className="my-2" isInvalid={description.length > 500} />
                                                        <FormControl.Feedback type="invalid" tooltip>Description must be between 1 and 500 characters</FormControl.Feedback>
                                                    </FloatingLabel>
                                                    <p>Opened {formatDate(new Date(incident.createdAt))}</p>
                                                    <p>
                                                        {
                                                            incident.resolvedAt
                                                                ? `Incident was resolved at ${formatDate(new Date(incident.resolvedAt))} by ${incident.resolvedBy?.name}`
                                                                : "Incident is Unresolved"
                                                        }
                                                    </p>
                                                    <h1 className="underline mt-3" style={{fontSize: 20}}>Teams required to resolve</h1>
                                                    {
                                                        incident.resolutionTeams.length == 0
                                                            ? (<div>
                                                                { /* This shouldn't ever show, but better to show something than nothing if a processor bug occurs */ }
                                                                <p className="mt-2">This incident currently does not require any teams to resolve</p>
                                                            </div>)
                                                            : (
                                                                <ListGroup className="mt-2">
                                                                    {
                                                                        incident.resolutionTeams.map((team: Team, index: number) => (
                                                                            <ListGroupItem key={`team-${index}`}>
                                                                                <Row>
                                                                                    <Col>{team.name}</Col>
                                                                                    <Col xs={2}>
                                                                                        <OverlayTrigger overlay={<Tooltip>Remove Team</Tooltip>}>
                                                                                            <Trash style={{color: "red", cursor: "pointer"}} onClick={() => {
                                                                                                incident.resolutionTeams = incident.resolutionTeams.filter((t: Team) => t.uuid != team.uuid);
                                                                                                const newTeams = teams.filter((team: Team) => !incident.resolutionTeams.map((t: Team) => t.uuid).includes(team.uuid));
                                                                                                setTeamID(newTeams.length > 0 ? newTeams[0].uuid : "");
                                                                                            }} />
                                                                                        </OverlayTrigger>
                                                                                    </Col>
                                                                                </Row>
                                                                            </ListGroupItem>
                                                                        ))
                                                                    }
                                                                </ListGroup>
                                                            )
                                                    }
                                                    <InputGroup className="mt-2">
                                                        <FloatingLabel controlId="team" label="Add a Team">
                                                            <FormSelect value={teamID} onChange={(e) => setTeamID(e.target.value)}>
                                                                {
                                                                    teams.filter((team: Team) => !incident.resolutionTeams.map((t: Team) => t.uuid).includes(team.uuid))
                                                                    .map((team: Team, index: number) => (
                                                                        <option key={`team-${index}`} value={team.uuid}>{team.name}</option>
                                                                    ))
                                                                }
                                                            </FormSelect>
                                                        </FloatingLabel>
                                                        <Button onClick={() => {
                                                            incident.resolutionTeams.push(teams.find((team: Team) => team.uuid == teamID) as Team);
                                                            const newTeams = teams.filter((team: Team) => !incident.resolutionTeams.map((t: Team) => t.uuid).includes(team.uuid));
                                                            setTeamID(newTeams.length > 0 ? newTeams[0].uuid : "");
                                                        }} disabled={teamID == ""}>Add</Button>
                                                    </InputGroup>
                                                </Col>

                                                { /* List of affected servers */ }
                                                <Col>
                                                    <h1 className="underline mb-2" style={{fontSize: 24}}>Affected Servers</h1>
                                                    {
                                                        incident.hostsAffected.length == 0
                                                            ? (<div>
                                                                {/* This shouldn't ever show, but better to show something than nothing if a processor bug occurs */}
                                                                <p className="mt-2">This incident currently does not impact any servers</p>
                                                            </div>)
                                                            : incident.hostsAffected.map((host: HostMachine, index: number) => (
                                                                <Card key={`host-${index}`} className="mt-2">
                                                                    <CardHeader>
                                                                        <Row>
                                                                            <Col>{host.hostname}</Col>
                                                                            <Col xs={1}>
                                                                                <OverlayTrigger overlay={<Tooltip>Remove Server</Tooltip>}>
                                                                                    <Trash style={{color: "red", cursor: "pointer"}} onClick={() => {
                                                                                            incident.hostsAffected.splice(index, 1);
                                                                                            const newHosts = hosts.filter((host: HostMachine) => !incident.hostsAffected.map((h: HostMachine) => h.uuid).includes(host.uuid));
                                                                                            setHostID(newHosts.length > 0 ? newHosts[0].uuid : "");
                                                                                        }} />
                                                                                </OverlayTrigger>
                                                                            </Col>
                                                                        </Row>
                                                                    </CardHeader>
                                                                    <CardLink href={`/hosts/${host.uuid}`} target="_blank">
                                                                        <CardBody>
                                                                            <FloatingLabel controlId="ip4" label="IPv4 Address" className="max-w-96 mx-auto">
                                                                                <FormControl readOnly value={host.ip4 ?? "Not Assigned"} disabled className="cursor-pointer" />
                                                                            </FloatingLabel>
                                                                            <FloatingLabel controlId="ip6" label="IPv6 Address" className="mt-2 max-w-96 mx-auto">
                                                                                <FormControl readOnly value={host.ip6 ?? "Not Assigned"} disabled className="cursor-pointer" />
                                                                            </FloatingLabel>
                                                                            <FloatingLabel controlId="os" label="Operating System" className="mt-2 max-w-96 mx-auto">
                                                                                <FormControl readOnly value={host.os} disabled className="cursor-pointer" />
                                                                            </FloatingLabel>
                                                                        </CardBody>
                                                                        <CardFooter>
                                                                            Managed by the {host.team.name} team
                                                                        </CardFooter>
                                                                    </CardLink>
                                                                </Card>
                                                            ))
                                                    }
                                                    <InputGroup className="mt-2">
                                                        <FloatingLabel controlId="host" label="Add a Server">
                                                            <FormSelect value={hostID} onChange={(e) => setHostID(e.target.value)}>
                                                                {
                                                                    hosts.filter((host: HostMachine) => !incident.hostsAffected.map((h: HostMachine) => h.uuid).includes(host.uuid))
                                                                    .map((host: HostMachine, index: number) => (
                                                                        <option key={`host-${index}`} value={host.uuid}>{host.hostname}</option>
                                                                    ))
                                                                }
                                                            </FormSelect>
                                                        </FloatingLabel>
                                                        <Button onClick={() => {
                                                            incident.hostsAffected.push(hosts.find((host: HostMachine) => host.uuid == hostID) as HostMachine);
                                                            const newHosts = hosts.filter((host: HostMachine) => !incident.hostsAffected.map((h: HostMachine) => h.uuid).includes(host.uuid));
                                                            setHostID(newHosts.length > 0 ? newHosts[0].uuid : "");
                                                        }} disabled={hostID == ""}>Add</Button>
                                                    </InputGroup>
                                                </Col>

                                                { /* Incident Comments */ }
                                                <Col className="text-center">
                                                    <h1 className="underline mb-2" style={{fontSize: 24}}>Comments</h1>
                                                    {
                                                        comments.length == 0
                                                            ? (<div>
                                                                <p className="mt-2">No comments have been posted on this incident</p>
                                                            </div>)
                                                            : comments.map((comment: IncidentComment, index: number) => (
                                                                <IncidentCommentCard key={comment.uuid} comment={comment} index={index} />
                                                            ))
                                                    }
                                                    <Card className="mt-4">
                                                        <CardHeader>Leave a Comment</CardHeader>
                                                        <CardBody>
                                                            <FloatingLabel controlId="comment" label={`Comment (${comment.length}/200 Characters)`}>
                                                                <FormControl as="textarea" rows={3} value={comment} onChange={(e) => setComment(e.target.value)} isInvalid={comment.length > 200} />
                                                                <FormControl.Feedback type="invalid" tooltip>Comment must be between 1 and 200 characters</FormControl.Feedback>
                                                            </FloatingLabel>
                                                            <Button className="mt-2" variant="primary" onClick={() => postComment()} disabled={pending || comment.length > 200 || comment.length == 0}>Post Comment</Button>
                                                        </CardBody>
                                                    </Card>
                                                </Col>
                                            </Row>
                                        )
                                }
                            </div>
                        )
                }
            </div>

            {/* Toasts for showing error messages */}
            <ToastContainerComponent
                errors={errors}
                successMessages={successMessages}
                setErrors={setErrors}
                setSuccessToastMessages={setSuccessMessages}
                />
        </main>
    );
}