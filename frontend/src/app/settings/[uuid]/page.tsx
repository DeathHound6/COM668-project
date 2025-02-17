"use client";

import type { APIError, SettingField, Settings } from "../../../interfaces";
import { DeleteSetting, GetSetting, UpdateSetting } from "../../../actions/settings";
import { useState, useEffect } from "react";
import ToastContainerComponent from "../../../components/toastContainer";
import { Modal, ModalHeader, ModalTitle, ModalBody, FloatingLabel, FormControl, FormSelect, FormCheck, ModalFooter, Button, Spinner, Row, Col, ButtonGroup, Card, CardBody, OverlayTrigger, Tooltip } from "react-bootstrap";
import { z } from "zod";
import { redirect, RedirectType } from "next/navigation";
import { GetMe } from "../../../actions/users";
import { Trash } from "react-bootstrap-icons";

const newFieldSchema = z.object({
    key: z.string().trim().min(1, "field key is required"),
    value: z.string().trim().optional(),
    type: z.string().trim().min(1, "field type is required"),
    required: z.boolean()
});
const settingSchema = z.object({
    name: z.string().trim().min(1, "setting name is required")
});

export default function SettingPage({ params }: { params: Promise<{uuid: string}> }) {
    const [loaded, setLoaded] = useState(false);
    const [pending, setPending] = useState(true);
    const [setting, setSetting] = useState({} as Settings);

    const [errors, setErrors] = useState([] as string[]);
    const [showErrors, setShowErrors] = useState([] as boolean[]);
    const [successMessage, setSuccessMessage] = useState(undefined as string | undefined);
    const [showSuccessMessage, setShowSuccessMessage] = useState(false);

    const [showNewFieldModal, setShowNewFieldModal] = useState(false);

    const [settingName, setSettingName] = useState("");
    const [fieldKey, setFieldKey] = useState("");
    const [fieldValue, setFieldValue] = useState("");
    const [fieldType, setFieldType] = useState("string" as "string"|"number"|"bool");
    const [fieldRequired, setFieldRequired] = useState(false);
    const [fields, setFields] = useState([] as SettingField[]);

    function handleError(error: APIError) {
        if ([400, 403, 500].includes(error.status))
            setErrors((prev) => [...prev, error.message]);
        setPending(false);
    }

    useEffect(() => {
        setShowErrors(new Array(errors.length).fill(true));
        setShowSuccessMessage(successMessage != undefined);
    }, [errors, successMessage]);

    useEffect(() => {
        async function fetchData() {
            setPending(true);
            const userResponse = await GetMe().catch(handleError);
            if (!userResponse)
                return;
            if (!userResponse.admin)
                redirect("/settings", RedirectType.replace);

            const settingResponse = await GetSetting({ uuid: (await params).uuid }).catch(handleError);
            setLoaded(true);
            setPending(false);
            if (!settingResponse)
                return;
            setSetting(settingResponse);
            setSettingName(settingResponse.name);
            setFields(settingResponse.fields);
        }
        fetchData();
    }, []);

    function validateField(data: SettingField): boolean {
        const validatedFields = newFieldSchema.safeParse(data);
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
            return false;
        }

        if (fieldRequired && ["", null, undefined].includes(fieldValue)) {
            setErrors((prev) => [...prev, "required field cannot have an empty value"]);
            setPending(false);
            return false;
        }

        if (fieldType == "number") {
            const parsedValue = parseFloat(fieldValue);
            if (isNaN(parsedValue)) {
                setErrors((prev) => [...prev, "invalid number value"]);
                setPending(false);
                return false;
            }
        } else if (fieldType == "bool") {
            const valid = ["true", "false", "1", "0"];
            if (!valid.includes(fieldValue.toLowerCase())) {
                setErrors((prev) => [...prev, "invalid boolean value"]);
                setPending(false);
                return false;
            }
        } else if (fieldType == "string") {}
        return true;
    }

    function createField() {
        setPending(true);

        if (fields.findIndex((f) => f.key == fieldKey) != -1) {
            setErrors((prev) => [...prev, "field key already exists"]);
            setPending(false);
            return;
        }

        const data = { key: fieldKey, value: fieldValue, type: fieldType, required: fieldRequired };
        if (!validateField(data))
            return;

        async function update() {
            const newFields: SettingField[] = [...fields, data];
            const updateResponse = await UpdateSetting({ ...setting, fields: newFields }).catch(handleError);
            if (updateResponse != undefined)
                return;
            setSuccessMessage("Field created successfully");
            setFieldKey("");
            setFieldValue("");
            setFieldType("string");
            setFieldRequired(false);
            setFields(newFields);
            setShowNewFieldModal(false);
            setPending(false);
        }
        update();
    }

    function deleteField(index: number) {
        setPending(true);
        async function deleteF() {
            if (index == -1) {
                setErrors((prev) => [...prev, "field not found"]);
                return;
            }
            const newFields = [...fields];
            newFields.splice(index, 1);
            const updateResponse = await UpdateSetting({ ...setting, name: settingName, fields: newFields }).catch(handleError);
            if (updateResponse != undefined)
                return;
            setFields(newFields);
            setSuccessMessage("Field deleted successfully");
            setPending(false);
        }
        deleteF();
    }

    function updateSetting() {
        setPending(true);
        const validatedName = settingSchema.safeParse({ name: setting.name });
        if (!validatedName.success) {
            const newErrors = validatedName.error.flatten().fieldErrors ?? { name: [] };
            setErrors((prev) => [...prev, ...(newErrors.name ?? [])]);
            return;
        }
        for (const field of fields)
            if (!validateField(field))
                return;
        async function update() {
            const updateResponse = await UpdateSetting({ ...setting, name: settingName, fields }).catch(handleError);
            if (updateResponse != undefined)
                return;
            setSetting({ ...setting, name: settingName, fields });
            setSuccessMessage("Setting updated successfully");
            setPending(false);
        }
        update();
    }

    function deleteSetting() {
        setPending(true);
        async function deleteS() {
            const deleteResponse = await DeleteSetting({ uuid: setting.uuid }).catch(handleError);
            if (deleteResponse != undefined)
                return;
            redirect("/settings", RedirectType.replace);
        }
        deleteS();
    }

    return (
        <main className="m-2 text-center">
            {/* Modal for creating a new setting field */}
            <Modal
               show={showNewFieldModal}
               onHide={() => setShowNewFieldModal(false)}
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
                        <FormSelect value={fieldType} onChange={(e) => setFieldType(e.target.value as "string"|"number"|"bool")}>
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

            {
                !loaded
                    ? (<Spinner role="status" animation="border" className="my-auto mx-auto" />)
                    : (
                        <div className="mt-3">
                            <Row>
                                <Col style={{textAlign: "left"}}></Col>
                                <Col style={{textAlign: "center"}}></Col>
                                <Col style={{textAlign: "right"}}>
                                    <ButtonGroup>
                                        <Button variant="secondary" onClick={() => setShowNewFieldModal(true)}>Create Field</Button>
                                        <Button variant="primary" onClick={() => updateSetting()} disabled={pending}>Update Setting</Button>
                                        <Button variant="danger" onClick={() => deleteSetting()} disabled={pending}>Delete Setting</Button>
                                    </ButtonGroup>
                                </Col>
                            </Row>

                            <div className="mx-auto max-w-96">
                                <h1 style={{fontSize: 24}}><b>{setting.name}</b></h1>
                                <FloatingLabel controlId="settingName" label="Setting Name" className="my-3">
                                    <FormControl type="text" value={settingName} onChange={(e) => setSettingName(e.target.value)} />
                                </FloatingLabel>
                            </div>
                            <h1 style={{fontSize: 20, textDecoration: "underline"}}><b>Fields</b></h1>
                            <Row xs={3} className="mx-80">
                                {
                                    fields.map((field: SettingField) => (
                                        <Col className="max-w-80" key={field.key}>
                                            <Card className="mt-3">
                                                <CardBody>
                                                    <OverlayTrigger overlay={<Tooltip>Delete Field</Tooltip>}>
                                                        <Trash onClick={() => pending ? null : deleteField(fields.findIndex((f) => f.key == field.key))} className="ms-auto me-2 mb-2 cursor-pointer" style={{color: pending ? "grey" : "red"}} />
                                                    </OverlayTrigger>
                                                    <FloatingLabel controlId="fieldKey" label="Field Key" className="mb-3">
                                                        <FormControl type="text" value={field.key} readOnly disabled />
                                                    </FloatingLabel>
                                                    <FloatingLabel controlId="fieldValue" label="Field Value" className="mb-3">
                                                        <FormControl type="text" value={field.value} onChange={(e) => setFields((prev) => prev.map((f) => f.key == field.key ? { ...f, value: e.target.value } : f))} />
                                                    </FloatingLabel>
                                                    <FloatingLabel controlId="fieldType" label="Field Data Type" className="mb-3">
                                                        <FormSelect value={field.type} onChange={(e) => setFields((prev) => prev.map((f) => f.key == field.key ? { ...f, type: e.target.value } : f))}>
                                                            <option value="string">String</option>
                                                            <option value="number">Number</option>
                                                            <option value="bool">Boolean</option>
                                                        </FormSelect>
                                                    </FloatingLabel>
                                                    <div className="mx-auto max-w-40">
                                                        <FormCheck className="" type="switch" label="Required Field" checked={field.required} onChange={(e) => setFields((prev) => prev.map((f) => f.key == field.key ? { ...f, required: e.target.checked} : f))} />
                                                    </div>
                                                </CardBody>
                                            </Card>
                                        </Col>
                                    ))
                                }
                            </Row>
                        </div>
                    )
            }

            {/* Toasts for showing error messages */}
            <ToastContainerComponent
                successMessage={successMessage}
                showSuccessMessage={showSuccessMessage}
                errors={errors}
                showErrors={showErrors}
                setErrors={setErrors}
                setSuccessToastMessage={setSuccessMessage}
                />
        </main>
    )
}