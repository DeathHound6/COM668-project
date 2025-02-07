"use client";

import type { APIError, SettingField, Settings } from "../../interfaces";
import { startTransition, Suspense, useActionState, useEffect, useState } from "react";
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
    OverlayTrigger,
    Tooltip,
    Modal,
    ModalHeader,
    ModalTitle,
    ModalBody,
    ModalFooter,
    FormSelect,
    Form,
    ToastContainer,
    Toast,
    ToastBody,
    ToastHeader,
    Spinner,
    FormCheck
} from "react-bootstrap";
import InputGroupText from "react-bootstrap/esm/InputGroupText";
import { XLg, Trash } from "react-bootstrap-icons";
import { set, z } from "zod";
import { CreateSetting, DeleteSetting, GetSetting, GetSettings, UpdateSetting } from "../../actions/settings";

const newFieldSchema = z.object({
    key: z.string().trim().min(1, "field key is required"),
    value: z.string().trim().min(1, "field value is required"),
    type: z.string().trim().min(1, "field type is required"),
    required: z.boolean()
});
const newSettingSchema = z.object({
    name: z.string().trim().min(1, "setting name is required")
});

export default function SettingsPage() {
    const [pending, setPending] = useState(false);
    const [loaded, setLoaded] = useState(false);

    const [settings, setSettings] = useState([] as Settings[]);
    const [providerType, setProviderType] = useState("log" as "alert"|"log");

    const [newfieldProviderIndex, setNewFieldProviderIndex] = useState(-1);
    const [showNewFieldModal, setShowNewFieldModal] = useState(false);

    const [fieldKey, setFieldKey] = useState("");
    const [fieldValue, setFieldValue] = useState("");
    const [fieldType, setFieldType] = useState("string");
    const [fieldRequired, setFieldRequired] = useState(false);

    const [showAPIError, setShowAPIError] = useState(false);
    const [apiError, setAPIError] = useState(undefined as string | undefined);
    const [errors, setErrors] = useState([] as string[]);
    const [showErrors, setShowErrors] = useState([] as boolean[]);

    const [showNewSettingModal, setShowNewSettingModal] = useState(false);
    const [settingName, setSettingName] = useState("");

    const [successToastMessage, setSuccessToastMessage] = useState(undefined as string | undefined);
    const [showSuccessToast, setShowSuccessToast] = useState(false);

    function handleError(error: APIError) {
        if ([400, 404, 500].includes(error.status))
            setAPIError(error.message);
        setLoaded(true);
        setPending(false);
    }

    useEffect(() => {
        async function fetchData() {
            setLoaded(false);
            GetSettings({ providerType })
                .then(
                    (data) => {
                        setSettings(data.data);
                        setLoaded(true);
                    },
                    (err) => {
                        handleError(err);
                        setSettings([]);
                    }
                );
        }
        fetchData();
    }, [providerType]);

    useEffect(() => {
        setShowAPIError(apiError != undefined);
        setShowErrors
        setShowSuccessToast(successToastMessage != undefined);
    }, [apiError, errors, successToastMessage]);

    function updateSetting(index: number) {
        setPending(true);
        if (settings.length <= index) {
            setAPIError("invalid provider index");
            return;
        }

        const setting = settings[index];
        UpdateSetting(setting)
            .then(
                () => {
                    setSuccessToastMessage("Setting updated successfully");
                    setPending(false);
                },
                handleError
            );
    }

    function createField() {
        setPending(true);

        if (newfieldProviderIndex == -1) {
            setAPIError("no provider selected");
            return;
        }
        if (settings.length <= newfieldProviderIndex) {
            setAPIError("invalid provider index");
            return;
        }

        const key = fieldKey;
        const value = fieldValue;
        const type = fieldType;
        const required = fieldRequired;

        const validatedFields = newFieldSchema.safeParse({ key, value, type });
        if (!validatedFields.success) {
            const newErrors = validatedFields.error.flatten().fieldErrors ?? { key: [], value: [], type: [] };
            const existingErrors = [
                ...newErrors.key ?? [],
                ...newErrors.value ?? [],
                ...newErrors.type ?? []
            ];
            setErrors(existingErrors);
            setPending(false);
            return;
        }

        if (type == "number") {
            const parsedValue = parseFloat(value);
            if (isNaN(parsedValue)) {
                return { errors: { key: undefined, value: ["invalid number"], type: undefined } };
            }
        } else if (type === "bool") {
            const valid = ["true", "false", "1", "0"];
            if (!valid.includes(value.toLowerCase()))
                return { errors: { key: undefined, value: ["invalid boolean"], type: undefined } };
        } else if (type == "string") {}

        const setting = settings[newfieldProviderIndex];
        setting.fields.push({
            key,
            value,
            type,
            required
        });
        setShowNewFieldModal(false);
        setPending(false);
    }

    function deleteField(providerIndex: number, fieldKey: string) {
        setPending(true);
        if (settings.length <= providerIndex) {
            setAPIError("invalid provider index");
            return;
        }

        const newSettings = [...settings];
        const setting = newSettings[providerIndex];
        setting.fields = setting.fields.filter((field: SettingField) => field.key != fieldKey);
        newSettings[providerIndex] = setting;
        setSettings(newSettings);
        setPending(false);
    }

    function onCloseToast(index: number) {
        const e = [...showErrors];
        if (e.length <= index) {
            setAPIError("invalid error index");
            return;
        }
        e[index] = false;
        setShowErrors(e);
    }

    function deleteSetting(index: number) {
        const setting = settings[index];
        DeleteSetting({ uuid: setting.uuid })
            .then(
                () => {
                    const newSettings = [...settings];
                    newSettings.splice(index, 1);
                    setSettings(newSettings);
                    setSuccessToastMessage("Setting deleted successfully");
                },
                handleError
            );
    }

    function createSetting() {
        setPending(true);
        const name = settingName;
        const validatedSetting = newSettingSchema.safeParse({ name });
        if (!validatedSetting.success) {
            const newErrors = validatedSetting.error.flatten().fieldErrors ?? { name: [] };
            setErrors(newErrors.name ?? []);
            return;
        }
        CreateSetting({ name, providerType })
            .then(
                async(data) => {
                    const setting = await GetSetting({ uuid: data }).catch(handleError);
                    if (!setting)
                        return;
                    const newSettings = [...settings];
                    newSettings.push(setting);
                    setSettings(newSettings);
                    setShowNewSettingModal(false);
                    setSuccessToastMessage("Setting created successfully");
                    setPending(false);
                },
                handleError
            );
    }

    return (
        <div className="m-2">
            {/* Modal for creating a new setting field */}
            <Modal
               show={showNewFieldModal}
               onHide={() => { setNewFieldProviderIndex(-1); setShowNewFieldModal(false); }}
               backdrop="static"
               centered={true}
               restoreFocus={false}>
                <Form onSubmit={() => createField()}>
                    <ModalHeader>
                        <ModalTitle>Create New Field</ModalTitle>
                    </ModalHeader>
                    <ModalBody>
                        <FloatingLabel controlId="newFieldKey" label="Field Key" className="mb-3">
                            <FormControl type="text" autoFocus value={fieldKey} onChange={(e) => setFieldKey(e.target.value)} />
                        </FloatingLabel>
                        <FloatingLabel controlId="newFieldValue" label="Field Value" className="mb-3">
                            <FormControl type="text" value={fieldValue} onChange={(e) => setFieldValue(e.target.value)} />
                        </FloatingLabel>
                        <FloatingLabel controlId="newFieldDescription" label="Field Data Type" className="mb-3">
                            <FormSelect value={fieldType} onChange={(e) => setFieldType(e.target.value)}>
                                <option value="string">String</option>
                                <option value="number">Number</option>
                                <option value="bool">Boolean</option>
                            </FormSelect>
                        </FloatingLabel>
                        <FormCheck className="mx-auto" label="Required Field" checked={fieldRequired} onChange={(e) => setFieldRequired(e.target.checked)} />
                    </ModalBody>
                    <ModalFooter>
                        <Button variant="primary" disabled={pending} type="submit">Create Field</Button>
                    </ModalFooter>
                </Form>
            </Modal>

            {/* Modal for creating a new setting */}
            <Modal
               show={showNewSettingModal}
               onHide={() => setShowNewSettingModal(false)}
               backdrop="static"
               centered={true}
               restoreFocus={false}>
                <Form onSubmit={() => createSetting()}>
                    <ModalHeader>
                        <ModalTitle>Create New Setting</ModalTitle>
                    </ModalHeader>
                    <ModalBody>
                        <FloatingLabel controlId="newSettingName" label="Setting Name" className="mb-3">
                            <FormControl type="text" autoFocus value={settingName} onChange={(e) => setSettingName(e.target.value)} />
                        </FloatingLabel>
                    </ModalBody>
                    <ModalFooter>
                        <Button variant="primary" disabled={pending} type="submit">Create Setting</Button>
                    </ModalFooter>
                </Form>
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
                        <>
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
                                                            <CardTitle key={`ct-${setting.uuid}`}>
                                                                <Row>
                                                                    <Col className="ms-5">{setting.name}</Col>
                                                                    <Col xs={2}>
                                                                        <OverlayTrigger overlay={<Tooltip>Delete Setting</Tooltip>}>
                                                                            <Trash style={{color: pending ? "gray" : "red", cursor: "pointer"}} onClick={() => pending ? null : deleteSetting(index)} />
                                                                        </OverlayTrigger>
                                                                    </Col>
                                                                </Row>
                                                            </CardTitle>
                                                            {
                                                                setting.fields.map((field: SettingField) => (
                                                                    <InputGroup key={`ig-${setting.uuid}-${field.key}`} className="m-2">
                                                                        <FloatingLabel controlId="floatingKey" label={field.key} key={`fl-${setting.uuid}-${field.key}`}>
                                                                            <FormControl type="text" defaultValue={field.value} key={`fc-${setting.uuid}-${field.key}`} />
                                                                        </FloatingLabel>
                                                                        <InputGroupText key={`igt-${setting.uuid}-${field.key}`}>{field.type}</InputGroupText>
                                                                        <OverlayTrigger overlay={<Tooltip>{field.required ? "Field is required" : "Delete Field"}</Tooltip>}>
                                                                            <InputGroupText
                                                                            style={{cursor: "pointer", color: field.required ? "grey" : "red"}}
                                                                            onClick={() => field.required ? null : deleteField(index, field.key)}>
                                                                                <XLg />
                                                                            </InputGroupText>
                                                                        </OverlayTrigger>
                                                                    </InputGroup>
                                                                ))
                                                            }
                                                            <Button variant="secondary" onClick={() => {setNewFieldProviderIndex(index); setShowNewFieldModal(true);}}>Create new field</Button>
                                                            <br />
                                                            <Button variant="primary" className="mt-2" onClick={() => updateSetting(index)} disabled={pending}>Save</Button>
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
