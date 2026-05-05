import type { NotificationFrequency } from "../../notifications/model/notifications.types";

export type LoginRequest = { email: string; password: string };

export type SignupPreferencesPayload = {
  keywords: string[];
  locations: string[];
  desiredTitles: string[];
  preferredCompanies: string[];
  minComp: number;
  emailOptIn: boolean;
  darkMode: boolean;
};

export type SignupRequest = {
  email: string;
  name: string;
  password: string;
  notificationFrequency: NotificationFrequency;
  preferences: SignupPreferencesPayload;
  resumeFileName?: string;
  resumeContentBase64?: string;
};

export type SessionUser = { id: string; email: string; username: string; token: string };
