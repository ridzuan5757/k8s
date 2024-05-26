Agenda:
- Shell will introduce us on the operation process and how it relates with proactive monitoring
    - SIP
    - Capillary
    - GHL
    - Others?

- Finalize the acceptance criteria on POC alerting mechanism:
    - Webhooks API:
        - Slack
        - Microsoft Team
        - Others (Discord etc.)
    - Mail:
        - MailGun (Quota limited to 50k per month)
        - Amazon SNS ($2.00 per 100,000 notifications + $0.09 per GB on first 10 TB)

- Finalize unclear requirements on POS instrumentation:
    - Proposed metrics for POS & CDS (final):
        - Battery level (%) `GAUGE | device.battery.level | float | 0.00 - 1.00` 
        - Average RAM Memory usage (%) `GAUGE | device.memory.utilization.avg | float | 0.00 - 1.00`
        - Available disk memory (byte) `GAUGE | device.disk.free | float | 0.00 - 1.00`
        - Average CPU usage `GAUGE | device.cpu.utilization.avg | float | 0.00 - 1.00`
    - Resource information to be added on each metrics:
        - `site.id | int`
        - `device.id | int`

- Clarification on HUB instrumentation
    - Proposed metrics for HUB (POC only):
        - Connection status to GHL `GAUGE | connection.status.ghl | int`
        - Connection status to SIP `GAUGE | connection.status.sip | int`
        - Connection status to Capillary `GAUGE | connection.status.capillary | int`
    - Resource information to be added on each metrics:
        - `site.id | int`
        - `device.id | int`

- Design idea on OpenSearch Dashboard for users (Silentmode Supports & Strateq)

Participants:
- Muhamad Ridzuan bin Mohd Yazid Goi
  ridzuan.yazid@silentmode.my
- Wan Maryam Qistina Binti Muhammad Firdaus
  maryam.firdaus@silentmode.my