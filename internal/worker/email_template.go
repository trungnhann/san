package worker

import "fmt"

func generateOTPEmailContent(otp string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verify Your Account</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background-color: #f1f5f9;
            color: #334155;
            margin: 0;
            padding: 0;
            line-height: 1.6;
            -webkit-font-smoothing: antialiased;
        }
        .wrapper {
            width: 100%%;
            table-layout: fixed;
            background-color: #f1f5f9;
            padding-bottom: 40px;
        }
        .main {
            background-color: #ffffff;
            margin: 0 auto;
            width: 100%%;
            max-width: 520px;
            border-spacing: 0;
            font-family: sans-serif;
            color: #171a1b;
            border-radius: 12px;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);
            border: 1px solid #e2e8f0;
            overflow: hidden;
        }
        .header {
            padding: 40px 0 24px;
            text-align: center;
            border-bottom: 1px solid #f1f5f9;
        }
        .logo {
            font-size: 28px;
            font-weight: 800;
            color: #2563eb;
            text-decoration: none;
            letter-spacing: -1px;
            display: inline-block;
        }
        .logo span {
            color: #0f172a;
        }
        .body {
            padding: 40px 40px 20px;
            text-align: center;
        }
        h1 {
            font-size: 22px;
            font-weight: 700;
            color: #0f172a;
            margin: 0 0 16px;
            letter-spacing: -0.5px;
        }
        p {
            color: #64748b;
            margin: 0 0 24px;
            font-size: 15px;
            line-height: 1.6;
        }
        .otp-container {
            margin: 32px 0;
            text-align: center;
        }
        .otp-code {
            font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, Courier, monospace;
            background-color: #eff6ff;
            color: #2563eb;
            font-size: 36px;
            font-weight: 700;
            letter-spacing: 8px;
            padding: 20px 32px;
            border-radius: 12px;
            display: inline-block;
            border: 1px dashed #bfdbfe;
        }
        .warning {
            font-size: 13px;
            color: #94a3b8;
            margin-top: 32px;
            padding-top: 24px;
            border-top: 1px solid #f1f5f9;
        }
        .footer {
            padding: 32px;
            text-align: center;
            font-size: 12px;
            color: #94a3b8;
        }
        .footer a {
            color: #64748b;
            text-decoration: none;
            font-weight: 500;
        }
        .footer a:hover {
            text-decoration: underline;
        }
        /* Mobile adjustments */
        @media screen and (max-width: 520px) {
            .main {
                border-radius: 0;
                border: none;
                box-shadow: none;
            }
            .body {
                padding: 32px 24px;
            }
            .otp-code {
                font-size: 28px;
                padding: 16px 24px;
                letter-spacing: 6px;
            }
        }
    </style>
</head>
<body>
    <div class="wrapper">
        <table align="center" class="main">
            <tr>
                <td class="header">
                    <div class="logo">San<span>.Blog</span></div>
                </td>
            </tr>
            <tr>
                <td class="body">
                    <h1>Verify your email address</h1>
                    <p>Welcome to San! Use the following one-time password (OTP) to complete your sign up procedures. This code is valid for 15 minutes.</p>
                    
                    <div class="otp-container">
                        <div class="otp-code">%s</div>
                    </div>
                    
                    <p class="warning">
                        If you didn't request this code, you can safely ignore this email. Someone might have typed your email address by mistake.
                    </p>
                </td>
            </tr>
        </table>
        
        <div class="footer">
            &copy; 2026 San Project. All rights reserved.<br>
            <a href="#">Privacy Policy</a> &bull; <a href="#">Terms of Service</a>
        </div>
    </div>
</body>
</html>
`, otp)
}
