---
name: email-client
description: 邮件客户端技能，支持发送和接收邮件。
openclaw:
  env:
    - EMAIL_SMTP_SERVER=smtp.qq.com
    - EMAIL_SMTP_PORT=587
    - EMAIL_IMAP_SERVER=imap.qq.com
    - EMAIL_IMAP_PORT=993
    - EMAIL_USER=
    - EMAIL_PASS=
  bins:
    - python
---

To use this skill, you can run the python script `email_client.py`.

## Send Email
Command:
```bash
python {{path}}/email_client.py send --to "receiver@example.com" --subject "Hello" --body "This is a test email."
```

## Receive Email
Command:
```bash
python {{path}}/email_client.py receive --limit 5
```

**Note**: Please update the `env` section in `SKILL.md` with your actual email credentials before using.
