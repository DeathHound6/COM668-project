"use client";

import { startTransition, useActionState, useEffect, useState } from "react";
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
    ToastHeader
} from "react-bootstrap";
import InputGroupText from "react-bootstrap/esm/InputGroupText";
import { XLg, Trash } from "react-bootstrap-icons";
import { z } from "zod";

const newFieldSchema = z.object({
    key: z.string().trim().min(1, "field key is required"),
    value: z.string().trim().min(1, "field value is required"),
    type: z.string().trim().min(1, "field type is required")
});
const newSettingSchema = z.object({
    name: z.string().trim().min(1, "setting name is required")
});
type FormState = {
    errors: {
        key?: string[] | undefined;
        value?: string[] | undefined;
        type?: string[] | undefined;
        name?: string[] | undefined;
    };
} | undefined;

export default function SettingsPage() {
    const [state, action, pending] = useActionState<FormState, FormData>(createNewField, { errors: { key: undefined, value: undefined, type: undefined, name: undefined } });
    const [settingState, settingAction, settingPending] = useActionState<FormState, FormData>(createSetting, { errors: { key: undefined, value: undefined, type: undefined, name: undefined } });

    const [settings, setSettings] = useState([] as any[]);
    const [providerType, setProviderType] = useState("log");

    const [newfieldProviderIndex, setNewFieldProviderIndex] = useState(-1);
    const [showNewFieldModal, setShowNewFieldModal] = useState(false);

    const [fieldKey, setFieldKey] = useState("");
    const [fieldValue, setFieldValue] = useState("");
    const [fieldType, setFieldType] = useState("string");

    const [showAPIError, setShowAPIError] = useState(false);
    const [apiError, setAPIError] = useState(undefined as string | undefined);

    const [showKeyError, setShowKeyError] = useState([] as boolean[]);
    const [showValueError, setShowValueError] = useState([] as boolean[]);
    const [showTypeError, setShowTypeError] = useState([] as boolean[]);
    const [showSettingNameError, setShowSettingNameError] = useState([] as boolean[]);

    const [showNewSettingModal, setShowNewSettingModal] = useState(false);
    const [settingName, setSettingName] = useState("");

    const [successToastMessage, setSuccessToastMessage] = useState(undefined as string | undefined);
    const [showSuccessToast, setShowSuccessToast] = useState(false);

    useEffect(() => {
        fetch(`/api/providers?provider_type=${providerType}`)
            .then(res => res.json())
            .then(data => {
                setSettings(data.providers);
            });
    }, [providerType]);

    useEffect(() => {
        setShowAPIError(apiError != undefined);
        setShowKeyError(state?.errors.key == undefined ? [] : new Array(state.errors.key?.length || 0).fill(true));
        setShowValueError(state?.errors.value == undefined ? [] : new Array(state.errors.value?.length || 0).fill(true));
        setShowTypeError(state?.errors.type == undefined ? [] : new Array(state.errors.type?.length || 0).fill(true));
        setShowSettingNameError(settingState?.errors.name == undefined ? [] : new Array(settingState.errors.name?.length || 0).fill(true));
        setShowSuccessToast(successToastMessage != undefined);
    }, [apiError, state?.errors.key, state?.errors.value, state?.errors.type, settingState?.errors.name, successToastMessage]);

    function updateSetting(index: number) {
        const setting = settings[index];
        fetch(`/api/providers/${setting.id}`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(setting)
        })
        .then(
            async(res) => {
                if (!res.ok)
                    return setAPIError((await res.json()).error);
                setSuccessToastMessage("Setting updated successfully");
            }
        );
    }

    function createNewField(state: FormState, form: FormData) {
        if (newfieldProviderIndex == -1) {
            setAPIError("no provider selected");
            return {errors: {key: undefined, value: undefined, type: undefined}};
        }

        const key = form.get("key") as string;
        const value = form.get("value") as string;
        const type = form.get("type") as string;
        const validatedFields = newFieldSchema.safeParse({ key, value, type });
        if (!validatedFields.success)
            return { errors: validatedFields.error.flatten().fieldErrors };

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
            type
        });
        setShowNewFieldModal(false);
    }

    function onFormSubmit(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault();
        const form = new FormData();
        form.append("key", fieldKey);
        form.append("value", fieldValue);
        form.append("type", fieldType);
        startTransition(() => action(form));
    }

    function deleteField(providerIndex: number, fieldKey: string) {
        const newSettings = [...settings];
        const setting = newSettings[providerIndex];
        setting.fields = setting.fields.filter((field: any) => field.key != fieldKey);
        newSettings[providerIndex] = setting;
        setSettings(newSettings);
    }

    function onCloseToast(index: number, type: string) {
        const errors = [] as boolean[];
        switch (type) {
            case "key":
                errors.push(...showKeyError);
                errors[index] = false;
                setShowKeyError(errors);
                break;
            case "value":
                errors.push(...showValueError);
                errors[index] = false;
                setShowValueError(errors);
                break;
            case "type":
                errors.push(...showTypeError);
                errors[index] = false;
                setShowTypeError(errors);
                break;
            case "setting":
                break;
        };
    }

    function deleteSetting(index: number) {
        const setting = settings[index];
        fetch(`/api/providers/${setting.id}`, {
            method: "DELETE"
        })
        .then(
            async(res) => {
                if (!res.ok)
                    return setAPIError((await res.json()).error);
                const newSettings = [...settings];
                newSettings.splice(index, 1);
                setSettings(newSettings);
                setSuccessToastMessage("Setting deleted successfully");
            }
        );
    }

    function onSettingFormSubmit(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault();
        const form = new FormData();
        form.append("name", settingName);
        startTransition(() => settingAction(form));
    }

    function createSetting(state: FormState, form: FormData) {
        const name = form.get("name") as string;
        const validatedSetting = newSettingSchema.safeParse({ name });
        if (!validatedSetting.success)
            return { errors: validatedSetting.error.flatten().fieldErrors };

        fetch(`/api/providers?provider_type=${providerType}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ name })
        })
        .then(
            async(res) => {
                if (!res.ok)
                    return setAPIError((await res.json()).error);
                const data = await res.json();
                console.log(data);
                data["fields"] = [];
                const newSettings = [...settings];
                newSettings.push(data);
                setSettings(newSettings);
                setShowNewSettingModal(false);
                setSuccessToastMessage("Setting created successfully");
            }
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
                <Form onSubmit={onFormSubmit}>
                    <ModalHeader closeButton={!pending}>
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
                <Form onSubmit={onSettingFormSubmit}>
                    <ModalHeader closeButton={!settingPending}>
                        <ModalTitle>Create New Setting</ModalTitle>
                    </ModalHeader>
                    <ModalBody>
                        <FloatingLabel controlId="newSettingName" label="Setting Name" className="mb-3">
                            <FormControl type="text" autoFocus value={settingName} onChange={(e) => setSettingName(e.target.value)} />
                        </FloatingLabel>
                    </ModalBody>
                    <ModalFooter>
                        <Button variant="primary" disabled={settingPending} type="submit">Create Setting</Button>
                    </ModalFooter>
                </Form>
            </Modal>

            <Row>
                <Col>
                    <ButtonGroup aria-label="Settings Provider Type">
                        <Button variant={providerType == "log" ? "primary" : "secondary"}
                            onClick={() => setProviderType("log")}
                            disabled={providerType == "log"}>Logs</Button>
                        <Button variant={providerType == "alert" ? "primary" : "secondary"}
                            onClick={() => setProviderType("alert")}
                            disabled={providerType == "alert"}>Alerts</Button>
                    </ButtonGroup>
                </Col>
                <Col style={{textAlign: "right"}}>
                    <Button variant="secondary" onClick={() => setShowNewSettingModal(true)}>Create Setting</Button>
                </Col>
            </Row>

            <Row style={{alignContent: "center", textAlign: "center"}} xs={2} md={4} className="mx-5">
                {
                    settings.map((setting: any, index: number) => (
                        <Col key={`col-${setting.id}`}>
                            <Card className="m-2 p-2 border rounded" key={`c-${setting.id}`}>
                                <CardBody key={`cb-${setting.id}`}>
                                    <CardTitle key={`ct-${setting.id}`}>
                                        <Row xs={12}>
                                            <Col xs={10}>{setting.name}</Col>
                                            <Col xs={2}>
                                                <OverlayTrigger overlay={<Tooltip>Delete Setting</Tooltip>}>
                                                    <Trash style={{color: "red", cursor: "pointer"}} onClick={() => deleteSetting(index)} />
                                                </OverlayTrigger>
                                            </Col>
                                        </Row>
                                    </CardTitle>
                                    {
                                        setting.fields.map((field: any) => (
                                            <InputGroup key={`ig-${setting.id}-${field.key}`} className="m-2">
                                                <FloatingLabel controlId="floatingKey" label={field.key} key={`fl-${setting.id}-${field.key}`}>
                                                    <FormControl type="text" defaultValue={field.value} key={`fc-${setting.id}-${field.key}`} />
                                                </FloatingLabel>
                                                <InputGroupText key={`igt-${setting.id}-${field.key}`}>{field.type}</InputGroupText>
                                                <OverlayTrigger overlay={<Tooltip>Delete Field</Tooltip>}>
                                                    <InputGroupText style={{cursor: "pointer", color: "red"}} onClick={() => deleteField(index, field.key)}>
                                                        <XLg />
                                                    </InputGroupText>
                                                </OverlayTrigger>
                                            </InputGroup>
                                        ))
                                    }
                                    <Button variant="secondary" onClick={() => {setNewFieldProviderIndex(index); setShowNewFieldModal(true);}}>Create new field</Button>
                                    <br />
                                    <Button variant="primary" className="mt-2" onClick={() => updateSetting(index)}>Save</Button>
                                </CardBody>
                            </Card>
                        </Col>
                    ))
                }
            </Row>

            {/* Toasts for showing error messages */}
            <ToastContainer position="bottom-end" className="p-3">
                { state?.errors.key?.map((error: string, index: number) => (
                    showKeyError[index] && (
                        <Toast bg="danger" onClose={() => onCloseToast(index, "key")} key={`k-${index}`}>
                            <ToastHeader>Error</ToastHeader>
                            <ToastBody>{error}</ToastBody>
                        </Toast>
                    ))
                )}
                { state?.errors.value?.map((error: string, index: number) => (
                    showValueError[index] && (
                        <Toast bg="danger" onClose={() => onCloseToast(index, "value")} key={`v-${index}`}>
                            <ToastHeader>Error</ToastHeader>
                            <ToastBody>{error}</ToastBody>
                        </Toast>
                    ))
                )}
                { state?.errors.type?.map((error: string, index: number) => (
                    showTypeError[index] && (
                        <Toast bg="danger" onClose={() => onCloseToast(index, "type")} key={`t-${index}`}>
                            <ToastHeader>Error</ToastHeader>
                            <ToastBody>{error}</ToastBody>
                        </Toast>
                    ))
                )}
                { showAPIError && (
                    <Toast bg="danger" onClose={() => { setAPIError(undefined); }} key={"error"}>
                        <ToastHeader>Error</ToastHeader>
                        <ToastBody>{apiError}</ToastBody>
                    </Toast>
                )}
                { settingState?.errors.name?.map((error: string, index: number) => (
                    showSettingNameError[index] && (
                        <Toast bg="danger" onClose={() => onCloseToast(index, "setting")} key={`s-${index}`}>
                            <ToastHeader>Error</ToastHeader>
                            <ToastBody>{error}</ToastBody>
                        </Toast>
                    ))
                )}
                { showSuccessToast && (
                    <Toast bg="success" onClose={() => { setSuccessToastMessage(undefined); }} key={"success"}>
                        <ToastHeader>Success</ToastHeader>
                        <ToastBody>{successToastMessage}</ToastBody>
                    </Toast>
                )}
            </ToastContainer>
        </div>
    );
}
