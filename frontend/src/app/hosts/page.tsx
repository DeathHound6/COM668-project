"use client";

import type { HostMachine } from "../../interfaces/hosts";
import type { Team } from "../../interfaces/user";
import {
    Button,
    Card,
    CardBody,
    CardTitle,
    Col,
    FloatingLabel,
    FormControl,
    FormSelect,
    Modal,
    ModalBody,
    ModalHeader,
    ModalTitle,
    OverlayTrigger,
    Row,
    Spinner,
    Toast,
    ToastBody,
    ToastContainer,
    ToastHeader,
    Tooltip
} from "react-bootstrap";
import { Suspense, useEffect, useState } from "react";
import { z } from "zod";
import { Trash } from "react-bootstrap-icons";
import { GetTeams } from "../../actions/teams";
import { GetHosts, UpdateHost, DeleteHost, CreateHost, GetHost } from "../../actions/hosts";
import { APIError } from "../../interfaces/error";
import { handleUnauthorized } from "../../actions/api";

export default function HostsPage() {
    const oses = [
        "Windows",
        "Linux",
        "MacOS"
    ];
    const hostSchema = z.object({
        hostname: z.string().trim().nonempty("Hostname is required"),
        ip4: z.string().trim().ip({ version: "v4", message: "Invalid IPv4 address" }).optional(),
        ip6: z.string().trim().ip({ version: "v6", message: "Invalid IPv6 address" }).optional(),
        os: z.string().trim().nonempty("OS is required"),
        teamID: z.string().trim().nonempty("Team is required").uuid("Invalid team ID")
    });

    const [hosts, setHosts] = useState([] as HostMachine[]);
    const [teams, setTeams] = useState([] as Team[]);
    const [errors, setErrors] = useState([] as string[]);
    const [pending, setPending] = useState(true);

    const [showCreateModal, setShowCreateModal] = useState(false);

    const [showErrors, setShowErrors] = useState([] as boolean[]);
    const [showAPIError, setShowAPIError] = useState(false);
    const [apiError, setAPIError] = useState(undefined as string | undefined);

    const [showSuccessToast, setShowSuccessToast] = useState(false);
    const [successToastMessage, setSuccessToastMessage] = useState(undefined as string | undefined);

    // used for creating new host
    const [hostname, setHostname] = useState("");
    const [ip4, setIp4] = useState("");
    const [ip6, setIp6] = useState("");
    const [os, setOs] = useState("");
    const [teamID, setTeamID] = useState("");

    function handleError(err: APIError) {
        handleUnauthorized({ err });
        if ([400, 404, 500].includes(err.status))
            setAPIError(err.message);
        setPending(false);
    }

    useEffect(() => {
        setOs(oses[0]);
        setPending(true);
        async function fetchData() {
            // Fetch all teams
            let page = 1;
            const teamsResponse = await GetTeams({ page }).catch(handleError);
            if (!teamsResponse)
                return;
            const pages = teamsResponse.meta.pages;
            setTeams(teamsResponse.data);
            setTeamID(teamsResponse.data.length > 0 ? teamsResponse.data[0].uuid : "");
            while (page < pages) {
                const teamsResponse = await GetTeams({ page: page + 1 }).catch(handleError);
                if (!teamsResponse)
                    break;
                setTeams((prev) => [...prev, ...teamsResponse.data]);
                page = teamsResponse.meta.page;
            }

            const hostsResponse = await GetHosts({}).catch(handleError);
            if (!hostsResponse)
                return;
            setHosts(hostsResponse.data);
            setPending(false);
        }
        fetchData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    useEffect(() => {
        setShowAPIError(apiError != undefined);
        setShowErrors(new Array(errors.length).fill(true));
        setShowSuccessToast(successToastMessage != undefined);
    }, [apiError, errors, successToastMessage]);

    function onCloseToast(index: number) {
        const e = [...showErrors];
        if (e.length <= index) {
            setAPIError("invalid error index");
            return;
        }
        e[index] = false;
        setShowErrors(e);
    }

    function createHost() {
        setPending(true);
        const validatedFields = hostSchema.safeParse({ hostname, ip4: ip4.length ? ip4 : undefined, ip6: ip6.length ? ip6 : undefined, os, teamID });
        if (!validatedFields.success || teams.length == 0 || (ip4 == "" && ip6 == "")) {
            const newErrors = validatedFields.error?.flatten().fieldErrors ?? { hostname: [], ip4: [], ip6: [], os: [], teamID: [] };
            const existingErrors = [
                ...newErrors.hostname ?? [],
                ...newErrors.ip4 ?? [],
                ...newErrors.ip6 ?? [],
                ...newErrors.os ?? [],
                ...(teams.length == 0 ? ["no teams found"] : newErrors.teamID ?? []),
                ...(ip4 == "" && ip6 == "" ? ["at least one of IPv4 or IPv6 is required"] : [])
            ];
            setErrors(existingErrors);
            setPending(false);
            return;
        }
        CreateHost({ hostname, os, ip4, ip6, teamID }).then(
            async(uuid) => {
                const host = await GetHost({ uuid }).catch(handleError);
                setPending(false);
                if (!host)
                    return;
                setHosts((prev) => [...prev, host]);
                setShowCreateModal(false);
                setSuccessToastMessage("Host created successfully");
            },
            handleError
        );
    }

    async function deleteHost(index: number) {
        setPending(true);
        const host = hosts[index];
        DeleteHost(host.uuid).then(
            () => {
                const newHosts = [...hosts];
                newHosts.splice(index, 1);
                setHosts(newHosts);
                setPending(false);
                setSuccessToastMessage("Host deleted successfully");
            },
            handleError
        );
    }

    function updateHost(index: number) {
        setPending(true);
        const host = hosts[index];
        const hostname = host.hostname;
        const ip4 = host.ip4.length > 0 ? host.ip4 : "";
        const ip6 = host.ip6.length > 0 ? host.ip6 : "";
        const os = host.os;
        const teamID = host.team.uuid;
        const body = { hostname, ip4: ip4.length > 0 ? ip4 : undefined, ip6: ip6.length > 0 ? ip6 : undefined, os, teamID };
        const validatedFields = hostSchema.safeParse(body);
        if (!validatedFields.success) {
            const newErrors = validatedFields.error.flatten().fieldErrors;
            const existingErrors = [
                ...newErrors.hostname ?? [],
                ...newErrors.ip4 ?? [],
                ...newErrors.ip6 ?? [],
                ...newErrors.os ?? [],
                ...(teams.length == 0 ? ["no teams found"] : newErrors.teamID ?? []),
                ...(host.ip4 == "" && host.ip6 == "" ? ["at least one of IPv4 or IPv6 is required"] : [])
            ];
            setErrors(existingErrors);
            setPending(false);
            return;
        }
        UpdateHost({ uuid: host.uuid, body }).then(
            () => {
                const newHosts = [...hosts];
                setHosts(newHosts);
                setPending(false);
                setSuccessToastMessage("Host updated successfully");
            },
            handleError
        );
    }

    return (
        <div className="m-2">
            <Modal
               show={showCreateModal}
               onHide={() => setShowCreateModal(false)}
               backdrop="static"
               centered={true}
               restoreFocus={false}>
                <ModalHeader closeButton>
                    <ModalTitle>Create Host</ModalTitle>
                </ModalHeader>
                <ModalBody>
                    <FloatingLabel label="Hostname" controlId="hostname" className="mt-2">
                        <FormControl type="text" value={hostname} onChange={(e) => setHostname(e.target.value)} />
                    </FloatingLabel>
                    <FloatingLabel label="IPv4" controlId="ip4" className="mt-2">
                        <FormControl type="text" value={ip4} onChange={(e) => setIp4(e.target.value)} />
                    </FloatingLabel>
                    <FloatingLabel label="IPv6" controlId="ip6" className="mt-2">
                        <FormControl type="text" value={ip6} onChange={(e) => setIp6(e.target.value)} />
                    </FloatingLabel>
                    <FloatingLabel label="OS" controlId="os" className="mt-2">
                        <FormSelect value={os} onChange={(e) => setOs(e.target.value)}>
                            {
                                oses.map((os: string) => (
                                    <option value={os} key={os}>{os}</option>
                                ))
                            }
                        </FormSelect>
                    </FloatingLabel>
                    <FloatingLabel label="Team" controlId="team" className="mt-2">
                        <FormSelect value={teamID} onChange={(e) => setTeamID(e.target.value)}>
                            {
                                teams.length > 0
                                    ? teams.map((team: Team) => (
                                        <option value={team.uuid} key={team.uuid}>{team.name}</option>
                                    ))
                                    : <option value="">No teams found</option>
                            }
                        </FormSelect>
                    </FloatingLabel>
                    <Button variant="primary" onClick={() => createHost()} style={{textAlign: "center"}} className="mt-2" disabled={pending}>Create</Button>
                </ModalBody>
            </Modal>

            <Row className="mt-3">
                <Col style={{textAlign: "right"}}>
                    <Button variant="secondary" onClick={() => setShowCreateModal(true)}>Add Host</Button>
                </Col>
            </Row>

            <Suspense fallback={<Spinner animation="border" role="status" />}>
                <Row xs={2} md={4} style={{textAlign: "center"}} className="mx-5 mt-3">
                    {
                        hosts.length > 0 && hosts.map((host: HostMachine, index: number) => (
                            <Col key={host.uuid}>
                                <Card>
                                    <CardBody>
                                        <CardTitle>
                                            <Row>
                                                <Col className="ms-5">{host.hostname}</Col>
                                                <Col xs={2}>
                                                    <OverlayTrigger overlay={<Tooltip>Delete Host</Tooltip>}>
                                                        <Trash style={{color: pending ? "grey" : "red", cursor: "pointer"}} onClick={() => pending ? null : deleteHost(index)} />
                                                    </OverlayTrigger>
                                                </Col>
                                            </Row>
                                        </CardTitle>
                                        <FloatingLabel label="Hostname" controlId="hostname" key="hostname" className="mt-2">
                                            <FormControl type="text" defaultValue={host.hostname} onChange={(e) => host.hostname = e.target.value} />
                                        </FloatingLabel>
                                        <FloatingLabel label="IPv4" controlId="ip4" key="ip4" className="mt-2">
                                            <FormControl type="text" defaultValue={host.ip4} onChange={(e) => host.ip4 = e.target.value ?? undefined} />
                                        </FloatingLabel>
                                        <FloatingLabel label="IPv6" controlId="ip6" key="ip6" className="mt-2">
                                            <FormControl type="text" defaultValue={host.ip6} onChange={(e) => host.ip6 = e.target.value ?? undefined} />
                                        </FloatingLabel>
                                        <FloatingLabel label="OS" controlId="os" key="os" className="mt-2">
                                            <FormSelect defaultValue={host.os} onChange={(e) => host.os = e.target.value}>
                                                {
                                                    oses.map((os: string) => (
                                                        <option value={os} key={os}>{os}</option>
                                                    ))
                                                }
                                            </FormSelect>
                                        </FloatingLabel>
                                        <FloatingLabel label="Team" controlId="team" key="team" className="mt-2">
                                            <FormSelect defaultValue={host.team.uuid} onChange={(e) => host.team.uuid = e.target.value}>
                                                {
                                                    teams.length > 0
                                                        ? teams.map((team: Team) => (
                                                            <option value={team.uuid} key={team.uuid}>{team.name}</option>
                                                        ))
                                                        : <option value="">No teams found</option>
                                                }
                                            </FormSelect>
                                        </FloatingLabel>
                                        <Button variant="primary" className="mt-2" onClick={() => updateHost(index)} disabled={pending}>Save</Button>
                                    </CardBody>
                                </Card>
                            </Col>
                        ))
                    }
                    {
                        hosts.length == 0 && apiError == undefined && (
                            <Col className="mx-auto my-3">
                                <h4 style={{fontSize: 24}}>No hosts found</h4>
                                <Button variant="primary" onClick={() => setShowCreateModal(true)} className="mt-2">Add Host</Button>
                            </Col>
                        )
                    }
                </Row>
            </Suspense>

            {/* Toasts for showing error messages */}
            <ToastContainer position="bottom-end" className="p-3">
                { errors.map((error: string, index: number) => (
                    showErrors[index] && (
                        <Toast bg="danger" onClose={() => onCloseToast(index)} key={`error-${index}`} autohide delay={5000}>
                            <ToastHeader>Error</ToastHeader>
                            <ToastBody>{error}</ToastBody>
                        </Toast>
                    ))
                )}
                { showAPIError && (
                    <Toast bg="danger" onClose={() => { setAPIError(undefined); }} key={"error"} autohide delay={5000}>
                        <ToastHeader>Error</ToastHeader>
                        <ToastBody>{apiError}</ToastBody>
                    </Toast>
                )}
                { showSuccessToast && (
                    <Toast bg="success" onClose={() => { setSuccessToastMessage(undefined); }} key={"success"} autohide delay={5000}>
                        <ToastHeader>Success</ToastHeader>
                        <ToastBody>{successToastMessage}</ToastBody>
                    </Toast>
                )}
            </ToastContainer>
        </div>
    );
}