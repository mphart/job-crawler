import { requestJson } from "../../../shared/api/client";
import { NotificationSettings } from "../model/notifications.types";

export async function fetchNotificationSettings(token: string): Promise<NotificationSettings> {
  return requestJson<NotificationSettings>("/api/notifications/settings", "GET", undefined, { token });
}

export async function updateNotificationSettings(settings: NotificationSettings, token: string): Promise<NotificationSettings> {
  return requestJson<NotificationSettings>("/api/notifications/settings", "PATCH", settings, { token });
}
