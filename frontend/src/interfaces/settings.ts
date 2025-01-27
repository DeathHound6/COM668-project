export interface Settings {
    uuid: string;
    name: string;
    type: string;
    fields: SettingField[];
}

export interface SettingField {
    key: string;
    type: string;
    value: string;
    required: boolean;
}
