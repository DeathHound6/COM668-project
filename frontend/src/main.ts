import { bootstrapApplication } from "npm:@angular/platform-browser";
import { appConfig } from "./app/app.config.ts";
import { AppComponent } from "./app/app.component.ts";

bootstrapApplication(AppComponent, appConfig).catch((_: Error) => {});
