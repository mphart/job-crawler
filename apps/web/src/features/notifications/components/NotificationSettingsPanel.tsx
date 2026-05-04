import React from "react";
import { useEffect, useState } from "react";
import { Select } from "../../../shared/components/Select";
import { fetchNotificationSettings, updateNotificationSettings } from "../api/notifications.api";
import { NotificationFrequency } from "../model/notifications.types";
import { Button } from "../../../shared/components/Button";

function coerceFrequency(value: string): NotificationFrequency {
  if (value === "twice-daily" || value === "weekly") {
    return value;
  }
  return "daily";
}

export function NotificationSettingsPanel({ emailOptIn, frequency, token }: { emailOptIn: boolean; frequency: NotificationFrequency; token: string }) {
  const [current, setCurrent] = useState<NotificationFrequency>(coerceFrequency(frequency));
  const [status, setStatus] = useState<string>("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let active = true;
    setLoading(true);
    fetchNotificationSettings(token)
      .then((settings) => {
        if (!active) return;
        setCurrent(coerceFrequency(settings.frequency));
      })
      .catch(() => {
        if (!active) return;
        setStatus("Using local frequency setting.");
      })
      .finally(() => {
        if (!active) return;
        setLoading(false);
      });
    return () => {
      active = false;
    };
  }, [token]);

  async function onSave() {
    try {
      await updateNotificationSettings({ emailOptIn, frequency: current }, token);
      setStatus("Saved.");
    } catch {
      setStatus("Failed to save settings.");
    }
  }

  return (
    <section style={{ border: "1px solid var(--border)", borderRadius: 8, padding: "0.75rem", margin: "1rem 0" }}>
      <h4 style={{ marginTop: 0 }}>Notifications</h4>
      <p style={{ marginTop: 0 }}>Email opt-in: {emailOptIn ? "Enabled" : "Disabled"}</p>
      <div style={{ display: "flex", gap: "0.5rem", alignItems: "center", flexWrap: "wrap" }}>
        <Select value={current} onChange={(e) => setCurrent(e.target.value as NotificationFrequency)}>
          <option value="daily">Once a day</option>
          <option value="twice-daily">Twice a day</option>
          <option value="weekly">Once a week</option>
        </Select>
        <Button onClick={onSave} disabled={loading}>
          Save
        </Button>
        {status ? <small>{status}</small> : null}
      </div>
    </section>
  );
}
