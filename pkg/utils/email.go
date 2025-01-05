package utils

import (
	"bytes"
	"cli-project/internal/config"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"
	"time"
)

//// sendEmail sends an email using the provided parameters
//func sendEmail(to, subject, body string) error {
//	from := "codesageofficialindia@gmail.com"
//	appPassword := config.APP_PASSWORD
//
//	// SMTP server configuration
//	smtpHost := "smtp.gmail.com"
//	smtpPort := "587"
//
//	// Construct the email headers, including the Content-Type for HTML
//	headers := map[string]string{
//		"From":         from,
//		"To":           to,
//		"Subject":      subject,
//		"MIME-Version": "1.0",
//		"Content-Type": "text/html; charset=\"UTF-8\"",
//	}
//
//	// Combine the headers and body into a single message
//	var message strings.Builder
//	for key, value := range headers {
//		message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
//	}
//	message.WriteString("\r\n" + body)
//
//	// Authentication
//	auth := smtp.PlainAuth("", from, appPassword, smtpHost)
//
//	// Sending email
//	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message.String()))
//	if err != nil {
//		return fmt.Errorf("failed to send email: %w", err)
//	}
//
//	return nil
//}
//
//// SendOTPEmail sends an OTP email with the provided OTP
//func SendOTPEmail(to, otp string) error {
//	subject := "Your OTP Code For Codesage"
//
//	body := fmt.Sprintf(`
//	<html>
//		<body>
//			<h2 style="color: #333;">Hello!</h2>
//			<p style="font-size: 16px; color: #555;">
//				You requested a one-time password (OTP) for verification. Please use the code below to proceed:
//			</p>
//			<p style="font-size: 20px; font-weight: bold; color: #000;">
//				OTP Code: %s
//			</p>
//			<p style="font-size: 16px; color: #555;">
//				This code will be valid for the next 15 minutes. If you did not request this, please disregard this email or contact support if you have any concerns.
//			</p>
//			<br>
//			<p style="font-size: 14px; color: #888;">
//				Best regards,<br>
//				The Codesage Team
//			</p>
//		</body>
//	</html>
//	`, otp)
//
//	return sendEmail(to, subject, body)
//}

// EmailConfig holds the configuration for sending emails
type EmailConfig struct {
	FromEmail    string
	AppPassword  string
	SMTPHost     string
	SMTPPort     string
	CompanyName  string
	SupportEmail string
}

// EmailTemplate represents an email template with styling
type EmailTemplate struct {
	Subject string
	HTML    string
}

// NewEmailConfig creates a new email configuration
func NewEmailConfig() *EmailConfig {
	return &EmailConfig{
		FromEmail:    "codesageofficialindia@gmail.com",
		AppPassword:  config.APP_PASSWORD,
		SMTPHost:     "smtp.gmail.com",
		SMTPPort:     "587",
		CompanyName:  "Codesage",
		SupportEmail: "pravincodes@outlook.com",
	}
}

// sendEmail sends an email using the provided parameters with simplified error handling
func (cfg *EmailConfig) sendEmail(to, subject, body string) error {
	// Add RFC 822 date header
	currentTime := time.Now().Format(time.RFC822)

	headers := map[string]string{
		"From":         fmt.Sprintf("%s <%s>", cfg.CompanyName, cfg.FromEmail),
		"To":           to,
		"Subject":      subject,
		"Date":         currentTime,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=\"UTF-8\"",
		"X-Priority":   "1", // High priority
	}

	var message strings.Builder
	for key, value := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	message.WriteString("\r\n" + body)

	auth := smtp.PlainAuth("", cfg.FromEmail, cfg.AppPassword, cfg.SMTPHost)

	// Simple email sending without timeout
	err := smtp.SendMail(
		cfg.SMTPHost+":"+cfg.SMTPPort,
		auth,
		cfg.FromEmail,
		[]string{to},
		[]byte(message.String()),
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// getBaseTemplate returns the base HTML template with consistent styling
func getBaseTemplate() string {
	return `
    <!DOCTYPE html>
    <html>
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <style>
                body {
                    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
                    line-height: 1.6;
                    color: #333;
                    max-width: 600px;
                    margin: 0 auto;
                    padding: 20px;
                }
                .header {
                    background: #2C3E50;
                    color: white;
                    padding: 20px;
                    border-radius: 8px 8px 0 0;
                    text-align: center;
                }
                .content {
                    background: #ffffff;
                    padding: 20px;
                    border-radius: 0 0 8px 8px;
                    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
                }
                .otp-code {
                    font-size: 32px;
                    font-weight: bold;
                    color: #2C3E50;
                    text-align: center;
                    padding: 20px;
                    background: #f8f9fa;
                    border-radius: 4px;
                    margin: 20px 0;
                    letter-spacing: 5px;
                }
                .footer {
                    text-align: center;
                    margin-top: 20px;
                    color: #666;
                    font-size: 14px;
                }
                .button {
                    display: inline-block;
                    padding: 10px 20px;
                    background: #2C3E50;
                    color: white;
                    text-decoration: none;
                    border-radius: 4px;
                    margin: 10px 0;
                }
                @media (max-width: 600px) {
                    body {
                        padding: 10px;
                    }
                    .otp-code {
                        font-size: 24px;
                    }
                }
            </style>
        </head>
        <body>
            {{.Content}}
            <div class="footer">
                <p>© {{.Year}} {{.CompanyName}}. All rights reserved.</p>
                <p>If you didn't request this email, please contact <a href="mailto:{{.SupportEmail}}">support</a>.</p>
            </div>
        </body>
    </html>`
}

// SendOTPEmail sends an enhanced OTP email with better styling and security notices
func (cfg *EmailConfig) SendOTPEmail(to, otp string) error {
	subject := fmt.Sprintf("Your Verification Code - %s", cfg.CompanyName)

	// Create the email content
	emailContent := fmt.Sprintf(`
        <div class="header">
            <h1>%s</h1>
        </div>
        <div class="content">
            <h2>Verification Required</h2>
            <p>Use the code below to complete your verification:</p>
            <div class="otp-code">%s</div>
            <p><strong>Important Security Notes:</strong></p>
            <ul>
                <li>This code expires in 15 minutes</li>
                <li>Never share this code with anyone</li>
            </ul>
            <p>If you didn't request this code, please ignore this email or contact our support team immediately.</p>
        </div>`,
		cfg.CompanyName,
		otp,
	)

	// Create template data
	templateData := struct {
		Content      template.HTML
		Year         int
		CompanyName  string
		SupportEmail string
	}{
		Content:      template.HTML(emailContent),
		Year:         time.Now().Year(),
		CompanyName:  cfg.CompanyName,
		SupportEmail: cfg.SupportEmail,
	}

	// Parse the base template
	tmpl, err := template.New("email").Parse(getBaseTemplate())
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute the template
	var body bytes.Buffer
	if err := tmpl.Execute(&body, templateData); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return cfg.sendEmail(to, subject, body.String())
}
