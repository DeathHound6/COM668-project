import { Routes } from '@angular/router';
import { AppComponent } from "./app.component.ts";

export const routes: Routes = [
    {
        title: "Dashboard",
        path: "/(dashboard)?",
        component: AppComponent
    },
    {
        title: "Alert Settings",
        path: "/settings/alerts"
    },
    {
        title: "Settings",
        path: "/settings/plugins"
    }
];