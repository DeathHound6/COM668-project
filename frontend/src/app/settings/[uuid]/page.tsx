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
    value: z.string().trim().min(1, "field value is required"),
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
    const [successMessages, setSuccessMessages] = useState([] as string[]);

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
        async function fetchData() {
            setPending(true);
            const userResponse = await GetMe().catch(handleError);
            if (!userResponse)
                return;
            if (!userResponse.admin)
                redirect("/dashboard", RedirectType.replace);

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
// eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    function validateField(data: SettingField): boolean {
        let newErrors: { key?: string[], value?: string[], type?: string[], required?: string[]} = { key: [], value: [], type: [], required: [] };
        const validatedFields = newFieldSchema.safeParse(data);
        if (!validatedFields.success)
            newErrors = validatedFields.error.flatten().fieldErrors ?? newErrors;
        if (fieldRequired && ["", null, undefined].includes(fieldValue))
            newErrors.value?.push("required field cannot have an empty value");
        else
            switch (fieldType) {
                case "number":
                    const parsedValue = parseFloat(fieldValue);
                    if (isNaN(parsedValue))
                        newErrors.value?.push("invalid number value");
                    break;
                case "bool":
                    const valid = ["true", "false", "1", "0"];
                    if (!valid.includes(fieldValue.toLowerCase()))
                        newErrors.value?.push("invalid boolean value");
                    break;
                case "string":
                default:
                    break;
            }
        const existingErrors = [
            ...newErrors.key ?? [],
            ...newErrors.value ?? [],
            ...newErrors.type ?? [],
            ...newErrors.required ?? []
        ];
        if (existingErrors.length > 0) {
            setErrors((prev) => [...prev, ...existingErrors]);
            setPending(false);
            return false;
        }
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
            if (!updateResponse)
                return;
            setSuccessMessages((prev) => [...prev, "Field created successfully"]);
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
            if (!updateResponse)
                return;
            setFields(newFields);
            setSuccessMessages((prev) => [...prev, "Field deleted successfully"]);
            setPending(false);
        }
        deleteF();
    }

    function updateSetting() {
        setPending(true);
        async function update() {
            const validatedName = settingSchema.safeParse({ name: settingName });
            if (!validatedName.success) {
                const newErrors = validatedName.error.flatten().fieldErrors ?? { name: [] };
                setErrors((prev) => [...prev, ...(newErrors.name ?? [])]);
                setPending(false);
                return;
            }
            for (const field of fields)
                if (!validateField(field))
                    return;
            const updateResponse = await UpdateSetting({ ...setting, name: settingName, fields }).catch(handleError);
            if (!updateResponse)
                return;
            setSetting({ ...setting, name: settingName, fields });
            setSuccessMessages((prev) => [...prev, "Setting updated successfully"]);
            setPending(false);
        }
        update();
    }

    function deleteSetting() {
        setPending(true);
        async function deleteS() {
            const deleteResponse = await DeleteSetting({ uuid: setting.uuid }).catch(handleError);
            if (!deleteResponse)
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
                    <FormCheck className="mx-auto" label="Required Field" checked={fieldRequired} onChange={(e) => setFieldRequired(e.target.checked)} type="switch" />
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
                            {
                                fields.length == 0
                                    ? (<h1 style={{fontSize: 20}} className="pt-3">No Fields Found</h1>)
                                    : (<Row xs={3} className="mx-80">
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
    )
}