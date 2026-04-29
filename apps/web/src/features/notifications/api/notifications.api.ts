import { mockRequestJson } from "../../../shared/api/client";
import { NotificationSettings } from "../model/notifications.types";

// MOCK: Pending backend notification preference endpoints.
export async function updateNotificationSettings(settings: NotificationSettings): Promise<NotificationSettings> {
  return mockRequestJson(() => settings);
}
