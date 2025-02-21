"use client";

import type { APIError, SettingField, Settings } from "../../interfaces";
import { useEffect, useState } from "react";
import {
    Button,
    ButtonGroup,
    Card,
    CardBody,
    CardTitle,
    Row,
    Col,
    FloatingLabel,
    FormControl,
    InputGroup,
    Modal,
    ModalHeader,
    ModalTitle,
    ModalBody,
    ModalFooter,
    Spinner
} from "react-bootstrap";
import InputGroupText from "react-bootstrap/esm/InputGroupText";
import { z } from "zod";
import { CreateSetting, GetSetting, GetSettings } from "../../actions/settings";
import ToastContainerComponent from "../../components/toastContainer";
import { GetMe } from "../../actions/users";
import { RedirectType, redirect } from "next/navigation";

const newSettingSchema = z.object({
    name: z.string().trim().min(1, "setting name is required")
});

export default function SettingsPage() {
    const [pending, setPending] = useState(true);
    const [loaded, setLoaded] = useState(false);

    const [settings, setSettings] = useState([] as Settings[]);
    const [providerType, setProviderType] = useState("log" as "alert"|"log");

    const [errors, setErrors] = useState([] as string[]);

    const [showNewSettingModal, setShowNewSettingModal] = useState(false);
    const [settingName, setSettingName] = useState("");

    const [successToastMessages, setSuccessToastMessages] = useState([] as string[]);

    function handleError(error: APIError) {
        if ([400, 404, 500].includes(error.status))
            setErrors((prev) => [...prev, error.message]);
        setPending(false);
    }

    useEffect(() => {
        setPending(true);
        async function fetchData() {
            setLoaded(false);
            const userResponse = await GetMe().catch(handleError);
            if (!userResponse || !userResponse.admin)
                redirect("/dashboard", RedirectType.replace);
            const settingsResponse = await GetSettings({ providerType }).catch(handleError);
            if (!settingsResponse)
                return;
            setSettings(settingsResponse.data);
            setLoaded(true);
            setPending(false);
        }
        fetchData();
    }, [providerType]);

    function createSetting() {
        setPending(true);
        const name = settingName;
        const validatedSetting = newSettingSchema.safeParse({ name });
        if (!validatedSetting.success) {
            const newErrors = validatedSetting.error.flatten().fieldErrors ?? { name: [] };
            setErrors((prev) => [...prev, ...(newErrors.name ?? [])]);
            setPending(false);
            return;
        }
        async function post() {
            const postResponse = await CreateSetting({ name, providerType }).catch(handleError);
            if (!postResponse)
                return;
            const setting = await GetSetting({ uuid: postResponse }).catch(handleError);
            if (!setting)
                return;
            const newSettings = [...settings];
            newSettings.push(setting);
            setSettings(newSettings);
            setShowNewSettingModal(false);
            setSuccessToastMessages((prev) => [...prev, "Setting created successfully"]);
            setPending(false);
        }
        post();
    }

    return (
        <div className="m-2" style={{textAlign: "center"}}>
            {/* Modal for creating a new setting */}
            <Modal
               show={showNewSettingModal}
               onHide={() => setShowNewSettingModal(false)}
               backdrop="static"
               centered={true}
               restoreFocus={false}>
                <ModalHeader closeButton>
                    <ModalTitle>Create New Setting</ModalTitle>
                </ModalHeader>
                <ModalBody>
                    <FloatingLabel controlId="newSettingName" label="Setting Name" className="mb-3">
                        <FormControl type="text" autoFocus value={settingName} onChange={(e) => setSettingName(e.target.value)} />
                    </FloatingLabel>
                </ModalBody>
                <ModalFooter>
                    <Button variant="primary" disabled={pending} onClick={() => createSetting()}>Create Setting</Button>
                </ModalFooter>
            </Modal>

            <Row className="mt-3">
                <Col style={{textAlign: "left"}}>
                    <ButtonGroup aria-label="Settings Provider Type">
                        <Button variant={providerType == "log" ? "primary" : "secondary"}
                            onClick={() => setProviderType("log")}
                            disabled={providerType == "log"}>Logs</Button>
                        <Button variant={providerType == "alert" ? "primary" : "secondary"}
                            onClick={() => setProviderType("alert")}
                            disabled={providerType == "alert"}>Alerts</Button>
                    </ButtonGroup>
                </Col>
                <Col style={{textAlign: "center"}}>
                    <h1 style={{fontSize: 24}}><b>Settings</b></h1>
                </Col>
                <Col style={{textAlign: "right"}}>
                    <Button variant="secondary" onClick={() => setShowNewSettingModal(true)}>Create Setting</Button>
                </Col>
            </Row>

            {
                !loaded
                    ? (<Spinner role="status" animation="border" className="my-auto mx-auto" />)
                    : (
                        <div>
                        {
                            settings.length == 0
                                ? (
                                    <div className="mx-auto mt-5">
                                        <h1 style={{fontSize: 40}}><b>No Settings</b></h1>
                                        <br />
                                        <p style={{fontSize: 20}}>No settings were found of the current type</p>
                                        <Button variant="secondary" onClick={() => setProviderType((prev) => prev == "alert" ? "log" : "alert")} className="mt-4">Switch Setting Type</Button>
                                        <Button onClick={() => setShowNewSettingModal(true)} className="mt-4">Create Setting</Button>
                                    </div>
                                )
                                : (
                                    <Row style={{textAlign: "center"}} xs={2} md={4} className="mx-5 mt-3">
                                        {
                                            /* Render settings */
                                            settings.length > 0 && settings.map((setting: Settings, index: number) => (
                                                <Col key={`col-${setting.uuid}`}>
                                                    <Card className="m-2 p-2 border rounded" key={`c-${setting.uuid}`}>
                                                        <CardBody key={`cb-${setting.uuid}`}>
                                                            <CardTitle key={`ct-${setting.uuid}`}>{setting.name}</CardTitle>
                                                            {
                                                                setting.fields.map((field: SettingField) => (
                                                                    <InputGroup key={`ig-${setting.uuid}-${field.key}`} className="m-2">
                                                                        <FloatingLabel controlId="floatingKey" label={field.key}>
                                                                            <FormControl type="text" value={field.value} disabled />
                                                                        </FloatingLabel>
                                                                        <InputGroupText>{field.type}</InputGroupText>
                                                                        <InputGroup.Checkbox checked={field.required} disabled />
                                                                    </InputGroup>
                                                                ))
                                                            }
                                                            <Button variant="primary" className="mt-2" href={`/settings/${setting.uuid}`}>Edit</Button>
                                                        </CardBody>
                                                    </Card>
                                                </Col>
                                            ))
                                        }
                                        </Row>
                                )
                        }
                        </div>
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
