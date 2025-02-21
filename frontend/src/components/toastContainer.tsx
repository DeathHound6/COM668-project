import { ToastContainer, Toast, ToastHeader, ToastBody } from "react-bootstrap";

export default function ToastContainerComponent(
    { errors, successMessages, setErrors, setSuccessToastMessages }:
    { errors: string[], successMessages: string[], setErrors: (errors: string[]) => void, setSuccessToastMessages: (messages: string[]) => void }
) {
    return (
        <ToastContainer position="bottom-end" className="p-3">
            { errors.map((error: string, index: number) => (
                <Toast bg="danger" onClose={() => {
                    const e = [...errors];
                    e.splice(index, 1);
                    setErrors(e);
                    }} key={`error-${index}`} autohide delay={5000}>
                    <ToastHeader>Error</ToastHeader>
                    <ToastBody>{error}</ToastBody>
                </Toast>
                )
            )}
            { successMessages.map((message: string, index: number) => (
                <Toast bg="success" onClose={() => {
                    const e = [...successMessages];
                    e.splice(index, 1);
                    setSuccessToastMessages(e);
                    }} key={`success-${index}`} autohide delay={5000}>
                    <ToastHeader>Success</ToastHeader>
                    <ToastBody>{message}</ToastBody>
                </Toast>
                )
            )}
        </ToastContainer>
    );
}