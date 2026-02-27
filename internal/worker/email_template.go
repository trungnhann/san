package worker

import "fmt"

func generateOTPEmailContent(otp string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Verify Your Account</title>
    <style>
        body {
            font-family: 'Courier New', Courier, monospace;
            background-color: #0f172a;
            color: #e2e8f0;
            padding: 20px;
            margin: 0;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #1e293b;
            border: 1px solid #334155;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
        }
        .header {
            background-color: #3b82f6;
            padding: 20px;
            text-align: center;
        }
        .header h1 {
            margin: 0;
            color: #ffffff;
            font-size: 24px;
            letter-spacing: 2px;
        }
        .content {
            padding: 40px 20px;
            text-align: center;
        }
        .otp-box {
            background-color: #0f172a;
            border: 2px dashed #3b82f6;
            color: #3b82f6;
            font-size: 32px;
            font-weight: bold;
            letter-spacing: 5px;
            padding: 20px;
            margin: 30px auto;
            display: inline-block;
            border-radius: 4px;
        }
        .footer {
            background-color: #0f172a;
            padding: 20px;
            text-align: center;
            font-size: 12px;
            color: #94a3b8;
            border-top: 1px solid #334155;
        }
        p {
            line-height: 1.6;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>SAN API SYSTEM</h1>
        </div>
        <div class="content">
            <p>Hello, Developer!</p>
            <p>We received a request to verify your email address. Use the secure code below to complete your registration:</p>
            
            <div class="otp-box">%s</div>
            
            <p>This code will expire in 15 minutes.</p>
            <p>If you didn't request this, please ignore this message.</p>
        </div>
        <div class="footer">
            &copy; 2026 San API. All rights reserved.<br>
            Secure Verification System
        </div>
    </div>
</body>
</html>
`, otp)
}
