import smtplib
from email.mime.text import MIMEText

def send_email():
    smtp_server = "smtp.qq.com"
    smtp_port = 465  # QQ邮箱SSL端口
    sender_email = "316969606@qq.com"
    sender_password = "791212yy"
    receiver_email = "316969606@qq.com"
    subject = "Hello"
    body = "hello world!"

    try:
        msg = MIMEText(body, 'plain', 'utf-8')
        msg['From'] = sender_email
        msg['To'] = receiver_email
        msg['Subject'] = subject

        # 使用SSL连接
        server = smtplib.SMTP_SSL(smtp_server, smtp_port)
        server.login(sender_email, sender_password)
        server.sendmail(sender_email, receiver_email, msg.as_string())
        server.quit()
        print("Email sent successfully!")
    except Exception as e:
        print(f"Error sending email: {e}")

if __name__ == "__main__":
    send_email()
