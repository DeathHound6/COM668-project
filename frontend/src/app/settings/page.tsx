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
    OverlayTrigger,
    Tooltip,
    Modal,
    ModalHeader,
    ModalTitle,
    ModalBody,
    ModalFooter,
    FormSelect,
    Spinner,
    FormCheck
} from "react-bootstrap";
import InputGroupText from "react-bootstrap/esm/InputGroupText";
import { XLg, Trash } from "react-bootstrap-icons";
import { z } from "zod";
import { CreateSetting, DeleteSetting, GetSetting, GetSettings, UpdateSetting } from "../../actions/settings";
import ToastContainerComponent from "../../components/toastContainer";

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
    const [pending, setPending] = useState(true);
    const [loaded, setLoaded] = useState(false);

    const [settings, setSettings] = useState([] as Settings[]);
    const [providerType, setProviderType] = useState("log" as "alert"|"log");

    const [newfieldProviderIndex, setNewFieldProviderIndex] = useState(-1);
    const [showNewFieldModal, setShowNewFieldModal] = useState(false);

    const [fieldKey, setFieldKey] = useState("");
    const [fieldValue, setFieldValue] = useState("");
    const [fieldType, setFieldType] = useState("string");
    const [fieldRequired, setFieldRequired] = useState(false);

    const [errors, setErrors] = useState([] as string[]);
    const [showErrors, setShowErrors] = useState([] as boolean[]);

    const [showNewSettingModal, setShowNewSettingModal] = useState(false);
    const [settingName, setSettingName] = useState("");

    const [successToastMessage, setSuccessToastMessage] = useState(undefined as string | undefined);
    const [showSuccessToast, setShowSuccessToast] = useState(false);

    function handleError(error: APIError) {
        if ([400, 404, 500].includes(error.status))
            setErrors((prev) => [...prev, error.message]);
        setLoaded(true);
        setPending(false);
    }

    useEffect(() => {
        setPending(true);
        async function fetchData() {
            setLoaded(false);
            GetSettings({ providerType })
                .then(
                    (data) => {
                        setSettings(data.data);
                        setLoaded(true);
                        setPending(false);
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
        setShowErrors(new Array(errors.length).fill(true));
        setShowSuccessToast(successToastMessage != undefined);
    }, [errors, successToastMessage]);

    function updateSetting(index: number) {
        setPending(true);
        if (settings.length <= index) {
            setErrors((prev) => [...prev, "invalid provider index"]);
            setPending(false);
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
            setPending(false);
            setErrors((prev) => [...prev, "no provider selected"]);
            return;
        }
        if (settings.length <= newfieldProviderIndex) {
            setPending(false);
            setErrors((prev) => [...prev, "invalid provider index"]);
            return;
        }

        const key = fieldKey;
        const value = fieldValue;
        const type = fieldType;
        const required = fieldRequired;

        const validatedFields = newFieldSchema.safeParse({ key, value, type, required });
        if (!validatedFields.success) {
            const newErrors = validatedFields.error.flatten().fieldErrors ?? { key: [], value: [], type: [], required: [] };
            const existingErrors = [
                ...newErrors.key ?? [],
                ...newErrors.value ?? [],
                ...newErrors.type ?? [],
                ...newErrors.required ?? []
            ];
            setErrors((prev) => [...prev, ...existingErrors]);
            setPending(false);
            return;
        }

        if (type == "number") {
            const parsedValue = parseFloat(value);
            if (isNaN(parsedValue)) {
                setErrors((prev) => [...prev, "invalid number value"]);
                setPending(false);
                return;
            }
        } else if (type === "bool") {
            const valid = ["true", "false", "1", "0"];
            if (!valid.includes(value.toLowerCase())) {
                setErrors((prev) => [...prev, "invalid boolean value"]);
                setPending(false);
                return;
            }
        } else if (type == "string") {}

        const setting = settings[newfieldProviderIndex];
        setting.fields.push({
            key,
            value,
            type,
            required
        });
        setSettings((prev) => { const temp = [...prev]; temp[newfieldProviderIndex] = setting; return temp; });
        setShowNewFieldModal(false);
        setPending(false);
    }

    function deleteField(providerIndex: number, fieldKey: string) {
        setPending(true);
        if (settings.length <= providerIndex) {
            setPending(false);
            setErrors((prev) => [...prev, "invalid provider index"]);
            return;
        }

        const newSettings = [...settings];
        const setting = newSettings[providerIndex];
        setting.fields = setting.fields.filter((field: SettingField) => field.key != fieldKey);
        newSettings[providerIndex] = setting;
        setSettings(newSettings);
        setPending(false);
    }

    function deleteSetting(index: number) {
        setPending(true);
        const setting = settings[index];
        DeleteSetting({ uuid: setting.uuid })
            .then(
                () => {
                    const newSettings = [...settings];
                    newSettings.splice(index, 1);
                    setPending(false);
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
            setErrors((prev) => [...prev, ...(newErrors.name ?? [])]);
            setPending(false);
            return;
        }
        CreateSetting({ name, providerType })
            .then(
                async(uuid) => {
                    const setting = await GetSetting({ uuid }).catch(handleError);
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
        <div className="m-2" style={{textAlign: "center"}}>
            {/* Modal for creating a new setting field */}
            <Modal
               show={showNewFieldModal}
               onHide={() => { setNewFieldProviderIndex(-1); setShowNewFieldModal(false); }}
               backdrop="static"
               centered={true}
               restoreFocus={false}>
                <ModalHeader closeButton>
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
                    <Button variant="primary" disabled={pending} onClick={() => createField()}>Create Field</Button>
                </ModalFooter>
            </Modal>

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
                        </div>
                    )
            }

            {/* Toasts for showing error messages */}
            <ToastContainerComponent
                errors={errors}
                showErrors={showErrors}
                successMessage={successToastMessage}
                showSuccessMessage={showSuccessToast}
                setErrors={setErrors}
                setSuccessToastMessage={setSuccessToastMessage}
                />
        </div>
    );
}
