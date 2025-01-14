export default function IncidentCard({ incident }: Readonly<{ incident: any }>) {
    return (
        <div>
            <h1>{incident.title}</h1>
            <p>{incident.description}</p>
        </div>
    );
}