import React from "react";
import { useState } from "react";
import { Button } from "../../../shared/components/Button";

export function ResumeUploader({
  current,
  onUpload,
}: {
  current?: string;
  onUpload: (resume: { fileName: string; contentBase64: string }) => Promise<void>;
}) {
  const [value, setValue] = useState(current ?? "");
  const [contentBase64, setContentBase64] = useState("");
  const [status, setStatus] = useState("");

  return (
    <section style={{ display: "flex", gap: "0.5rem", alignItems: "center", flexWrap: "wrap" }}>
      <input
        type="file"
        accept="application/pdf,.pdf"
        onChange={async (event) => {
          const file = event.target.files?.[0];
          if (!file) return;
          if (file.type !== "application/pdf" && !file.name.toLowerCase().endsWith(".pdf")) return;
          setValue(file.name);
          const base64 = await fileToBase64(file);
          setContentBase64(base64);
        }}
      />
      <Button
        onClick={async () => {
          await onUpload({ fileName: value || "resume.pdf", contentBase64 });
          setStatus("Uploaded.");
          window.alert("Resume uploaded.");
        }}
      >
        Upload Resume
      </Button>
      {value ? <small>Selected: {value}</small> : null}
      {status ? <small>{status}</small> : null}
    </section>
  );
}

function fileToBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => {
      const dataUrl = String(reader.result ?? "");
      const base64 = dataUrl.includes(",") ? dataUrl.split(",")[1] : "";
      resolve(base64);
    };
    reader.onerror = () => reject(reader.error);
    reader.readAsDataURL(file);
  });
}
