import { useState } from "react";
import { Button } from "../../../shared/components/Button";
import { Select } from "../../../shared/components/Select";
import { updateNotificationSettings } from "../api/notifications.api";
import { NotificationFrequency } from "../model/notifications.types";
export function NotificationSettingsPanel({emailOptIn,frequency}:{emailOptIn:boolean;frequency:NotificationFrequency;}){const [current,setCurrent]=useState<NotificationFrequency>(frequency);return <section style={{border:"1px solid var(--border)",borderRadius:8,padding:"0.75rem",margin:"1rem 0"}}><h4 style={{marginTop:0}}>Notifications</h4><p style={{marginTop:0}}>Email opt-in: {emailOptIn?"Enabled":"Disabled"}</p><div style={{display:"flex",gap:"0.5rem",alignItems:"center"}}><Select value={current} onChange={(e)=>setCurrent(e.target.value as NotificationFrequency)}><option value="daily">Daily</option><option value="twice-daily">Twice daily</option><option value="instant">Instant</option></Select><Button onClick={()=>void updateNotificationSettings({emailOptIn,frequency:current})}>Save</Button></div></section>;}
