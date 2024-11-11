const Sentry = require("@sentry/node");
const { nodeProfilingIntegration } = require("@sentry/profiling-node");

Sentry.init({
    dsn: process.env.SENTRY_DSN,
    integrations: [nodeProfilingIntegration()],
    tracesSampleRate: 1.0,
    profilesSampleRate: 1.0
});

const express = require("express");
const app = express();

app.listen(parseInt(process.env.PORT), () => {
    console.log(`App is listening on port ${process.env.PORT}`);
});

app.get("/error", (req, res) => {
    throw new Error("example error");
});

Sentry.setupExpressErrorHandler(app);
