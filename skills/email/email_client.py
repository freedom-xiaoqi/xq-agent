import argparse
import smtplib
import imaplib
import email
import os
import sys
from email.header import decode_header
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart

def send_email(smtp_server, smtp_port, sender_email, sender_password, receiver_email, subject, body):
    try:
        msg = MIMEMultipart()
        msg['From'] = sender_email
        msg['To'] = receiver_email
        msg['Subject'] = subject

        msg.attach(MIMEText(body, 'plain'))

        server = smtplib.SMTP(smtp_server, int(smtp_port))
        server.starttls()
        server.login(sender_email, sender_password)
        text = msg.as_string()
        server.sendmail(sender_email, receiver_email, text)
        server.quit()
        print("Email sent successfully.")
    except Exception as e:
        print(f"Error sending email: {e}")
        sys.exit(1)

def receive_email(imap_server, imap_port, email_user, email_pass, limit=5):
    try:
        mail = imaplib.IMAP4_SSL(imap_server, int(imap_port))
        mail.login(email_user, email_pass)
        mail.select("inbox")

        status, messages = mail.search(None, "ALL")
        if status != "OK":
             print("No messages found!")
             return

        email_ids = messages[0].split()
        
        # Get latest emails
        latest_email_ids = email_ids[-limit:]
        
        results = []

        for e_id in reversed(latest_email_ids):
            res, msg_data = mail.fetch(e_id, "(RFC822)")
            for response_part in msg_data:
                if isinstance(response_part, tuple):
                    msg = email.message_from_bytes(response_part[1])
                    
                    subject, encoding = decode_header(msg["Subject"])[0]
                    if isinstance(subject, bytes):
                        subject = subject.decode(encoding if encoding else "utf-8")
                    
                    from_ = msg.get("From")
                    
                    body = ""
                    if msg.is_multipart():
                        for part in msg.walk():
                            content_type = part.get_content_type()
                            content_disposition = str(part.get("Content-Disposition"))
                            try:
                                body = part.get_payload(decode=True).decode()
                            except:
                                pass
                            if content_type == "text/plain" and "attachment" not in content_disposition:
                                break
                    else:
                        body = msg.get_payload(decode=True).decode()

                    print(f"From: {from_}")
                    print(f"Subject: {subject}")
                    print(f"Body: {body[:200]}...")
                    print("---")

        mail.close()
        mail.logout()
            
    except Exception as e:
        print(f"Error receiving email: {e}")
        sys.exit(1)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Email Skill Script")
    subparsers = parser.add_subparsers(dest="command", help="Command to run")

    # Send Command
    send_parser = subparsers.add_parser("send", help="Send an email")
    send_parser.add_argument("--to", required=True, help="Receiver email")
    send_parser.add_argument("--subject", required=True, help="Email subject")
    send_parser.add_argument("--body", required=True, help="Email body")

    # Receive Command
    receive_parser = subparsers.add_parser("receive", help="Receive emails")
    receive_parser.add_argument("--limit", type=int, default=5, help="Number of emails to fetch")

    args = parser.parse_args()

    # Configuration from Environment Variables
    # Users should set these in SKILL.md env section
    SMTP_SERVER = os.getenv("EMAIL_SMTP_SERVER", "smtp.gmail.com")
    SMTP_PORT = os.getenv("EMAIL_SMTP_PORT", "587")
    IMAP_SERVER = os.getenv("EMAIL_IMAP_SERVER", "imap.gmail.com")
    IMAP_PORT = os.getenv("EMAIL_IMAP_PORT", "993")
    EMAIL_USER = os.getenv("EMAIL_USER")
    EMAIL_PASS = os.getenv("EMAIL_PASS")

    if not EMAIL_USER or not EMAIL_PASS:
        print("Error: EMAIL_USER and EMAIL_PASS environment variables must be set.")
        sys.exit(1)

    if args.command == "send":
        send_email(SMTP_SERVER, SMTP_PORT, EMAIL_USER, EMAIL_PASS, args.to, args.subject, args.body)
    elif args.command == "receive":
        receive_email(IMAP_SERVER, IMAP_PORT, EMAIL_USER, EMAIL_PASS, args.limit)
    else:
        parser.print_help()
