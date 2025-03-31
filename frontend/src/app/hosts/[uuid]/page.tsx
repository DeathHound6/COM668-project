"use client";

import { Button, Card, CardBody, CardTitle, Col, FloatingLabel, FormControl, FormSelect, Row, Spinner } from "react-bootstrap";
import { DeleteHost, GetHost, UpdateHost } from "../../../actions/hosts";
import { APIError, type Team, type HostMachine, type User } from "../../../interfaces";
import { useEffect, useState } from "react";
import { GetTeams } from "../../../actions/teams";
import { z } from "zod";
import ToastContainerComponent from "../../../components/toastContainer";
import { GetMe } from "../../../actions/users";
import { redirect, RedirectType } from "next/navigation";

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

export default function HostDetailsPage({ params }: { params: Promise<{ uuid: string }> }) {
    const [user, setUser] = useState(undefined as User | undefined);
    const [host, setHost] = useState(undefined as HostMachine | undefined);
    const [teams, setTeams] = useState([] as Team[]);

    const [pending, setPending] = useState(true);
    const [loaded, setLoaded] = useState(false);

    const [successMessages, setSuccessMessages] = useState([] as string[]);
    const [errors, setErrors] = useState([] as string[]);

    const [hostname, setHostname] = useState("");
    const [ip4, setIp4] = useState("");
    const [ip6, setIp6] = useState("");
    const [os, setOs] = useState("");
    const [teamID, setTeamID] = useState("");

    function handleError(err: APIError) {
        if ([400, 500].includes(err.status))
            setErrors((prev) => [...prev, err.message]);
        setLoaded(true);
        setPending(false);
    }

    useEffect(() => {
        setLoaded(false);
        async function fetchData() {
            setPending(true);
            const userResponse = await GetMe().catch(handleError);
            if (userResponse == undefined)
                return;
            setUser(userResponse);
            const teamsResponse = await GetTeams({ pageSize: 1000 }).catch(handleError);
            if (teamsResponse == undefined)
                return;
            setTeams(teamsResponse.data);
            const hostResponse = await GetHost({ uuid: (await params).uuid }).catch(handleError);
            if (hostResponse == undefined)
                return;
            setHost(hostResponse);
            setHostname(hostResponse.hostname);
            setIp4(hostResponse.ip4 ?? "");
            setIp6(hostResponse.ip6 ?? "");
            setOs(hostResponse.os);
            setTeamID(hostResponse.team.uuid);
            setLoaded(true);
            setPending(false);
        };
        fetchData();
// eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    async function deleteHost() {
        // this shouldnt hit, this is just to keep typescript happy
        if (host == undefined)
            return;
        setPending(true);
        DeleteHost(host.uuid).then(
            async() => {
                setPending(false);
                redirect("/hosts", RedirectType.replace);
            },
            handleError
        );
    }

    function updateHost() {
        // this shouldnt hit, this is just to keep typescript happy
        if (host == undefined)
            return;
        setPending(true);
        const hn = hostname;
        const ipv4 = ip4.length > 0 ? ip4 : "";
        const ipv6 = ip6.length > 0 ? ip6 : "";
        const OS = os;
        const teamUUID = teamID;
        const body = { hostname: hn, ip4: ipv4.length > 0 ? ipv4 : undefined, ip6: ipv6.length > 0 ? ipv6 : undefined, os: OS, teamID: teamUUID };
        const validatedFields = hostSchema.safeParse(body);
        if (!validatedFields.success) {
            const newErrors = validatedFields.error.flatten().fieldErrors;
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
        UpdateHost({ uuid: host.uuid, body }).then(
            () => {
                setPending(false);
                setSuccessMessages((prev) => [...prev, "Host updated successfully"]);
            },
            handleError
        );
    }

    return (
        <main className="text-center mx-2">
            {
                !loaded
                    ? (<Spinner animation="border" role="status" className="mt-5" />)
                    : (
                        <div>
                            {
                                host == undefined
                                    ? (
                                        <div></div>
                                    )
                                    : (
                                        <div>
                                            <Row className="mt-3">
                                                <Col style={{textAlign: "left"}}></Col>
                                                <Col style={{textAlign: "center"}}></Col>
                                                <Col style={{textAlign: "right"}}>
                                                    <Button variant="danger" onClick={() => deleteHost()} disabled={pending || !user?.admin}>Delete Host</Button>
                                                </Col>
                                            </Row>

                                            <Card className="mt-4 mx-auto max-w-96">
                                                <CardBody>
                                                    <CardTitle>{hostname}</CardTitle>
                                                    <FloatingLabel label="Hostname" controlId="hostname" key="hostname" className="mt-2">
                                                        <FormControl type="text" value={hostname} onChange={(e) => setHostname(e.target.value)} />
                                                    </FloatingLabel>
                                                    <FloatingLabel label="IPv4" controlId="ip4" key="ip4" className="mt-2">
                                                        <FormControl type="text" value={ip4 ?? ""} onChange={(e) => setIp4(e.target.value)} />
                                                    </FloatingLabel>
                                                    <FloatingLabel label="IPv6" controlId="ip6" key="ip6" className="mt-2">
                                                        <FormControl type="text" value={ip6 ?? ""} onChange={(e) => setIp6(e.target.value)} />
                                                    </FloatingLabel>
                                                    <FloatingLabel label="OS" controlId="os" key="os" className="mt-2">
                                                        <FormSelect value={os} onChange={(e) => setOs(e.target.value)}>
                                                            {
                                                                oses.map((os: string) => (
                                                                    <option value={os} key={os}>{os}</option>
                                                                ))
                                                            }
                                                        </FormSelect>
                                                    </FloatingLabel>
                                                    <FloatingLabel label="Team" controlId="team" key="team" className="mt-2">
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
                                                    <Button variant="primary" className="mt-2" onClick={() => updateHost()} disabled={pending || !user?.admin}>Save</Button>
                                                </CardBody>
                                            </Card>
                                        </div>
                                    )
                            }
                        </div>
                    )
            }

            {/* Toasts for showing error messages */}
            <ToastContainerComponent
                successMessages={successMessages}
                errors={errors}
                setErrors={setErrors}
                setSuccessToastMessages={setSuccessMessages}
                />
        </main>
    );
}
