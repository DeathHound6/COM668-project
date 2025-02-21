"use client";

import type { HostMachine, Team, APIError, User } from "../../interfaces";
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
    Row,
    Spinner
} from "react-bootstrap";
import { useEffect, useState } from "react";
import { z } from "zod";
import { GetTeams } from "../../actions/teams";
import { GetHosts, CreateHost, GetHost } from "../../actions/hosts";
import ToastContainerComponent from "../../components/toastContainer";
import { GetMe } from "../../actions/users";

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

export default function HostsPage() {
    const [hosts, setHosts] = useState([] as HostMachine[]);
    const [teams, setTeams] = useState([] as Team[]);
    const [errors, setErrors] = useState([] as string[]);
    const [user, setUser] = useState(undefined as User | undefined);
    const [pending, setPending] = useState(true);
    const [loaded, setLoaded] = useState(false);

    const [showCreateModal, setShowCreateModal] = useState(false);

    const [successToastMessages, setSuccessToastMessages] = useState([] as string[]);

    const [hostname, setHostname] = useState("");
    const [ip4, setIp4] = useState("");
    const [ip6, setIp6] = useState("");
    const [os, setOs] = useState("");
    const [teamID, setTeamID] = useState("");

    function handleError(err: APIError) {
        if ([400, 404, 500].includes(err.status))
            setErrors((prev) => [...prev, err.message]);
        setPending(false);
    }

    useEffect(() => {
        setOs(oses[0]);
        setLoaded(false);
        async function fetchData() {
            setPending(true);
            const userResponse = await GetMe().catch(handleError);
            setPending(false);
            if (!userResponse)
                return;
            setUser(userResponse);

            const teamsResponse = await GetTeams({ pageSize: 1000 }).catch(handleError);
            setPending(false);
            if (!teamsResponse)
                return;
            setTeams(teamsResponse.data);

            setPending(true);
            const hostsResponse = await GetHosts({}).catch(handleError);
            setPending(false);
            if (!hostsResponse)
                return;
            setHosts(hostsResponse.data);
            setLoaded(true);
        }
        fetchData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

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
            setErrors((prev) => [...prev, ...existingErrors]);
            setPending(false);
            return;
        }
        CreateHost({ hostname, os, ip4: ip4.length == 0 ? null : ip4, ip6: ip6.length == 0 ? null : ip6, teamID }).then(
            async(uuid) => {
                const host = await GetHost({ uuid }).catch(handleError);
                setPending(false);
                if (!host)
                    return;
                setHosts((prev) => [...prev, host]);
                setShowCreateModal(false);
                setSuccessToastMessages((prev) => [...prev, "Host created successfully"]);
            },
            handleError
        );
    }

    return (
        <div className="m-2 text-center">
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
                <Col style={{textAlign: "left"}}></Col>
                <Col style={{textAlign: "center", fontSize: 24}}><b>Host Inventory</b></Col>
                <Col style={{textAlign: "right"}}>
                    <Button variant="secondary" onClick={() => setShowCreateModal(true)}>Add Host</Button>
                </Col>
            </Row>

            {
                !loaded
                    ? (<Spinner role="status" animation="border" className="my-auto mx-auto" />)
                    : (
                        <>
                            {
                                hosts.length == 0
                                    ? (
                                        <div className="mx-auto mt-5">
                                            <h1 style={{fontSize: 40}}><b>No Hosts</b></h1>
                                            <br />
                                            <p style={{fontSize: 20}}>There are currently no hosts</p>
                                            <Button onClick={() => setShowCreateModal(true)} className="mt-4">Add Host</Button>
                                        </div>
                                    )
                                    : (
                                        <Row xs={2} md={4} style={{textAlign: "center"}} className="mx-5 mt-3">
                                            {
                                                hosts.map((host: HostMachine) => (
                                                    <Col key={host.uuid}>
                                                        <Card>
                                                            <CardBody>
                                                                <CardTitle>{host.hostname}</CardTitle>
                                                                <FloatingLabel label="Hostname" controlId="hostname" key="hostname" className="mt-2">
                                                                    <FormControl type="text" value={host.hostname} readOnly disabled />
                                                                </FloatingLabel>
                                                                <FloatingLabel label="IPv4" controlId="ip4" key="ip4" className="mt-2">
                                                                    <FormControl type="text" value={host.ip4 ?? ""} readOnly disabled />
                                                                </FloatingLabel>
                                                                <FloatingLabel label="IPv6" controlId="ip6" key="ip6" className="mt-2">
                                                                    <FormControl type="text" value={host.ip6 ?? ""} readOnly disabled />
                                                                </FloatingLabel>
                                                                <FloatingLabel label="OS" controlId="os" key="os" className="mt-2">
                                                                    <FormControl type="text" value={host.os} readOnly disabled />
                                                                </FloatingLabel>
                                                                <FloatingLabel label="Team" controlId="team" key="team" className="mt-2">
                                                                    <FormControl type="text" value={host.team.name} readOnly disabled />
                                                                </FloatingLabel>
                                                                <Button variant="primary" className="mt-2" href={`/hosts/${host.uuid}`} disabled={pending || !user?.admin}>Edit</Button>
                                                            </CardBody>
                                                        </Card>
                                                    </Col>
                                                ))
                                            }
                                        </Row>
                                    )
                            }
                        </>
                    )
            }

            {/* Toasts for showing error messages */}
            <ToastContainerComponent
                errors={errors}
                successMessages={successToastMessages}
                setErrors={setErrors}
                setSuccessToastMessages={setSuccessToastMessages}
                />
        </div>
    );
}