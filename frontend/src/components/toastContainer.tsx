import { ToastContainer, Toast, ToastHeader, ToastBody } from "react-bootstrap";

export default function ToastContainerComponent(
    { errors, showErrors, successMessage, showSuccessMessage, setErrors, setSuccessToastMessage }:
    { errors: string[], showErrors: boolean[], successMessage: string | undefined, showSuccessMessage: boolean, setErrors: (errors: string[]) => void, setSuccessToastMessage: (message: string|undefined) => void }
) {
    return (
        <ToastContainer position="bottom-end" className="p-3">
            { errors.map((error: string, index: number) => (
                showErrors[index] && (
                    <Toast bg="danger" onClose={() => {
                       const e = [...errors];
                       e.splice(index, 1);
                       setErrors(e);
                       }} key={`error-${index}`} autohide delay={5000}>
                        <ToastHeader>Error</ToastHeader>
                        <ToastBody>{error}</ToastBody>
                    </Toast>
                ))
            )}
            { showSuccessMessage && (
                <Toast bg="success" onClose={() => setSuccessToastMessage(undefined)} key={"success"} autohide delay={5000}>
                    <ToastHeader>Success</ToastHeader>
                    <ToastBody>{successMessage}</ToastBody>
                </Toast>
            )}
        </ToastContainer>
    );
}