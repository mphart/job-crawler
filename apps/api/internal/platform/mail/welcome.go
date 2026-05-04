package mail

import (
	"html"
	"strings"
)

const productName = "Job Crawler"

func welcomeRegistrationSubject() string {
	// ASCII subject avoids rare SMTP header issues with Unicode punctuation.
	return "Registration Confirmation - " + productName
}

func welcomeRegistrationHTML(recipientName string) string {
	display := strings.TrimSpace(recipientName)
	if display == "" {
		display = "Valued Registrant"
	}
	safe := html.EscapeString(display)

	var b strings.Builder
	b.WriteString(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="utf-8"><title>Registration Confirmation</title></head>
<body style="margin:0;padding:24px;background:#f5f5f5;font-family:Georgia,'Times New Roman',serif;color:#1a1a1a;line-height:1.6;">
<table role="presentation" width="100%" cellspacing="0" cellpadding="0"><tr><td align="center">
<table role="presentation" width="600" cellspacing="0" cellpadding="0" style="max-width:600px;background:#ffffff;border:1px solid #d9d9d9;border-radius:2px;">
<tr><td style="padding:32px 40px 24px;border-bottom:1px solid #e5e5e5;">
<p style="margin:0;font-size:11px;letter-spacing:0.12em;text-transform:uppercase;color:#555;">Official Notice</p>
<h1 style="margin:12px 0 0;font-size:20px;font-weight:400;color:#111;">Account Registration Acknowledged</h1>
</td></tr>
<tr><td style="padding:32px 40px 40px;">
<p style="margin:0 0 16px;">Dear ` + safe + `,</p>
<p style="margin:0 0 16px;">We are writing to confirm that your registration with <strong>` + html.EscapeString(productName) + `</strong> has been completed successfully. Your account credentials are now active, and you may access the platform using the electronic mail address and password you provided at enrollment.</p>
<p style="margin:0 0 16px;">Please retain this message for your records. Should you require assistance with your account, you may contact our support representatives through the application interface.</p>
<p style="margin:0 0 24px;">We thank you for your patronage and remain at your service.</p>
<p style="margin:0;">Respectfully yours,</p>
<p style="margin:4px 0 0;"><strong>Office of Account Services</strong><br>
<span style="font-size:14px;color:#444;">` + html.EscapeString(productName) + `</span></p>
</td></tr>
<tr><td style="padding:16px 40px;background:#fafafa;border-top:1px solid #e5e5e5;font-size:12px;color:#666;">
This is an automated, system-generated transmission. Replies to this message are not monitored.
</td></tr>
</table>
</td></tr></table>
</body>
</html>`)
	return b.String()
}
