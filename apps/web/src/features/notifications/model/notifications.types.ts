export type NotificationFrequency = "daily" | "twice-daily" | "weekly";
export type NotificationSettings = { emailOptIn: boolean; frequency: NotificationFrequency };
