# API Logging
---

#### C01. chapter1
1. meta
- event_id: string
- event_at: rfc3999
- event_level: enum=[debug, info, warning, error, critical]
- app_name: string
- app_version: string

2. api
- name: string
- at: rfc3339
- client: string
- values: map[string]string, keys=[query, status]
- identities: map[string]string, keys=[accountId, tokenId, ip, role]
- bizName: string
- code: string, list=[ok, warn, error, panic...]
- latencyMilli: float64
- meta: []byte

#### C02. chapter2

#### C03. chapter3
